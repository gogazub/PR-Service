package userhttp

import (
	"PRService/internal/app"
	"net/http"
)

type Handler struct {
	app.Services
}

// POST /users/setIsActive
func (h *Handler) SetIsActive(w http.ResponseWriter, r *http.Request) {
	
}

// GET /users/getReview?user_id=...
func (h *Handler) GetReview(w http.ResponseWriter, r *http.Request) {
	
}