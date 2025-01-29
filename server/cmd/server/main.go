package main

import (
	"log"
	"net/http"
	"os"

	"github.com/aedifex/FortiFi/config"
	"github.com/aedifex/FortiFi/internal/database"
	"github.com/aedifex/FortiFi/internal/handler"
	"github.com/aedifex/FortiFi/internal/middleware"
	"go.uber.org/zap"
)

func main() {

	// Setup environment
	config := config.SetConfig()

	// Create new FortifiServer
	server := NewServer(config)

	// Dump logs on crash
	defer server.logger.Sync()

	// Start Server
	err := server.HttpServer.ListenAndServe()
	if err != nil {
		log.Fatalf("Failed to start server on port: %s\n", config.Port)
	}

}

// ------------- Server Struct Logic ------------

type FortifiServer struct {
	HttpServer *http.Server
	DBConn *database.DatabaseConn // Database wrapper
	config *config.Config
	logger *zap.SugaredLogger
}

func NewServer(config *config.Config) *FortifiServer {

	// Structured Logger 
	zapLogger := zap.Must(zap.NewProduction()).Sugar()
	if os.Getenv("config") == "dev" {
		zapLogger = zap.Must(zap.NewDevelopment()).Sugar()
	}

	httpServer := &http.Server{
		Addr: config.Port,
	}

	// connect to mysql database
	db := database.ConnectDatabase(zapLogger, config)
	
	// Route handling wrapper
	routeHandler := &handler.RouteHandler{
		Log: zapLogger,
		Db: db,
		Config: config,
	}
	
	// Register the Routes
	// All routes should be wrapped by middleware.Logging
	mux := http.NewServeMux()
	mux.HandleFunc("/NotifyIntrusion", routeHandler.NotifyIntrusionHandler)
	mux.HandleFunc("/CreateUser", routeHandler.CreateUser)
	mux.HandleFunc("/Login", routeHandler.Login)
	mux.HandleFunc("/Refresh", routeHandler.Refresh)
	mux.HandleFunc("/Protected", middleware.Auth(config.SIGNING_KEY, zapLogger, routeHandler.Protected))
	loggingMiddleware := middleware.Logging(zapLogger)
	httpServer.Handler = loggingMiddleware(mux)

	return &FortifiServer{
		HttpServer: httpServer,
		DBConn: db,
		config: config,
		logger: zapLogger,
	}

}
