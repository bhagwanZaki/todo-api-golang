package common

import (
	"encoding/json"
	"log"
	"net/http"
	"todoGoApi/types"
)

func ErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	log.Println("ERROR: ",message)
	errResponse := types.ErrorResponse{
		Message: message,
		Code: statusCode,
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errResponse)
}