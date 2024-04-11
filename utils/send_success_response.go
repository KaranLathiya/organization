package utils

import (
	"encoding/json"
	"net/http"
)

func SuccessMessageResponse(w http.ResponseWriter, statusCode int, response interface{}) {
	w.WriteHeader(statusCode)
	responseData, _ := json.MarshalIndent(response, "", "  ")
	w.Write(responseData)
}
