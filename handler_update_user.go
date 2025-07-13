package main

import (
	"time"
	"context"
	"net/http"
	"encoding/json"
	
	"github.com/google/uuid"
	"github.com/kmilanbanda/chirpy/internal/database"
	"github.com/kmilanbanda/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		handleErrorResponse(w, http.StatusUnauthorized, "Error getting bearer token")
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		handleErrorResponse(w, http.StatusUnauthorized, "Invalid")
		return
	}
	
	type request struct{
		Password	string	`json:"password"`
		Email		string	`json:"email"`
	}

	var reqBody request
	if err = json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		handleErrorResponse(w, http.StatusInternalServerError, "Error decoding request parameters")
		return
	}

	newHashedPassword, err := auth.HashPassword(reqBody.Password)
	if err != nil {
		handleErrorResponse(w, http.StatusInternalServerError, "Error hashing password")
		return
	}
	updateUserParams := database.UpdateUserParams {
		ID:		userID,
		Email:		reqBody.Email,
		HashedPassword:	newHashedPassword,
	}

	user, err := cfg.db.UpdateUser(context.Background(), updateUserParams)
	if err != nil {
		handleErrorResponse(w, http.StatusInternalServerError, "Error updating user")
		return	
	}

	resp := struct{
		ID		uuid.UUID 	`json:"id"`
		CreatedAt	time.Time	`json:"created_at"`
		UpdatedAt	time.Time	`json:"updated_at"`
		Email		string		`json:"email"`
	}{
		ID:		user.ID,
		CreatedAt:	user.CreatedAt,
		UpdatedAt:	user.UpdatedAt,
		Email:		user.Email,
	}
	dat, _  := json.Marshal(resp)
	w.Write(dat)
}
