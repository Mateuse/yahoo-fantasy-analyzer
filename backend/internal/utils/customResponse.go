package utils

import (
	"encoding/json"
	"net/http"

	"github.com/mateuse/yahoo-fantasy-analyzer/internal/models"
)

func CustomResponse(w http.ResponseWriter, code int, message string, details interface{}) {
	response := models.CustomResponse{
		Code:    code,
		Message: message,
		Details: details,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}
