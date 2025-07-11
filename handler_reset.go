package main

import (
	"net/http"
	"context"
	"log"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(""))
		return
	}

	cfg.fileserverHits.Store(0)

	err := cfg.db.ResetUsers(context.Background())
	if err != nil {
		log.Fatalf("Failed to reset user database: %v", err)
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(""))
	if err != nil {
		log.Fatalf("Failed to write response for hits reset")
	}
}
