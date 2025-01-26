package database

import (
	"database/sql"
	"fmt"

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

func (db *DatabaseConn) InsertUser(user *User) error {

    query := "INSERT into USERS VALUES (?,?,?,?,?)"

    preparedStatement, err := db.conn.Prepare(query)
    if err != nil {
        return fmt.Errorf("failed to create prepared statement: %s", err)
    }
    defer preparedStatement.Close()
    preparedStatement.Exec(user.Id, user.FirstName, user.LastName, user.Email, user.Password)
    return nil
}