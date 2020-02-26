package handlers

import (
	"encoding/json"
	"net/http"
)

type ErrorResponseBody struct {
	Error string `json:"error"`
}

type SuccessResponseBody struct {
	Message string `json:"error"`
}

func HttpResponseBody(w http.ResponseWriter, error string, code int) {
	var data interface{}

	if error == "" {
		data = SuccessResponseBody{
			Message: "success",
		}
	} else {
		data = &ErrorResponseBody{
			Error: error,
		}
	}

	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(data)
}
