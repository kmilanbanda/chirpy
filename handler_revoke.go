package main

import (
	"context"
	"net/http"
	
	"github.com/kmilanbanda/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		handleErrorResponse(w, http.StatusInternalServerError, "Error getting bearer token")
		return
	}
	err = cfg.db.RevokeToken(context.Background(), token)
	if err != nil {
		handleErrorResponse(w, http.StatusInternalServerError, "Error revoking token")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
