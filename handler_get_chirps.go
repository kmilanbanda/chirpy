package main

import (
	"sort"
	"net/http"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/kmilanbanda/chirpy/internal/database"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	authorIDString := req.URL.Query().Get("author_id")
	sortOrder := req.URL.Query().Get("sort")

	var chirps []database.Chirp
	var err error
	if authorIDString != "" {
		userID, err := uuid.Parse(authorIDString)
		if err != nil {
			handleErrorResponse(w, http.StatusInternalServerError, "Error getting author ID")
		}
		chirps, err = cfg.db.GetChirpsByUser(context.Background(), userID)
	} else {
		chirps, err = cfg.db.GetChirps(context.Background())
	}
	if err != nil {
		handleErrorResponse(w, http.StatusInternalServerError, "Error getting chirps")
		return	
	}

	if sortOrder == "desc" {
		sort.Slice(chirps, func(i, j int) bool { return chirps[i].CreatedAt.After(chirps[j].CreatedAt) })
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
