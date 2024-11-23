package common

import (
	"encoding/json"
	"log"
	"net/http"
	"todoGoApi/types"
)

func ErrorResponse(w http.ResponseWriter, message string, statusCode int,functionName string) {
	log.Println("\033[31m" + "FUNCTION: " + functionName + " ERROR: " ,message + "\033[0m" )
	errResponse := types.ErrorResponse{
		Message: message,
		Code: statusCode,
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errResponse)
}