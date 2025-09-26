package auth

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/middleware"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/usecase/auth"
)

type Handler struct {
	authUC *auth.UseCase
}

func NewHandler(authUC *auth.UseCase) *Handler {
	return &Handler{authUC: authUC}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response, err := h.authUC.Register(req)
	if err != nil {
		log.Printf("Registration error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.sendJSONResponse(w, response)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response, err := h.authUC.Login(req)
	if err != nil {
		log.Printf("Login error: %v", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	h.sendJSONResponse(w, response)
}

func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	user, err := h.authUC.GetUserByID(userID)
	if err != nil {
		log.Printf("Get profile error: %v", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	h.sendJSONResponse(w, user)
}

func (h *Handler) sendJSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
