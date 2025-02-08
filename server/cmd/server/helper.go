package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aedifex/FortiFi/config"
	"github.com/aedifex/FortiFi/internal/database"
	"github.com/aedifex/FortiFi/internal/handler"
	"github.com/aedifex/FortiFi/internal/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ------------- Server Struct Logic ------------

type fortifiServer struct {
	httpServer *http.Server
	dbconn *database.DatabaseConn // Database wrapper
	config *config.Config
	logger *zap.SugaredLogger
}

func newServer(config *config.Config) *fortifiServer {

	// Structured Logger 
	zapConfig := zap.NewProductionConfig()
	if os.Getenv("config") == "dev" {
		zapConfig = zap.NewDevelopmentConfig()
	}
	zapConfig.OutputPaths = []string{"server.log.json", os.Stdout.Name()}
	zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	zapLogger := zap.Must(zapConfig.Build()).Sugar()
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
	mux.HandleFunc("/RefreshUser", routeHandler.RefreshUser)
	mux.HandleFunc("/RefreshPi", routeHandler.RefreshPi)
	mux.HandleFunc("/PiInit", routeHandler.PiInit)
	mux.HandleFunc("/Protected", middleware.Auth(config.SIGNING_KEY, zapLogger, routeHandler.Protected))
	loggingMiddleware := middleware.Logging(zapLogger)

	httpServer := &http.Server{
		Addr: config.Port,
		Handler: loggingMiddleware(mux),
	}

	return &fortifiServer{
		httpServer: httpServer,
		dbconn: db,
		config: config,
		logger: zapLogger,
	}

}

func (s *fortifiServer) shutdown() {
    s.logger.Info("Starting server shutdown...")

    // HTTP server shutdown
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.logger.Error("Error during HTTP server shutdown:", err)
	}
	
	s.logger.Info("Http server shutdown complete")
	if s.dbconn != nil {
		s.logger.Info("closing database connection")
		if err := s.dbconn.Conn.Close(); err != nil {
			s.logger.Error("Error during HTTP server shutdown:", err)
		} else {
			s.logger.Info("Database connection closed")
		}
	}
    s.logger.Info("Server shutdown complete.")

	if err := s.logger.Sync(); err != nil {
		fmt.Printf("err logger: %s", err)
	}
}