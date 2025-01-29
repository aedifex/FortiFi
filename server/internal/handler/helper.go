package handler

import (
	"fmt"
	"net/http"
)

func (h *RouteHandler) writeResponse(writer http.ResponseWriter, res string) {
	_,err := writer.Write([]byte(res));
	if err != nil {
		h.writeErr(writer, fmt.Errorf("error writing response: %s", err.Error()), http.StatusInternalServerError)
	}
}

func (h *RouteHandler) writeErr(writer http.ResponseWriter, err error, status int) {
	h.Log.Error(err)
	http.Error(writer, err.Error(), status)
}