package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	db "github.com/aedifex/FortiFi/internal/database"
	"github.com/aedifex/FortiFi/internal/middleware"
	"github.com/aedifex/FortiFi/internal/requests"
	"github.com/aedifex/FortiFi/pkg/utils"
)

func (h *RouteHandler) CreateUser(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// parse body
	body := &requests.CreateUserRequest{}
	if request.Body == nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	err := json.NewDecoder(request.Body).Decode(body)
	if err != nil {
		http.Error(writer, "Failed to parse request body", http.StatusBadRequest)
		return
	}
	if body.User == nil {
		http.Error(writer, "invalid request", http.StatusBadRequest)
		return
	}

	user := body.User

	// Validate all required fields
	if user.Id == "" || user.FirstName == "" || user.LastName == "" || user.Email == "" || user.Password == "" {
		http.Error(writer, "missing fields", http.StatusBadRequest)
		return
	}

	insertErr := h.Db.InsertUser(user)
	if insertErr != nil {
		h.Log.Errorf("Error creating a new user: %s", insertErr.Err)
		http.Error(writer, "Failed to Create User", insertErr.HttpStatus)
		return
	}

	res := "Account Created"
	h.Log.Infof("New account created: %s", user.Email)
	writer.WriteHeader(http.StatusCreated)
	h.writeResponse(writer, res)
}

