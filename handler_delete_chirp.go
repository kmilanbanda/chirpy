package main

import (
	"context"
	"net/http"
	"github.com/kmilanbanda/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, req *http.Request) {
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
	chirpID, err := uuid.Parse(req.PathValue("chirpID"))
	if err != nil {
		handleErrorResponse(w, http.StatusInternalServerError, "Failed to parse chirp ID")
		return
	}

	chirp, err := cfg.db.GetChirp(context.Background(), chirpID)
	if err != nil {
		handleErrorResponse(w, http.StatusNotFound, "Error finding chirp")
		return
	}

	if chirp.UserID != validatedUserID {
		handleErrorResponse(w, http.StatusForbidden, "User does not own chirp")
		return
	}

	cfg.db.DeleteChirp(context.Background(), chirpID)
	w.WriteHeader(http.StatusNoContent)
}
