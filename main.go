package main

import _ "github.com/lib/pq"

import (
	"net/http"
	"log"
	"sync/atomic"
	"fmt"
	"encoding/json"	
	"strings"
	"os"
	"context"
	"database/sql"
	"github.com/kmilanbanda/chirpy/internal/database"
	"github.com/joho/godotenv"
	"github.com/google/uuid"
	"time"
)

type Chirp struct {
	body	string `json:"body"`
}

type apiConfig struct  {
	fileserverHits 	atomic.Int32
	db		*database.Queries
	platform	string
}

func getProfaneWords() map[string]struct{} {
	words := map[string]struct{}{
		"kerfuffle": 	{},
		"sharbert":	{},
		"fornax":	{},
	}

	return words
}

func censorProfanity(text string, profanities map[string]struct{}) string {
	words := strings.Split(text, " ")
	for i, word := range words {
		_, exists := profanities[strings.ToLower(word)]
		if exists {
			words[i] = "****"
		}
	}

	return strings.Join(words, " ")
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func handlerFunc(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(200)

	_, err := writer.Write([]byte("OK"))
	if err != nil {
		log.Fatalf("Failed to write response")
	}
} 

func handlerValidateChirp(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	type request struct {
		Body string `json:"body"`
	}

	var reqBody request
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp := struct{ Error string `json:"error"` }{"Error decoding request parameters"}
		dat, _ := json.Marshal(resp)
		w.Write(dat)
		return
	} 

	if len(reqBody.Body) > 140 {
		w.WriteHeader(http.StatusBadRequest)
		resp := struct{ Error string `json:"error"` }{"Chirp is too long"}
		dat, _ := json.Marshal(resp)
		w.Write(dat)
		return
	}
	
	w.WriteHeader(http.StatusOK)
	resp := struct{ CleanedBody string `json:"cleaned_body"` }{censorProfanity(reqBody.Body, getProfaneWords())}
	dat, _ := json.Marshal(resp)
	w.Write(dat)
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	type request struct {
		Email string `json:"email"`
	}

	var reqBody request
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp := struct{ Error string `json:"error"` }{"Error decoding request parameters"}
		dat, _ := json.Marshal(resp)
		w.Write(dat)
		return	
	}

	user, err := cfg.db.CreateUser(context.Background(), reqBody.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp := struct{ Error string `json:"error"` }{fmt.Sprintf("Error creating user: %v", err)}
		dat, _ := json.Marshal(resp)
		w.Write(dat)
		return
	}

	w.WriteHeader(http.StatusCreated)
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

func (cfg *apiConfig) handlerHits(w http.ResponseWriter, r *http.Request) {
	htmlRaw := "<html>\n  <body>\n    <h1>Welcome, Chirpy Admin</h1>\n    <p>Chirpy has been visited %d times!</p>\n  </body>\n</html>"
	htmlResponse := fmt.Sprintf(htmlRaw, cfg.fileserverHits.Load())
	
	w.Header().Set("Content-Type", "text/html")
	_, err := w.Write([]byte(htmlResponse))
	if err != nil {
		log.Fatalf("Failed to write response for hits")
	}
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(""))
		return
	}
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

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	pform := os.Getenv("PLATFORM")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Fatal error occured during connection to database: %w", err)
	}
	dbQueries := database.New(db)

	serveMux := http.NewServeMux()
	server := &http.Server{
		Addr:		":8080",
		Handler: 	serveMux,
	}
	cfg := &apiConfig{
		db:		dbQueries,
		platform:	pform,
	}


	fileHandler := http.FileServer(http.Dir("."))
	serveMux.HandleFunc("GET /api/healthz", handlerFunc)
	serveMux.Handle("/app/", http.StripPrefix("/app/", cfg.middlewareMetricsInc(fileHandler)))
	serveMux.HandleFunc("GET /admin/metrics", cfg.handlerHits)
	serveMux.HandleFunc("POST /api/users", cfg.handlerCreateUser)
	serveMux.HandleFunc("POST /admin/reset", cfg.handlerReset)
	serveMux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)
	server.ListenAndServe()
}
