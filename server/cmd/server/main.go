package main

import (
	"log"
	"net/http"

	"github.com/aedifex/FortiFi/config"
	"github.com/aedifex/FortiFi/internal/handler"
)

func main() {

	config := config.SetConfig()

	server := newServer(config)
	
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Failed to start server on %s\n", config.Port)
	}

}

func newServer(config *config.Config) *http.Server{

	http.HandleFunc("/NotifyIntrusion", handler.NotifyIntrusionHandler)

	return &http.Server{
		Addr: config.Port,
	}
}