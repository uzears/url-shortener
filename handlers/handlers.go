package handlers

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/uzears/url-shortener/models"
)

type URLRequest struct {
	LongURL string `json:"long_url"`
}

type URLResponse struct {
	ShortURL string `json:"short_url"`
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateShortID(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func ShortenURL(w http.ResponseWriter, r *http.Request) {
	var req URLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	shortID := generateShortID(6)
	models.URLStore[shortID] = req.LongURL

	resp := URLResponse{ShortURL: "http://localhost:8080/" + shortID}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func RedirectURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortID := vars["shortID"]

	longURL, exists := models.URLStore[shortID]
	if !exists {
		http.Error(w, "Short URL not found", http.StatusNotFound)
		return
	}

	// 🔁 Perform redirect
	http.Redirect(w, r, longURL, http.StatusFound)
}
