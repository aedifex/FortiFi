package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/aedifex/FortiFi/config"
	"github.com/go-sql-driver/mysql"
)

func ConnectDatabase(config *config.Config) *sql.DB {
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
        log.Fatal(err)
    }

    pingErr := db.Ping()
    if pingErr != nil {
        log.Fatal(pingErr)
    }
    fmt.Println("Connected!")

	return db
}

