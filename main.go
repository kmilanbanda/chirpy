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

	"github.com/kmilanbanda/chirpy/internal/database"
	"github.com/joho/godotenv"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type apiConfig struct  {
	fileserverHits 	atomic.Int32
	db		*database.Queries
	platform	string
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

func main() {
	const filepathRoot = "."
	const port = "8080"


	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	pform := os.Getenv("PLATFORM")
	if pform == "" {
		log.Fatal("PLATFORM must be set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Fatal error occured during connection to database: %w", err)
	}
	dbQueries := database.New(db)

	cfg := &apiConfig{
		fileserverHits: atomic.Int32{},
		db:		dbQueries,
		platform:	pform,
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
	serveMux.HandleFunc("POST /admin/reset", cfg.handlerReset)
	serveMux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)
	server.ListenAndServe()
}
