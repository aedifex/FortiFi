package database

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/aedifex/FortiFi/config"
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
    hashedPassword, err := HashPassword(user.Password)
    if err != nil {
        return http.StatusInternalServerError, fmt.Errorf("failed to hash password: %s", err.Error())
    }
    fmt.Printf("User: %v", user)
    preparedStatement, err := db.conn.Prepare(query)
    if err != nil {
        return http.StatusInternalServerError, fmt.Errorf("failed to create prepared statement: %s", err)
    }
    defer preparedStatement.Close()
    preparedStatement.Exec(user.Id, user.FirstName, user.LastName, user.Email, hashedPassword)
    return http.StatusCreated, nil
}


func (db *DatabaseConn) Login(user *User) (int, error) {
    query := "SELECT email,password FROM  USERS WHERE email = ?;"
    preparedStatement, err := db.conn.Prepare(query)
    if err != nil {
        return http.StatusInternalServerError, fmt.Errorf("failed to create prepared statement: %s", err)
    }
    defer preparedStatement.Close()
    if user.Email == "" {
        return http.StatusBadRequest, errors.New("email not provided")
    }
    if user.Password == "" {
        return http.StatusBadRequest, errors.New("password not provided")
    }

    res, err := preparedStatement.Query(user.Email)
    if err != nil {
        return http.StatusInternalServerError, err
    }
    if !res.Next() { return http.StatusNotFound, errors.New("email does not exist")}
    
    var storedPass string
    var email string
    if err := res.Scan(&email, &storedPass); err != nil {
        return http.StatusInternalServerError, err
    }
    if !ValidatePassword(storedPass, user.Password) {
        return http.StatusUnauthorized, errors.New("passwords do not match")
    }
    return http.StatusOK, nil
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