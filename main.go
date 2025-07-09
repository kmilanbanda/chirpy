package main

import (
	"net/http"
	"log"
	"sync/atomic"
	"fmt"
	"encoding/json"
)

type Chirp struct {
	body	string `json:"body"`
}

type apiConfig struct  {
	fileserverHits atomic.Int32
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
	resp := struct{ Valid bool `json:"valid"` }{true}
	dat, _ := json.Marshal(resp)
	w.Write(dat)
	return
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

func (cfg *apiConfig) handlerHitsReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	_, err := w.Write([]byte(""))
	if err != nil {
		log.Fatalf("Failed to write response for hits reset")
	}
}

func main() {
	serveMux := http.NewServeMux()
	server := &http.Server{
		Addr:		":8080",
		Handler: 	serveMux,
	}
	cfg := &apiConfig{}


	fileHandler := http.FileServer(http.Dir("."))
	serveMux.HandleFunc("GET /api/healthz", handlerFunc)
	serveMux.Handle("/app/", http.StripPrefix("/app/", cfg.middlewareMetricsInc(fileHandler)))
	serveMux.HandleFunc("GET /admin/metrics", cfg.handlerHits)
	serveMux.HandleFunc("POST /admin/reset", cfg.handlerHitsReset)
	serveMux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)
	server.ListenAndServe()
}
