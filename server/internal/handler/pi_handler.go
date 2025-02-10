package handler

import (
	"encoding/json"
	"net/http"

	db "github.com/aedifex/FortiFi/internal/database"
	"github.com/aedifex/FortiFi/internal/requests"
	"github.com/aedifex/FortiFi/pkg/utils"
)

func (h *RouteHandler) PiInit(writer http.ResponseWriter, request *http.Request) {

	if request.Method != http.MethodPost {
		http.Error(writer, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get pi id
	pi := &requests.PiInitRequest{}
	err := json.NewDecoder(request.Body).Decode(pi)
	if err != nil {
		h.Log.Errorf("Decoding error in pi init: %s", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	if pi.Id == "" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	// Gen new tokens
	jwt, refresh, err := utils.GenTokenPair(h.Config.SIGNING_KEY, pi.Id)
	if err != nil {
		h.Log.Errorf("GenJwt Error: %s", err.Error())
		http.Error(writer, "pi init Error", http.StatusInternalServerError)
		return
	}

	// store the token
	storeErr := h.Db.StoreRefresh(refresh, db.PiRefreshTable)
	if storeErr != nil {
		h.Log.Errorf("error storing refresh token: %s", storeErr.Err)
		http.Error(writer, "login error", storeErr.HttpStatus)
		return
	}

	writer.Header().Add("jwt", jwt)
	writer.Header().Add("refresh", refresh.Token)
	writer.WriteHeader(http.StatusOK)
	h.writeResponse(writer, "init success")
}

func (h *RouteHandler) RefreshPi(writer http.ResponseWriter, request *http.Request) {

	suppliedId := request.URL.Query().Get("id")

	if request.Method != http.MethodGet {
		http.Error(writer, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	
	//construct token
	token := &db.RefreshToken{
		Token: 	request.Header.Get("Refresh"),
		Id: 	suppliedId,
	}

	//validate token
	err := h.Db.ValidateRefresh(token, db.PiRefreshTable)
	if err != nil {
		h.Log.Errorf("Refresh Token Err: %s", err.Err)
		writer.WriteHeader(err.HttpStatus)
		return
	}

	// Generate new tokens
	jwt, refresh, genTokenErr := utils.GenTokenPair(h.Config.SIGNING_KEY, suppliedId)
	if genTokenErr != nil {
		h.Log.Errorf("GenJwt Error: %s", genTokenErr.Error())
		http.Error(writer, "pi init Error", http.StatusInternalServerError)
		return
	}

	// store the token
	storeErr := h.Db.StoreRefresh(refresh, db.PiRefreshTable)
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
