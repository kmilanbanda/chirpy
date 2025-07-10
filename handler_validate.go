package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func getProfaneWords() map[string]struct{} {
	words := map[string]struct{}{
		"kerfuffle": 	{},
		"sharbert":	{},
		"fornax":	{},
	}

	return words
}

func censorProfanity(text string, profanities map[string]struct{}) string {
	words := strings.Split(text, " ")
	for i, word := range words {
		_, exists := profanities[strings.ToLower(word)]
		if exists {
			words[i] = "****"
		}
	}

	return strings.Join(words, " ")
}

func handlerValidateChirp(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	type request struct {
		Body string `json:"body"`
	}

	var reqBody request
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp := struct{ Error string `json:"error"` }{"Error decoding request parameters"}
		dat, _ := json.Marshal(resp)
		w.Write(dat)
		return
	} 

	if len(reqBody.Body) > 140 {
		w.WriteHeader(http.StatusBadRequest)
		resp := struct{ Error string `json:"error"` }{"Chirp is too long"}
		dat, _ := json.Marshal(resp)
		w.Write(dat)
		return
	}
	
	w.WriteHeader(http.StatusOK)
	resp := struct{ CleanedBody string `json:"cleaned_body"` }{censorProfanity(reqBody.Body, getProfaneWords())}
	dat, _ := json.Marshal(resp)
	w.Write(dat)
}
