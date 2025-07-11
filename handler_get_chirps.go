package main

import (
	"net/http"
	"context"
	"encoding/json"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	chirps, err := cfg.db.GetChirps(context.Background())
	if err != nil {
		handleErrorResponse(w, http.StatusInternalServerError, "Error getting chirps")
		return	
	}

	w.WriteHeader(http.StatusOK)
	dat, _ := json.Marshal(chirps)
	w.Write(dat)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	chirpID, err := uuid.Parse(req.PathValue("chirpID"))
	if err != nil {
		handleErrorResponse(w, http.StatusNotFound, "Error parsing UUID")
		return
	}

	chirp, err := cfg.db.GetChirp(context.Background(), uuid.UUID(chirpID))
	if err != nil {
		handleErrorResponse(w, http.StatusNotFound, "Error getting chirp: id not found")	
		return	
	}

	w.WriteHeader(http.StatusOK)
	dat, _ := json.Marshal(chirp)
	w.Write(dat)
}
