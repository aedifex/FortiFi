package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/aedifex/FortiFi/config"
	"github.com/aedifex/FortiFi/internal/database"
	"github.com/aedifex/FortiFi/internal/handler"
)

func main() {

	config := config.SetConfig()

	server := NewServer(config)
  
	err := server.HttpServer.ListenAndServe()
	if err != nil {
		log.Fatalf("Failed to start server on port: %s\n", config.Port)
	}

}

// ------------- Server Struct Logic ------------

type FortifiServer struct {
	HttpServer *http.Server
	DBConn *sql.DB
}

func NewServer(config *config.Config) *FortifiServer {

	http.HandleFunc("/NotifyIntrusion", handler.NotifyIntrusionHandler)

	httpServer := &http.Server{
		Addr: config.Port,
	}

	return &FortifiServer{
		HttpServer: httpServer,
		DBConn: database.ConnectDatabase(config),
	}
  
}