func (h *RouteHandler) Login(writer http.ResponseWriter, request *http.Request) {

	if request.Method != http.MethodPost {
		http.Error(writer, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse body
	body := &requests.LoginUserRequest{}
	if request.Body == nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	decodeErr := json.NewDecoder(request.Body).Decode(body)
	if decodeErr != nil {
		http.Error(writer, "Failed to parse request body", http.StatusBadRequest)
		return
	}
	if body.User == nil {
		http.Error(writer, "invalid request", http.StatusBadRequest)
		return
	}
	
	// Validate user
	user := body.User
	foundUser, err := h.Db.ValidateLogin(user)
	if err != nil {
		h.Log.Errorf("login error: %s", err.Err)
		http.Error(writer, "Login failed", err.HttpStatus)
		return
	}
	h.Log.Infof("Successful login for user %s", foundUser.Email)

	// generate auth tokens
	auth, refresh, tokenGenErr := utils.GenTokenPair(h.Config.SIGNING_KEY, foundUser.Id)
	if tokenGenErr != nil {
		h.Log.Errorf("GenJwt Error: %s", tokenGenErr.Error())
		http.Error(writer, "Login Error", http.StatusInternalServerError)
		return
	}

	// store the token
	storeErr := h.Db.StoreRefresh(refresh, db.UserRefreshTable)
	if storeErr != nil {
		h.Log.Errorf("error storing refresh token: %s", storeErr.Err)
		http.Error(writer, "login error", http.StatusInternalServerError)
		return
	}

	// set tokens in headers
	writer.Header().Add("jwt", auth)
	writer.Header().Add("refresh", refresh.Token)
	writer.WriteHeader(http.StatusOK)
	h.writeResponse(writer, "Login Success")
}

func (h *RouteHandler) RefreshUser(writer http.ResponseWriter, request *http.Request) {

	suppliedId := request.URL.Query().Get("id")
	if request.Method != http.MethodGet {
		http.Error(writer, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	

	// Serialize token
	refresh := &db.RefreshToken{
		Id: suppliedId,
		Token: request.Header.Get("Refresh"),
	}

	//Validate token
	err := h.Db.ValidateRefresh(refresh, db.UserRefreshTable)
	if err != nil {
		h.Log.Errorf("Refresh Token Err: %s", err.Err)
		writer.WriteHeader(err.HttpStatus)
		return
	}

	// generate auth tokens
	jwt, refresh, genTokenErr := utils.GenTokenPair(h.Config.SIGNING_KEY, suppliedId)
	if genTokenErr != nil {
		h.Log.Errorf("Token Gen Error: %s", genTokenErr.Error())
		http.Error(writer, "Refresh Error", http.StatusInternalServerError)
		return
	}

	// store the token
	storeErr := h.Db.StoreRefresh(refresh, db.UserRefreshTable)
	if storeErr != nil {
		h.Log.Errorf("error storing refresh token: %s", storeErr.Err)
		http.Error(writer, "login error", storeErr.HttpStatus)
		return
	}

	writer.Header().Add("jwt", jwt)
	writer.Header().Add("refresh", refresh.Token)
	writer.WriteHeader(http.StatusOK)
	res := "Valid Refresh Token"
	h.writeResponse(writer, res)

}

func (h *RouteHandler) UpdateFcmToken(writer http.ResponseWriter, request *http.Request){

	if request.Method != http.MethodPost {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	
	subjectId, ok := request.Context().Value(middleware.UserIdContextKey).(string)
	if !ok {
		h.Log.Errorf("could not assert subjectId from context as string: %v", subjectId)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	// parse body
	body := &requests.UpdateFcmRequest{}
	if request.Body == nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	err := json.NewDecoder(request.Body).Decode(body)
	if err != nil {
		h.Log.Errorf("json decode error: %s", err.Error())
		http.Error(writer, "unable to parse body", http.StatusBadRequest)
		return
	}
	if body.FcmToken == "" {
		http.Error(writer, "invalid request", http.StatusBadRequest)
		return
	}
	
	// Store new fcm token with user entry in database
	fcmErr := h.Db.UpdateFcmToken(subjectId, body.FcmToken)
	if fcmErr != nil {
		h.Log.Errorf("error updating fcm token: %s", fcmErr.Err)
		http.Error(writer, "unable to update fcm token", fcmErr.HttpStatus)
		return
	}
	h.Log.Infof("updated fcm for user: %s", subjectId)
	writer.WriteHeader(http.StatusAccepted)
	h.writeResponse(writer, "notifications token updated")
}

func (h *RouteHandler) GetUserEvents(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user ID from context (set by Auth middleware)
	userId, ok := request.Context().Value(middleware.UserIdContextKey).(string)
	if !ok {
		h.Log.Error("could not assert userId from context")
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Get events from database
	events, err := h.Db.GetUserEvents(userId)
	if err != nil {
		h.Log.Errorf("Error fetching user events: %s", err.Err)
		http.Error(writer, "Failed to fetch events", err.HttpStatus)
		return
	}

	// Write response
	writer.Header().Set("Content-Type", "application/json")
	encodeErr := json.NewEncoder(writer).Encode(map[string][]*db.Event{
		"events": events,
	})
	if encodeErr != nil {
		h.Log.Errorf("error encoding events: %s", encodeErr.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (h *RouteHandler) GetWeeklyDistribution(writer http.ResponseWriter, request *http.Request) {

	if request.Method != http.MethodGet {
		http.Error(writer, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	subjectId, ok := request.Context().Value(middleware.UserIdContextKey).(string)
	if !ok {
		h.Log.Errorf("could not assert subjectId from context as string: %v", subjectId)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	weeklyDistribution, err := h.Db.GetWeeklyDistribution(subjectId)
	if err != nil {
		h.Log.Errorf("error getting weekly distribution: %s", err.Err)
		http.Error(writer, "failed to get weekly distribution", err.HttpStatus)
		return
	}
	
	writer.Header().Set("Content-Type", "application/json")
	encodeErr := json.NewEncoder(writer).Encode(weeklyDistribution)
	if encodeErr != nil {
		h.Log.Errorf("error encoding weekly distribution: %s", encodeErr.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	
}

func (h *RouteHandler) GetDevices(writer http.ResponseWriter, request *http.Request) {

	if request.Method != http.MethodGet {
		http.Error(writer, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	subjectId, ok := request.Context().Value(middleware.UserIdContextKey).(string)
	if !ok {
		h.Log.Errorf("could not assert subjectId from context as string: %v", subjectId)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	devices, err := h.Db.GetDevices(subjectId)
	if err != nil {
		h.Log.Errorf("error getting devices: %s", err.Err)
		http.Error(writer, "failed to get devices", err.HttpStatus)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	encodeErr := json.NewEncoder(writer).Encode(devices)
	if encodeErr != nil {
		h.Log.Errorf("error encoding devices: %s", encodeErr.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	
}

func (h *RouteHandler) GetThreatAssistance(writer http.ResponseWriter, request *http.Request) {

	if request.Method != http.MethodGet {
		http.Error(writer, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	subjectId, ok := request.Context().Value(middleware.UserIdContextKey).(string)
	if !ok {
		h.Log.Errorf("could not assert subjectId from context as string: %v", subjectId)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	threatId, err := strconv.Atoi(request.URL.Query().Get("threatId"))
	if err != nil {
		h.Log.Errorf("error converting threatId to int: %s", err.Error())
		http.Error(writer, "invalid threatId", http.StatusBadRequest)
		return
	}
	
	event, dbErr := h.Db.GetThreatById(threatId, subjectId)
	if dbErr != nil {
		h.Log.Errorf("error getting threat by id: %s", dbErr.Err)
		http.Error(writer, "failed to get threat", dbErr.HttpStatus)
		return
	}

	llmResponse, err := h.OpenaiClient.GetHelpWithThreat(event)
	if err != nil {
		h.Log.Errorf("error getting help with threat: %s", err.Error())
		http.Error(writer, "failed to get help with threat", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	encodeErr := json.NewEncoder(writer).Encode(map[string]string{
		"response": llmResponse,
	})
	writer.WriteHeader(http.StatusOK)
	if encodeErr != nil {
		h.Log.Errorf("error encoding response: %s", encodeErr.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (h *RouteHandler) GetRecommendations(writer http.ResponseWriter, request *http.Request) {

	if request.Method != http.MethodGet {
		http.Error(writer, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	subjectId, ok := request.Context().Value(middleware.UserIdContextKey).(string)
	if !ok {
		h.Log.Errorf("could not assert subjectId from context as string: %v", subjectId)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	
	threatId, err := strconv.Atoi(request.URL.Query().Get("threatId"))
	if err != nil {
		h.Log.Errorf("error converting threatId to int: %s", err.Error())
		http.Error(writer, "invalid threatId", http.StatusBadRequest)
		return
	}

	event, dbErr := h.Db.GetThreatById(threatId, subjectId)
	if dbErr != nil {
		h.Log.Errorf("error getting threat by id: %s", dbErr.Err)
		http.Error(writer, "failed to get threat", dbErr.HttpStatus)
		return
	}

	llmResponse, err := h.OpenaiClient.GetRecommendations(event)
	if err != nil {
		h.Log.Errorf("error getting recommendations: %s", err.Error())
		http.Error(writer, "failed to get recommendations", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	encodeErr := json.NewEncoder(writer).Encode(map[string]string{
		"response": llmResponse,
	})
	if encodeErr != nil {
		h.Log.Errorf("error encoding response: %s", encodeErr.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusOK)

}

func (h *RouteHandler) GetMoreAssistance(writer http.ResponseWriter, request *http.Request) {

	if request.Method != http.MethodGet {
		http.Error(writer, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	subjectId, ok := request.Context().Value(middleware.UserIdContextKey).(string)
	if !ok {
		h.Log.Errorf("could not assert subjectId from context as string: %v", subjectId)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}		

	threatId, err := strconv.Atoi(request.URL.Query().Get("threatId"))
	if err != nil {
		h.Log.Errorf("error converting threatId to int: %s", err.Error())
		http.Error(writer, "invalid threatId", http.StatusBadRequest)
		return
	}
	
	event, dbErr := h.Db.GetThreatById(threatId, subjectId)
	if dbErr != nil {
		h.Log.Errorf("error getting threat by id: %s", dbErr.Err)
		http.Error(writer, "failed to get threat", dbErr.HttpStatus)
		return
	}
	
	body := &requests.GetMoreAssistanceRequest{}
	if request.Body == nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	err = json.NewDecoder(request.Body).Decode(body)
	if err != nil {
		h.Log.Errorf("error decoding request body: %s", err.Error())
		http.Error(writer, "invalid request", http.StatusBadRequest)
		return
	}

	if body.Query == ""{
		http.Error(writer, "invalid request", http.StatusBadRequest)
		return
	}
	
	llmResponse, err := h.OpenaiClient.GetMoreAssistance(body.Query, event)
	if err != nil {
		h.Log.Errorf("error getting more assistance: %s", err.Error())
		http.Error(writer, "failed to get more assistance", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	encodeErr := json.NewEncoder(writer).Encode(map[string]string{
		"response": llmResponse,
	})
	if encodeErr != nil {
		h.Log.Errorf("error encoding response: %s", encodeErr.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusOK)

}
// TODO implement this -- Should revoke refresh tokens
// func (h *RouteHandler) Logout(writer http.ResponseWriter, request *http.Request){
// 	return
// }