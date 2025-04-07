package handler

import (
	"errors"
	"fmt"
	"net/http"
	"url-shortener/proto"

	"github.com/gorilla/mux"
	"url-shortener/internal/service"
	"url-shortener/internal/storage"
)

// Handler обрабатывает HTTP-запросы для сервиса сокращения ссылок
type Handler struct {
	service *service.Service
}

// NewHandler создаёт экземпляр обработчика с переданным сервисом
func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

// CreateURL обрабатывает POST-запрос для создания короткой ссылки
func (h *Handler) CreateURL(w http.ResponseWriter, r *http.Request) {
	originalURL := r.FormValue("url")
	if originalURL == "" {
		http.Error(w, "Отсутствует параметр url", http.StatusBadRequest)
		return
	}

	resp, err := h.service.CreateURL(r.Context(), &proto.CreateURLRequest{
		OriginalUrl: originalURL,
	})
	if err != nil {
		http.Error(w, "Не удалось создать короткую ссылку: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if resp.Error != "" {
		http.Error(w, resp.Error, http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, resp.ShortUrl)
}

// GetURL обрабатывает GET-запрос для получения оригинальной ссылки по короткой
func (h *Handler) GetURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortURL := vars["shortURL"]

	resp, err := h.service.GetURL(r.Context(), &proto.GetURLRequest{
		ShortUrl: shortURL,
	})
	if err != nil {
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}
	if resp.Error != "" {
		if errors.Is(err, storage.ErrNotFound) {
			http.Error(w, "Ссылка не найдена", http.StatusNotFound)
		} else {
			http.Error(w, resp.Error, http.StatusInternalServerError)
		}
		return
	}

	fmt.Fprintln(w, resp.OriginalUrl)
}

// SetupRoutes настраивает маршруты API с использованием маршрутизатора gorilla/mux
func (h *Handler) SetupRoutes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", h.CreateURL).Methods("POST")
	r.HandleFunc("/{shortURL}", h.GetURL).Methods("GET")
	return r
}
