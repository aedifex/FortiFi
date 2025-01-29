package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aedifex/FortiFi/config"
	"github.com/aedifex/FortiFi/internal/database"
	"github.com/aedifex/FortiFi/internal/middleware"
	"github.com/aedifex/FortiFi/pkg/utils"
	"go.uber.org/zap"
)

// Handler Wrapper
type RouteHandler struct {
	Log *zap.SugaredLogger
	Db 	*database.DatabaseConn
	Config *config.Config
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
	h.Log.Infof("New account created: %s", user.Email)
	writer.WriteHeader(status)
	h.writeResponse(writer, res)
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

	foundUser, status, err := h.Db.ValidateLogin(user)
	if err != nil {
		h.Log.Warnf("login error: %s", err.Error())
		http.Error(writer, "Login failed", status)
		return
	}
	h.Log.Infof("Successful login for user %s", foundUser.Email)
	res := "Login Succcess"

	// generate auth tokens
	auth, refresh, refreshExp, err := utils.GenJwt(h.Config.SIGNING_KEY, foundUser.Id)
	if err != nil {
		h.Log.Warnf("GenJwt Error: %s", err.Error())
		http.Error(writer, "Login Error", http.StatusInternalServerError)
		return
	}

	// store the token
	storeErr := h.Db.StoreRefresh(refresh, foundUser.Id, refreshExp)
	if storeErr != nil {
		h.Log.Errorf("error storing refresh token: %s", storeErr.Error())
		http.Error(writer, "login error", http.StatusInternalServerError)
		return
	}

	// set tokens in headers
	writer.Header().Add("jwt", auth)
	writer.Header().Add("refresh", refresh)
	writer.WriteHeader(status)
	h.writeResponse(writer,res)
}

func (h *RouteHandler) Refresh(writer http.ResponseWriter, request *http.Request) {
	token := request.Header.Get("Refresh")
	jwt, refresh, err := h.Db.ValidateRefresh(h.Config.SIGNING_KEY,token)
	if err != nil {
		h.Log.Warnf("Refresh Token Err: %s", err.Error())
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	writer.Header().Add("jwt", jwt)
	writer.Header().Add("refresh", refresh)
	writer.WriteHeader(http.StatusOK)
	res := "Valid Refresh Token"
	h.writeResponse(writer, res)
}

func (h *RouteHandler) Protected(writer http.ResponseWriter, request *http.Request) {
	userId := request.Context().Value(middleware.UserIdContextKey)
	res := fmt.Sprintf("You have reached this endpoint as user: %s", userId)
	writer.WriteHeader(http.StatusOK)
	h.writeResponse(writer,res)
}