package response

import (
	"encoding/json"
	"log"
	"net/http"
)

func RespondWithError(w http.ResponseWriter, code int, message string) {
	errorMessage := struct {
		Error string `json:"error"`
	}{
		Error: message,
	}

	RespondWithJson(w, code, errorMessage)

}

func RespondWithJson(w http.ResponseWriter, code int, message any) {
	respJson, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error encoding response message: %s", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	w.Write(respJson)
}
