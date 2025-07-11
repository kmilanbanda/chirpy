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
		w.WriteHeader(http.StatusInternalServerError)
		resp := struct{ Error string `json:"error"` }{"Error getting chirps"}
		dat, _ := json.Marshal(resp)
		w.Write(dat)
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
		w.WriteHeader(http.StatusInternalServerError)
		resp := struct{ Error string `json:"error"` }{"Error parsing UUID"}
		dat, _ := json.Marshal(resp)
		w.Write(dat)
		return		
	}

	chirp, err := cfg.db.GetChirp(context.Background(), uuid.UUID(chirpID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp := struct{ Error string `json:"error"` }{"Error getting chirp: id not found"}
		dat, _ := json.Marshal(resp)
		w.Write(dat)
		return	
	}

	w.WriteHeader(http.StatusOK)
	dat, _ := json.Marshal(chirp)
	w.Write(dat)
}
