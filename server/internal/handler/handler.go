package handler

import (
	"encoding/json"
	"net/http"

	"github.com/aedifex/FortiFi/internal/database"
	"go.uber.org/zap"
)

// Handler Wrapper
type RouteHandler struct {
	Log *zap.SugaredLogger
	Db 	*database.DatabaseConn
}

func (h *RouteHandler) NotifyIntrusionHandler(writer http.ResponseWriter, request *http.Request){
	if request.Method != "POST" {
		http.Error(writer, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// Correct flow:
	// Request body should include information about the intrusion and device uid as it is entered in database 
	// Server should update database with event info
	// Server should send notification to user accordingly
	//    Need to get user associated with specific pi from database
}

func (h *RouteHandler) CreateUser(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		http.Error(writer, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	user := &database.User{}
	err := json.NewDecoder(request.Body).Decode(&user)

	if err != nil {
		http.Error(writer, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	status, err := h.Db.InsertUser(user) // fix logic so handlers belong to server object
	if err != nil {
		h.Log.Warnf("Error creating a new user: %s", err.Error())
		http.Error(writer, "Failed to Create User", status)
		return
	}

	res := "Account Created"
	writer.WriteHeader(status)
	writer.Write([]byte(res))
}

func (h *RouteHandler) Login(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		http.Error(writer, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	user := &database.User{}
	err := json.NewDecoder(request.Body).Decode(&user)

	if err != nil {
		http.Error(writer, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	status, err := h.Db.Login(user) // fix logic so handlers belong to server object
	if err != nil {
		h.Log.Warnf("login error: %s", err.Error())
		http.Error(writer, "Login failed", status)
		return
	}
	h.Log.Infof("Successful login for user %s", user.Email)
	res := "Login Succcess" // should send some form of auth token eventually
	writer.WriteHeader(status)
	writer.Write([]byte(res))
}