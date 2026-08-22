package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type errorRes struct {
		Error string `json:"error"`
	}
	respondWithJson(w, code, errorRes{
		msg,
	})
}

func respondWithJson(w http.ResponseWriter, code int, payload any) {
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("error handling json data: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}

func cleanProfanity(body string) string {
	profaneWords := []string{
		"kerfuffle",
		"sharbert",
		"fornax",
	}
	words := strings.Split(body, " ")

	for i, word := range words {
		lowered := strings.ToLower(word)
		if slices.Contains(profaneWords, lowered) {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}
