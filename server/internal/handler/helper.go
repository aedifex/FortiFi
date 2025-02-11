package handler

import (
	"net/http"
)

func (h *RouteHandler) writeResponse(writer http.ResponseWriter, res string) {
	_, err := writer.Write([]byte(res))
	if err != nil {
		h.Log.Error(err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}
