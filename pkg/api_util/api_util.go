package api_util

import (
	"encoding/json"
	"io"
	"net/http"
)

// ErrorResponse represents a response containing an error message.
type ErrorResponse struct {
	Error string `json:"error"`
}

func RenderErrorResponse(w http.ResponseWriter, msg string, status int) {
	RenderResponse(w, ErrorResponse{Error: msg}, status)
}

func RenderResponse(w http.ResponseWriter, res interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")

	content, err := json.Marshal(res)
	if err != nil {
		// Do something with the error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)

	if _, err = w.Write(content); err != nil {
		// Do something with the error
	}
}

func ReadBody(r *http.Request, im any) error {
	js, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	err = json.Unmarshal(js, im)
	if err != nil {
		return err
	}
	return nil
}
