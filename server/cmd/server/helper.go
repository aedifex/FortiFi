package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aedifex/FortiFi/config"
	"github.com/aedifex/FortiFi/internal/database"
	"github.com/aedifex/FortiFi/internal/firebase"
	"github.com/aedifex/FortiFi/internal/handler"
	"github.com/aedifex/FortiFi/internal/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ------------- Server Struct Logic ------------

type fortifiServer struct {
	httpServer 	*http.Server
	dbconn 		*database.DatabaseConn // Database wrapper
	config 		*config.Config
	logger 		*zap.SugaredLogger
	fcmClient 	*firebase.FcmClient
}

func newServer(config *config.Config) *fortifiServer {

	// Structured Logger 
	zapConfig := zap.NewProductionConfig()
	if os.Getenv("config") == "dev" {
		zapConfig = zap.NewDevelopmentConfig()
	}
	verbose := flag.Bool("verbose", false, "Enable logging output")
	flag.Parse()
	if *verbose {
		zapConfig.OutputPaths = []string{"server.log.json", os.Stdout.Name()}
	} else {
		zapConfig.OutputPaths = []string{"/dev/null"}
	}
	zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	zapLogger := zap.Must(zapConfig.Build()).Sugar()
	// connect to mysql database
	db, err := database.ConnectDatabase(config)
	if err != nil {
		zapLogger.Panicf("database connection error: %s", err)
	}
	zapLogger.Info("database connection successful")

	fcmClient,err := firebase.NewFirebaseMessagingClient(config)
	if err != nil {
		zapLogger.Panicf("error connecting to firebase: %s", err)
	}
	zapLogger.Info("connected to firebase client")

	// Route handling wrapper
	routeHandler := &handler.RouteHandler{
		Log: zapLogger,
		Db: db,
		Config: config,
		FcmClient: fcmClient,
	}
	
	// Register the Routes
	// All routes should be wrapped by middleware.Logging
	mux := http.NewServeMux()
	mux.HandleFunc("/NotifyIntrusion", middleware.Auth(config.SIGNING_KEY, zapLogger, routeHandler.NotifyIntrusion))
	// ? Should CreateUser be wrapped by Auth? Use the Pi init token to create a user
	mux.HandleFunc("/CreateUser", routeHandler.CreateUser)
	mux.HandleFunc("/Login", routeHandler.Login)
	mux.HandleFunc("/RefreshUser", routeHandler.RefreshUser)
	mux.HandleFunc("/RefreshPi", routeHandler.RefreshPi)
	mux.HandleFunc("/PiInit", routeHandler.PiInit)
	mux.HandleFunc("/UpdateFcm",  middleware.Auth(config.SIGNING_KEY, zapLogger, routeHandler.UpdateFcmToken))
	mux.HandleFunc("/GetUserEvents", middleware.Auth(config.SIGNING_KEY, zapLogger, routeHandler.GetUserEvents))
	mux.HandleFunc("/GetWeeklyDistribution", middleware.Auth(config.SIGNING_KEY, zapLogger, routeHandler.GetWeeklyDistribution))
	mux.HandleFunc("/UpdateWeeklyDistribution", middleware.Auth(config.SIGNING_KEY, zapLogger, routeHandler.UpdateWeeklyDistribution))
	mux.HandleFunc("/ResetWeeklyDistribution", middleware.Auth(config.SIGNING_KEY, zapLogger, routeHandler.ResetWeeklyDistribution))
	loggingMiddleware := middleware.Logging(zapLogger)
	corsMiddleware := middleware.CORSMiddleware(config.CORS_ORIGIN)
	serverHandler := corsMiddleware(loggingMiddleware(mux))
	httpServer := &http.Server{
		Addr: config.Port,
		Handler: serverHandler,
	}

	return &fortifiServer{
		httpServer: httpServer,
		dbconn: db,
		config: config,
		logger: zapLogger,
		fcmClient: fcmClient,
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
		if err := s.dbconn.Close(); err != nil {
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