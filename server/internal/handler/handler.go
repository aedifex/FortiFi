package handler

import (
	"encoding/json"
	"net/http"

	"github.com/aedifex/FortiFi/config"
	db "github.com/aedifex/FortiFi/internal/database"
	"github.com/aedifex/FortiFi/internal/firebase"
	"github.com/aedifex/FortiFi/internal/middleware"
	"github.com/aedifex/FortiFi/internal/requests"
	"go.uber.org/zap"
)

// Handler Wrapper
type RouteHandler struct {
	Log       	*zap.SugaredLogger
	Db    	  	*db.DatabaseConn
	Config 	  	*config.Config
	FcmClient 	*firebase.FcmClient
}

func (h *RouteHandler) NotifyIntrusion(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	subjectId, ok := request.Context().Value(middleware.UserIdContextKey).(string)
	if !ok  {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	
	// Decode body
	body := &requests.NotifyIntrusionRequest{}
	if request.Body == nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	err := json.NewDecoder(request.Body).Decode(body)
	if err != nil {
		h.Log.Errorf("error decoding body: %s", err)
		http.Error(writer, "failed to parse body", http.StatusBadRequest)
		return
	}
	if body.Event == nil {
		http.Error(writer, "invalid request", http.StatusBadRequest)
		return
	}

	// Store event in database
	event := body.Event
	event.Id = subjectId
	if event.Id == "" || event.Details == "" || event.Expires == "" || event.TS == "" || event.Type == "" || event.SrcIP == "" || event.DstIP == "" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	
	storeErr := h.Db.StoreEvent(event)
	if storeErr != nil {
		h.Log.Errorf("error storing event: %s", storeErr.Err)
		http.Error(writer, "failed to store event", storeErr.HttpStatus)
		return
	}
	h.Log.Infof("new event stored for user %s", subjectId)

	// Get notifications token
	fcmToken, fcmTokenErr := h.Db.GetFcmToken(subjectId)
	if fcmTokenErr != nil {
		h.Log.Errorf("error getting fcm token: %s", fcmTokenErr.Err)
		http.Error(writer, "failed to send notification", fcmTokenErr.HttpStatus)
		return
	}
	
	// Send Notification
	response, sendErr := h.FcmClient.SendMessage(fcmToken)
	if sendErr != nil {
		h.Log.Errorf("error sending notification: %s", sendErr)
		http.Error(writer, "failed to send notification", http.StatusInternalServerError)
		return
	}
	h.Log.Infof("Notification Sent: %s", response)
	writer.WriteHeader(http.StatusOK)
}

func (h *RouteHandler) UpdateWeeklyDistribution(writer http.ResponseWriter, request *http.Request) {

	if request.Method != http.MethodPost {
		http.Error(writer, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	
	subjectId, ok := request.Context().Value(middleware.UserIdContextKey).(string)
	if !ok {
		h.Log.Errorf("could not assert subjectId from context as string: %v", subjectId)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	
	// parse body
	body := &requests.UpdateWeeklyDistributionRequest{}
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
	if body.Benign < 0 || body.PortScan < 0 || body.DDoS < 0 {
		http.Error(writer, "invalid request", http.StatusBadRequest)
		return
	}
	
	updateErr := h.Db.UpdateWeeklyDistribution(subjectId, body.Benign, body.PortScan, body.DDoS)
	if updateErr != nil {
		h.Log.Errorf("error updating weekly distribution: %s", updateErr.Err)
		http.Error(writer, "unable to update weekly distribution", updateErr.HttpStatus)
		return
	}
	
	writer.WriteHeader(http.StatusOK)

}

func (h *RouteHandler) ResetWeeklyDistribution(writer http.ResponseWriter, request *http.Request) {

	if request.Method != http.MethodPost {
		http.Error(writer, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	
	subjectId, ok := request.Context().Value(middleware.UserIdContextKey).(string)
	if !ok {
		h.Log.Errorf("could not assert subjectId from context as string: %v", subjectId)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	body := &requests.ResetWeeklyDistributionRequest{}
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
	
	resetErr := h.Db.ResetWeeklyDistribution(subjectId, body.WeekTotal)
	if resetErr != nil {
		h.Log.Errorf("error resetting weekly distribution: %s", resetErr.Err)
		http.Error(writer, "unable to reset weekly distribution", resetErr.HttpStatus)
		return
	}
	
	writer.WriteHeader(http.StatusOK)
}
