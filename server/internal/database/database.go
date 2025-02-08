package database

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/aedifex/FortiFi/config"
	"github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
)

type DatabaseConn struct {
    Conn *sql.DB
}

const (
    UserRefreshTable =  "UserRefreshTokens"
    PiRefreshTable   =  "PiRefreshTokens"
)

func ConnectDatabase(log *zap.SugaredLogger, config *config.Config) *DatabaseConn {
	sqlConfig := mysql.Config{
        User:   config.DB_USER,
        Passwd: config.DB_PASS,
        Net:    "tcp",
        Addr:   config.DB_URL,
        DBName: config.DB_NAME,
    }

    // Get a database handle.
    db, err := sql.Open("mysql", sqlConfig.FormatDSN())
    if err != nil {
        log.Errorf("Error opening the database: %s", err.Error())
    }

    err = db.Ping()
    if err != nil {
        log.Errorf("Could not connect to DB: %s", err.Error())
    }
    log.Info("Database connection successful")

	return &DatabaseConn{
        Conn: db,
    }

}

func (db *DatabaseConn) InsertUser(user *User) (int,error) {
    //Check if user exists
    userExists := db.userExists(user)
    if userExists != nil {
        return http.StatusConflict, userExists
    }
    
    query := "INSERT into USERS VALUES (?,?,?,?,?);"

    // Hash the password
    hashedPassword, err := HashString(user.Password)
    if err != nil {
        return http.StatusInternalServerError, fmt.Errorf("failed to hash password: %s", err.Error())
    }

    // prepare statement
    preparedStatement, err := db.Conn.Prepare(query)
    if err != nil {
        return http.StatusInternalServerError, fmt.Errorf("failed to create prepared statement: %s", err)
    }
    defer preparedStatement.Close()

    // execute statement
    _, err = preparedStatement.Exec(user.Id, user.FirstName, user.LastName, user.Email, hashedPassword)
    if err != nil {
        return http.StatusInternalServerError, fmt.Errorf("failed to insert user: %s", err.Error())
    }

    return http.StatusCreated, nil
}

// Takes in a user object and compares the passwords to the matching user in the database.
//
// Returns a found user, http status code, and nil if the user is validated.
// Otherwise returns nil, error code, and a corresponding error.
func (db *DatabaseConn) ValidateLogin(user *User) (*User, int, error) {

    query := "SELECT * FROM  USERS WHERE email = ?;"

    // create prepared statment
    preparedStatement, err := db.Conn.Prepare(query)
    if err != nil {
        return nil, http.StatusInternalServerError, fmt.Errorf("failed to create prepared statement for login: %s", err)
    }
    defer preparedStatement.Close()

    // input validation
    if user.Email == "" {
        return nil, http.StatusBadRequest, errors.New("email not provided")
    }
    if user.Password == "" {
        return nil, http.StatusBadRequest, errors.New("password not provided")
    }

    // Query database
    res, err := preparedStatement.Query(user.Email)
    if err != nil {
        return nil, http.StatusInternalServerError, err
    }

    // check if user exists
    if !res.Next() { return nil, http.StatusUnauthorized, errors.New("email does not exist")}

    // store the user the query returned
    foundUser := &User{}
    if err := res.Scan(&foundUser.Id, &foundUser.FirstName, &foundUser.LastName, &foundUser.Email, &foundUser.Password); err != nil {
        return nil, http.StatusInternalServerError, err
    }

    // validate the user
    if !HashMatch(foundUser.Password, user.Password) {
        return nil, http.StatusUnauthorized, errors.New("passwords do not match")
    }
    return foundUser, http.StatusOK, nil
}

// Validates a refresh token by finding a match in the database.
// Returns error if token is invalid, otherwise nil.
func (db *DatabaseConn) ValidateRefresh(token string, tablename string, subjectId string) error {
    
    // drop expired tokens from respective table
    query := fmt.Sprintf("DELETE FROM %s WHERE expires < Now();", tablename)
    _, err := db.Conn.Exec(query)
    if err != nil {
        return fmt.Errorf("error deleting expired tokens: %s", err.Error())
    }
    
    // check if matching refresh token is in the database
    query = fmt.Sprintf("SELECT token_hash FROM %s WHERE id = ? LIMIT 1;",tablename)

    preparedStatement, err := db.Conn.Prepare(query)
    if err != nil {
        return fmt.Errorf("error preparing refreshtoken query: %s", err)
    }

    rows, err := preparedStatement.Query(subjectId)
    if err != nil {
        return fmt.Errorf("error executing query: %s",err.Error())
    }

    if !rows.Next() {
        return errors.New("token for user does not exist")
    }

    // Get token info
    storedHashedToken := ""
    err = rows.Scan(&storedHashedToken)
    if err != nil {
        return err
    }
    
    if !HashMatch(storedHashedToken, token) {
        return errors.New("invalid token")
    }

    return nil
}

// Stores a given token hash and subject id with exp time 7 days from current time the given tablename.
// Returns nil on success or non-nil error on failure.
func (db *DatabaseConn) StoreRefresh(token string, subject string, tablename string) error {

    expTime := time.Now().Add(time.Hour*24*7)

    // Delete existing keys for this user
    query := fmt.Sprintf("DELETE FROM %s WHERE id = ? LIMIT 1;", tablename)
    preparedStatement, err := db.Conn.Prepare(query)
    if err != nil {
        return fmt.Errorf("failed to prepare drop duplicates statement: %s", err.Error())
    }
    defer preparedStatement.Close()

    _, err = preparedStatement.Exec(subject)
    if err != nil {
        return fmt.Errorf("error executing drop dup query: %s", err.Error())
    }
    
    // Insert new token into database
    query = fmt.Sprintf("INSERT INTO %s VALUES (?, ?, ?);", tablename) // token, id, expires
    preparedStatement, err = db.Conn.Prepare(query)
    if err != nil {
        return fmt.Errorf("failed to prepare INSERT token statement: %s", err.Error())
    }
    defer preparedStatement.Close()

    formattedTime := expTime.Format("2006-01-02 15:04:05")
    tokenHash, err := HashString(token)
    if err != nil {
        return err
    }
    _, err = preparedStatement.Exec(tokenHash, subject, formattedTime)
    if err != nil {
        return fmt.Errorf("failed to store refresh token: %s", err.Error())
    }
    return nil
}


// Checks if a given user exists in the database.
// Returns error if the user email or id exists, nil otherwise.
func (db *DatabaseConn) userExists(user *User) error {
    // check email
    queryEmail := "SELECT * FROM USERS WHERE email = ?;"
    preparedStatementEmail, _ := db.Conn.Prepare(queryEmail)
    defer preparedStatementEmail.Close()
    res, err := preparedStatementEmail.Query(user.Email)
    if err != nil { return fmt.Errorf("error executing query: %s", err.Error()) }
    defer res.Close()
    if res.Next() { return errors.New("user already exists")}
    
    // check id
    queryId := "SELECT * FROM USERS WHERE id = ?;"
    preparedStatementId, _ := db.Conn.Prepare(queryId)
    defer preparedStatementId.Close()
    res, err = preparedStatementId.Query(user.Id)
    if err != nil { return fmt.Errorf("error executing query: %s", err.Error()) }
    defer res.Close()
    if res.Next() { return errors.New("user already exists")}

    return nil
}