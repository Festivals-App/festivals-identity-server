package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

// respondJSON makes the response with payload as json format
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {

	resultMap := map[string]interface{}{"data": payload}
	response, err := json.Marshal(resultMap)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to retrieve content type")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(response))
}

// respondError makes the error response with payload as json format
func respondError(w http.ResponseWriter, code int, message string) {
	resultMap := map[string]interface{}{"error": message}
	response, err := json.Marshal(resultMap)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write([]byte(response))
}

func respondString(w http.ResponseWriter, code int, message string) {

	response := []byte(message)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(code)
	_, err := w.Write(response)
	if err != nil {
		log.Print(err.Error())
	}
}

func respondCode(w http.ResponseWriter, code int) {

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(code)
}
