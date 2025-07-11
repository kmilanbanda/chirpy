package main

import (
	"encoding/json"
	"net/http"
)

func handleErrorResponse(w http.ResponseWriter, httpStatusCode int, errorMessage string) {
	w.WriteHeader(httpStatusCode)
	resp := struct{ Error string `json:"error"` }{errorMessage}
	dat, _ := json.Marshal(resp)
	w.Write(dat)	
}
