package main

import (
	"net/http"
	"log"
	"sync/atomic"
	"fmt"
	"encoding/json"
	"os"
	"context"
	"database/sql"
	"time"
	"strconv"

	"github.com/kmilanbanda/chirpy/internal/auth"
	"github.com/kmilanbanda/chirpy/internal/database"
	"github.com/joho/godotenv"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type apiConfig struct  {
	fileserverHits 	atomic.Int32
	db		*database.Queries
	platform	string
	maxChirpLength	int
	secret		string
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



func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	type request struct {
		Password string `json:"password"`
		Email string `json:"email"`
	}

	var reqBody request
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		handleErrorResponse(w, http.StatusInternalServerError, "Error decoding request parameters")
		return	
	}

	hashedPassword, err := auth.HashPassword(reqBody.Password)
	if err != nil {
		handleErrorResponse(w, http.StatusInternalServerError, "Error hashing password")
		return
	}
	createUserParams := database.CreateUserParams {
		Email:		reqBody.Email,
		HashedPassword:	hashedPassword,
	}

	user, err := cfg.db.CreateUser(context.Background(), createUserParams)
	if err != nil {
		handleErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Error creating user: %v", err))
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
		log.Fatalf("Error: failed to write response for hits")
	}
}

func main() {
	const filepathRoot = "."
	const port = "8080"


	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}

	envPlatform := os.Getenv("PLATFORM")
	if envPlatform == "" {
		log.Fatal("PLATFORM must be set")
	}

	envMaxChirpLength, err := strconv.Atoi(os.Getenv("MAX_CHIRP_LENGTH"))
	if err != nil {
		log.Fatalf("Fatal error occured converting a string if you can believe it: %v", err)
	}
	if envMaxChirpLength == 0 {
		log.Fatal("MAX_CHIRP_LENGTH must be set")
	}

	secret := os.Getenv("SECRET")
	if secret == "" {
		log.Fatal("SECRET must be set")
	}



	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Fatal error occured during connection to database: %v", err)
	}
	dbQueries := database.New(db)

	cfg := &apiConfig{
		fileserverHits: atomic.Int32{},
		db:		dbQueries,
		platform:	envPlatform,
		maxChirpLength: envMaxChirpLength,
		secret:		secret,
	}

	serveMux := http.NewServeMux()
	server := &http.Server{
		Addr:		":" + port,
		Handler: 	serveMux,
	}
	


	fileHandler := http.FileServer(http.Dir(filepathRoot))
	serveMux.HandleFunc("GET /api/healthz", handlerFunc)
	serveMux.Handle("/app/", http.StripPrefix("/app", cfg.middlewareMetricsInc(fileHandler)))
	serveMux.HandleFunc("GET /admin/metrics", cfg.handlerHits)
	serveMux.HandleFunc("POST /api/users", cfg.handlerCreateUser)
	serveMux.HandleFunc("POST /api/login", cfg.handlerLogin)
	serveMux.HandleFunc("POST /admin/reset", cfg.handlerReset)
	serveMux.HandleFunc("POST /api/chirps", cfg.handlerPostChirp)
	serveMux.HandleFunc("GET /api/chirps", cfg.handlerGetChirps)
	serveMux.HandleFunc("GET /api/chirps/{chirpID}", cfg.handlerGetChirp)
	serveMux.HandleFunc("POST /api/refresh", cfg.handlerRefresh)
	serveMux.HandleFunc("POST /api/revoke", cfg.handlerRevoke)
	serveMux.HandleFunc("PUT /api/users", cfg.handlerUpdateUser)
	server.ListenAndServe()
}
