package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"fmt"
	"context"
	"github.com/kmilanbanda/chirpy/internal/database"
	"github.com/kmilanbanda/chirpy/internal/auth"
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

func (cfg *apiConfig) handlerPostChirp(w http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		handleErrorResponse(w, http.StatusUnauthorized, "Failed to read header")
		return
	}
	validatedUserID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		handleErrorResponse(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	w.Header().Set("Content-Type", "application/json")

	type request struct {
		Body 	string `json:"body"`
	}

	var reqBody request
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		handleErrorResponse(w, http.StatusInternalServerError, "Error decoding request parameters")
		return
	}

	if len(reqBody.Body) > cfg.maxChirpLength {
		handleErrorResponse(w, http.StatusBadRequest, "Chirp is too long")	
		return
	}

	createChirpParams := database.CreateChirpParams{
		Body:	reqBody.Body,
		UserID:	validatedUserID,
	}

	chirp, err := cfg.db.CreateChirp(context.Background(), createChirpParams)
	if err != nil {
		handleErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Error creating chirp: %v", err))
		return
	}
	
	w.WriteHeader(http.StatusCreated)
	dat, _ := json.Marshal(chirp)
	w.Write(dat)
}
