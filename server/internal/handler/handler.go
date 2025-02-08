package handler

import (
	"fmt"
	"net/http"

	"github.com/aedifex/FortiFi/config"
	db "github.com/aedifex/FortiFi/internal/database"
	"github.com/aedifex/FortiFi/internal/middleware"
	"go.uber.org/zap"
)

// Handler Wrapper
type RouteHandler struct {
	Log    *zap.SugaredLogger
	Db     *db.DatabaseConn
	Config *config.Config
}

func (h *RouteHandler) NotifyIntrusionHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	// TODO Implement this as protected route
	// Correct flow:
	// Request body should include information about the intrusion and device uid as it is entered in database
	// Server should update database with event info
	// Server should send notification to user accordingly
	//    Need to get user associated with specific pi from database
}

func (h *RouteHandler) Protected(writer http.ResponseWriter, request *http.Request) {
	userId := request.Context().Value(middleware.UserIdContextKey)
	res := fmt.Sprintf("You have reached this endpoint as user: %s", userId)
	writer.WriteHeader(http.StatusOK)
	h.writeResponse(writer, res)
}
