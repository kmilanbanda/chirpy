package main

import (
	"net/http"
	"encoding/json"
	"time"
	"context"

	"github.com/google/uuid"
	"github.com/kmilanbanda/chirpy/internal/auth"
	"github.com/kmilanbanda/chirpy/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	type request struct {
		Password 		string	`json:"password"`
		Email			string	`json:"email"`
	}

	var reqBody request
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		handleErrorResponse(w, http.StatusInternalServerError, "Error decoding request parameters")
		return
	}

	user, err := cfg.db.GetUserByEmail(context.Background(), reqBody.Email)
	if err != nil {
		handleErrorResponse(w, http.StatusNotFound, "Error finding user")
		return
	}

	err = auth.CheckPasswordHash(user.HashedPassword, reqBody.Password)
	if err != nil {
		handleErrorResponse(w, http.StatusUnauthorized, "Incorrect Password")
		return
	}

	accessTokenDuration := time.Hour * 1
	refreshTokenDuration := time.Hour * 24 * 60

	token, err := auth.MakeJWT(user.ID, cfg.secret, accessTokenDuration)
	if err != nil {
		handleErrorResponse(w, http.StatusInternalServerError, "Error making JSON web token")
		return
	}
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		handleErrorResponse(w, http.StatusInternalServerError, "Error getting refresh token")
		return
	}
	createRefreshTokenParams := database.CreateRefreshTokenParams{
		Token:		refreshToken,
		UserID:		user.ID,
		ExpiresAt:	time.Now().Add(refreshTokenDuration),
	}
	cfg.db.CreateRefreshToken(context.Background(), createRefreshTokenParams)
	
	w.WriteHeader(http.StatusOK)
	resp := struct{
		ID		uuid.UUID 	`json:"id"`
		CreatedAt	time.Time	`json:"created_at"`
		UpdatedAt	time.Time	`json:"updated_at"`
		Email		string		`json:"email"`
		Token		string		`json:"token"`
		RefreshToken	string		`json:"refresh_token"`
	}{
		ID:		user.ID,
		CreatedAt:	user.CreatedAt,
		UpdatedAt:	user.UpdatedAt,
		Email:		user.Email,
		Token:		token,
		RefreshToken:	refreshToken,
	}
	dat, _  := json.Marshal(resp)
	w.Write(dat)	
}
