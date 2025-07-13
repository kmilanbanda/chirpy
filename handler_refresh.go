package main

import (
	"time"
	"context"
	"net/http"
	"encoding/json"
	"github.com/kmilanbanda/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		handleErrorResponse(w, http.StatusUnauthorized, "Error getting bearer token")
		return
	}
	refreshToken, err := cfg.db.GetRefreshTokenByToken(context.Background(), token)
	if err != nil {
		handleErrorResponse(w, http.StatusUnauthorized, "Unable to find refresh token")
		return
	} else if refreshToken.RevokedAt.Valid {
		handleErrorResponse(w, http.StatusUnauthorized, "Refresh token revoked")
		return
	} else if time.Now().After(refreshToken.ExpiresAt) {
		handleErrorResponse(w, http.StatusUnauthorized, "Refresh token expired")
		return
	}

	userID, err := cfg.db.GetUserFromRefreshToken(context.Background(), refreshToken.Token)
	if err != nil {
		handleErrorResponse(w, http.StatusInternalServerError, "Error getting user")
		return
	}

	user, err := cfg.db.GetUserByID(context.Background(), userID)
	if err != nil {
		handleErrorResponse(w, http.StatusInternalServerError, "Error getting user")
		return
	}
	accessToken, err := auth.MakeJWT(user.ID, cfg.secret, time.Hour)
	if err != nil {
		handleErrorResponse(w, http.StatusInternalServerError, "Error making JWT")
		return
	}

	resp := struct{
		Token	string	`json:"token"`
	}{
		Token:	accessToken,
	}
	dat, _ := json.Marshal(resp)
	w.Write(dat)
}
