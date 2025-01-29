package database

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/aedifex/FortiFi/config"
	"github.com/aedifex/FortiFi/pkg/utils"
	"github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
)

type DatabaseConn struct {
    conn *sql.DB
}

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
        conn: db,
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
    hashedPassword, err := HashPassword(user.Password)
    if err != nil {
        return http.StatusInternalServerError, fmt.Errorf("failed to hash password: %s", err.Error())
    }

    // prepare statement
    preparedStatement, err := db.conn.Prepare(query)
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


func (db *DatabaseConn) ValidateLogin(user *User) (*User, int, error) {

    query := "SELECT * FROM  USERS WHERE email = ?;"

    // create prepared statment
    preparedStatement, err := db.conn.Prepare(query)
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
    if !res.Next() { return nil, http.StatusNotFound, errors.New("email does not exist")}

    // store the user the query returned
    foundUser := &User{}
    if err := res.Scan(&foundUser.Id, &foundUser.FirstName, &foundUser.LastName, &foundUser.Email, &foundUser.Password); err != nil {
        return nil, http.StatusInternalServerError, err
    }

    // validate the user
    if !ValidatePassword(foundUser.Password, user.Password) {
        return nil, http.StatusUnauthorized, errors.New("passwords do not match")
    }
    return foundUser, http.StatusOK, nil
}

func (db *DatabaseConn) ValidateRefresh(key string, signedToken string) (string, string, error) {
    // drop expired tokens
    query := "DELETE FROM RefreshTokens WHERE expires < Now();"
    _, err := db.conn.Exec(query)
    if err != nil {
        return "","",fmt.Errorf("error deleting expired tokens: %s", err.Error())
    }

    // check if signed token is in the database
    signedTokenQuery := "SELECT token,FK_UserId FROM RefreshTokens WHERE token = ?;"
    preparedStatement, err := db.conn.Prepare(signedTokenQuery)
    if err != nil {
        return "","",fmt.Errorf("error preparing refreshtoken query: %s", err)
    }
    rows, err := preparedStatement.Query(signedToken)
    if err != nil {
        return "","", fmt.Errorf("error executing query: %s",err.Error())
    }
    if !rows.Next() {
        return "","", errors.New("invalid token")
    }

    // Get token info
    token := &Token{}
    scanErr := rows.Scan(&token.Token,&token.FK_UserId)
    if scanErr != nil {
        return "","",fmt.Errorf("error scanning results from refreshtokens: %s", scanErr)
    }

    // userId <- Get FK_UserId from refreshtokens table where token == signedToken
    userId := token.FK_UserId
    auth, refresh, expAt, err := utils.GenJwt(key, userId)
    if err != nil {
        return "","",err
    }
    // Store token
    storeTokenErr := db.StoreRefresh(refresh, userId, expAt)
    if storeTokenErr != nil {
        return "","",storeTokenErr
    }
    return auth, refresh, nil
}

func (db *DatabaseConn) StoreRefresh(token string, userId string, exp time.Time) error {

    // Delete existing keys for this user
    dropDup := "DELETE FROM RefreshTokens WHERE FK_UserId = ?;"
    dropDupPrep, dropDupPrepErr := db.conn.Prepare(dropDup)
    if dropDupPrepErr != nil {
        return fmt.Errorf("failed to prepare drop duplicates statement: %s", dropDupPrepErr.Error())
    }
    defer dropDupPrep.Close()
    _, dropDupErr := dropDupPrep.Exec(userId)
    if dropDupErr != nil {
        return fmt.Errorf("error executing drop dup query: %s", dropDupErr.Error())
    }
    
    // Insert new token into database
    query := "INSERT INTO RefreshTokens VALUES (?, ?, ?)" // token, id, expires
    preparedStatement, _ := db.conn.Prepare(query)
    defer preparedStatement.Close()
    formattedTime := exp.Format("2006-01-02 15:04:05")
    _, err := preparedStatement.Exec(token, userId, formattedTime)
    if err != nil {
        return fmt.Errorf("failed to store refresh token: %s", err.Error())
    }
    return nil
}

func (db *DatabaseConn) userExists(user *User) error {
    // check email
    queryEmail := "SELECT * FROM USERS WHERE email = ?;"
    preparedStatementEmail, _ := db.conn.Prepare(queryEmail)
    defer preparedStatementEmail.Close()
    res, err := preparedStatementEmail.Query(user.Email)
    if err != nil { return fmt.Errorf("error executing query: %s", err.Error()) }
    defer res.Close()
    if res.Next() { return errors.New("user already exists")}
    
    // check id
    queryId := "SELECT * FROM USERS WHERE id = ?;"
    preparedStatementId, _ := db.conn.Prepare(queryId)
    defer preparedStatementId.Close()
    res, err = preparedStatementId.Query(user.Id)
    if err != nil { return fmt.Errorf("error executing query: %s", err.Error()) }
    defer res.Close()
    if res.Next() { return errors.New("user already exists")}

    return nil
}