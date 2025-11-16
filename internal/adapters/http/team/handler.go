package teamhttp

import (
	"PRService/internal/app"
	"net/http"
)

type Handler struct {
	app.Services
}

// POST /team/add
func (h *Handler) AddTeam(w http.ResponseWriter, r *http.Request) {
}

// GET /team/get?team_name=...
func (h *Handler) GetTeam(w http.ResponseWriter, r *http.Request) {

}
