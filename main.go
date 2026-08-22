package main

import (
	"log"
	"net/http"
)

func main() {
	apiCfg := &apiConfig{}
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/healthz", apiCfg.handlerHealth)
	mux.HandleFunc("POST /api/validate_chirp", apiCfg.handlerValidateChirp)
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
