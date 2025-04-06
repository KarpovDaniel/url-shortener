package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"url-shortener/internal/storage"

	"github.com/gorilla/mux"
	"url-shortener/internal/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateURL(w http.ResponseWriter, r *http.Request) {
	var req struct {
		OriginalURL string `json:"original_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	shortURL, err := h.service.Create(req.OriginalURL)
	if err != nil {
		http.Error(w, "Failed to create short URL: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"short_url": shortURL})
}

func (h *Handler) GetURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortURL := vars["shortURL"]
	originalURL, err := h.service.Get(shortURL)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			http.Error(w, "URL not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"original_url": originalURL})
}

func (h *Handler) SetupRoutes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", h.CreateURL).Methods("POST")
	r.HandleFunc("/{shortURL}", h.GetURL).Methods("GET")
	return r
}
