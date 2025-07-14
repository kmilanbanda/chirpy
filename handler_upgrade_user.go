package main

import (
	"context"
	"net/http"
	"encoding/json"

	"github.com/kmilanbanda/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUpgradeUser(w http.ResponseWriter, req *http.Request) {
	apiKey, err := auth.GetAPIKey(req.Header)
	if err != nil {
		handleErrorResponse(w, http.StatusUnauthorized, "Error getting API Key")
		return
	} else if apiKey != cfg.polkaKey {
		handleErrorResponse(w, http.StatusUnauthorized, "API Key does not match")
		return
	}
	
	
	type data struct {
		UserID	string	`json:"user_id"`
	}

	type request struct {
		Event	string	`json:"event"`
		Data	data	`json:"data"`
	}

	var reqBody request
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		handleErrorResponse(w, http.StatusInternalServerError, "Error decoding request parameters")
		return
	}

	if reqBody.Event != "user.upgraded" {
		handleErrorResponse(w, http.StatusNoContent, "No Content")
		return
	}

	userID, err := uuid.Parse(reqBody.Data.UserID)
	if err != nil {
		handleErrorResponse(w, http.StatusInternalServerError, "Error parsing user ID")
		return
	}

	_, err = cfg.db.UpgradeUser(context.Background(), userID)
	if  err != nil {
		handleErrorResponse(w, http.StatusNotFound, "Error upgrading user")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
