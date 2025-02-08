package main

import (
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/aedifex/FortiFi/config"
)

func main() {

	// Setup environment
	config := config.SetConfig()

	// Create new FortifiServer
	server := newServer(config)

	// shutdown handling channel
	shutdownChan := make(chan os.Signal,1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Start server
	go func() {
		if err := server.httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed){
			server.logger.Fatalf("server err: %s\n", err)
			os.Exit(1)
		}
	}()
	server.logger.Infof("Server running on port: %s", server.httpServer.Addr)

	// block until ctrl+c
	// * Program exits with code 0 when programatically interrupted but not via ctrl+c
	// go func () {shutdownChan <- os.Interrupt}()
	<-shutdownChan
	server.shutdown()
	os.Exit(0)
}