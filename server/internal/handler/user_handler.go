package handler

import (
	"encoding/json"
	"net/http"

	db "github.com/aedifex/FortiFi/internal/database"
	"github.com/aedifex/FortiFi/pkg/utils"
)


func (h *RouteHandler) CreateUser(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	user := &db.User{}
	err := json.NewDecoder(request.Body).Decode(user)
	if err != nil {
		http.Error(writer, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// Validate all required fields
	if user.FirstName == "" || user.LastName == "" || user.Email == "" || user.Password == "" {
		http.Error(writer, "Missing required fields: first name, last name, email, and password", http.StatusBadRequest)
		return
	}

	status, err := h.Db.InsertUser(user)
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
	if request.Method != http.MethodPost {
		http.Error(writer, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	user := &db.User{}
	err := json.NewDecoder(request.Body).Decode(user)

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
	auth, refresh, err := utils.GenTokenPair(h.Config.SIGNING_KEY, foundUser.Id)
	if err != nil {
		h.Log.Warnf("GenJwt Error: %s", err.Error())
		http.Error(writer, "Login Error", http.StatusInternalServerError)
		return
	}

	// store the token
	storeErr := h.Db.StoreRefresh(refresh, foundUser.Id, db.UserRefreshTable)
	if storeErr != nil {
		h.Log.Errorf("error storing refresh token: %s", storeErr.Error())
		http.Error(writer, "login error", http.StatusInternalServerError)
		return
	}

	// set tokens in headers
	writer.Header().Add("jwt", auth)
	writer.Header().Add("refresh", refresh)
	writer.WriteHeader(http.StatusOK)
	h.writeResponse(writer, res)
}

func (h *RouteHandler) RefreshUser(writer http.ResponseWriter, request *http.Request) {

	suppliedId := request.URL.Query().Get("id")
	if request.Method != http.MethodGet {
		http.Error(writer, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	token := request.Header.Get("Refresh")
	err := h.Db.ValidateRefresh(token, db.UserRefreshTable, suppliedId)
	if err != nil {
		h.Log.Warnf("Refresh Token Err: %s", err.Error())
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	// generate auth tokens
	jwt, refresh, err := utils.GenTokenPair(h.Config.SIGNING_KEY, suppliedId)
	if err != nil {
		h.Log.Warnf("Token Gen Error: %s", err.Error())
		http.Error(writer, "Refresh Error", http.StatusInternalServerError)
		return
	}

	// store the token
	storeErr := h.Db.StoreRefresh(refresh, suppliedId, db.UserRefreshTable)
	if storeErr != nil {
		h.Log.Errorf("error storing refresh token: %s", storeErr.Error())
		http.Error(writer, "login error", http.StatusInternalServerError)
		return
	}

	writer.Header().Add("jwt", jwt)
	writer.Header().Add("refresh", refresh)
	writer.WriteHeader(http.StatusOK)
	res := "Valid Refresh Token"
	h.writeResponse(writer, res)

}


//TODO implement this -- Should revoke refresh tokens
// func (h *RouteHandler) Logout(writer http.ResponseWriter, request *http.Request){
// 	return
// }