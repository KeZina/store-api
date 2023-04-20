package helper

import (
	"encoding/json"
	"log"
	"net/http"
)

func SendError(w http.ResponseWriter, code int, message string) {
	log.Println(message)

	w.WriteHeader(code)
	w.Write([]byte(message))
}

func SendJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		SendError(w, http.StatusInternalServerError, err.Error())

		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
