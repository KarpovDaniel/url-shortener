package handler

import (
	"errors"
	"fmt"
	"net/http"
	"url-shortener/internal/storage"

	"github.com/gorilla/mux"
	"url-shortener/internal/service"
)

// Handler обрабатывает HTTP-запросы для сервиса сокращения URL
type Handler struct {
	service *service.Service
}

// NewHandler создаёт экземпляр обработчика с зависимостью от сервиса
func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

// CreateURL обрабатывает POST-запрос для создания сокращённой ссылки
func (h *Handler) CreateURL(w http.ResponseWriter, r *http.Request) {
	originalURL := r.FormValue("url")
	if originalURL == "" {
		http.Error(w, "Missing original_url", http.StatusBadRequest)
		return
	}
	shortURL, err := h.service.Create(originalURL)
	if err != nil {
		http.Error(w, "Failed to create short URL: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, shortURL) //nolint:errcheck
}

// GetURL обрабатывает GET-запрос для получения оригинальной ссылки по короткому ключу
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
	fmt.Fprint(w, originalURL) //nolint:errcheck
}

// SetupRoutes настраивает маршруты API с использованием роутера gorilla/mux
func (h *Handler) SetupRoutes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", h.CreateURL).Methods("POST")
	r.HandleFunc("/{shortURL}", h.GetURL).Methods("GET")
	return r
}
