package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/tsuyoshi64/tsuyoshi-first-http-server/internal/database"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading environmental variables")
	}

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("DB_URL environment variale is not set")
	}

	platform := os.Getenv("PLATFORM")

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("error opening database link: %s", err)
	}
	defer db.Close()

	dbQueries := database.New(db)

	apiCfg := &apiConfig{
		DB:       dbQueries,
		platform: platform,
	}
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirpByID)

	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)

	mux.HandleFunc("GET /api/healthz", apiCfg.handlerHealth)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	fileServer := http.FileServer(http.Dir("."))

	mux.Handle(
		"/app/",
		apiCfg.middlewareMetricsInc(http.StripPrefix("/app", fileServer)),
	)

	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	log.Print("serving on port 8080...")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("couldn't start the server: %v", err)
	}
}
