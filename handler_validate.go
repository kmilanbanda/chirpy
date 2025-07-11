package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"fmt"
	"context"
	"github.com/kmilanbanda/chirpy/internal/database"

	"github.com/google/uuid"
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

func (cfg *apiConfig) handlerValidateChirp(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	type request struct {
		Body 	string `json:"body"`
		UserID	uuid.UUID `json:"user_id"`
	}

	var reqBody request
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		handleErrorResponse(w, http.StatusInternalServerError, "Error decoding request parameters")
		return
	}

	if len(reqBody.Body) > 140 {
		handleErrorResponse(w, http.StatusBadRequest, "Chirp is too long")	
		return
	}

	createChirpParams := database.CreateChirpParams{
		Body:	reqBody.Body,
		UserID:	reqBody.UserID,
	}

	chirp, err := cfg.db.CreateChirp(context.Background(), createChirpParams)
	if err != nil {
		handleErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Error creating chirp: %w", err))
		return
	}
	
	w.WriteHeader(http.StatusCreated)
	dat, _ := json.Marshal(chirp)
	w.Write(dat)
}
