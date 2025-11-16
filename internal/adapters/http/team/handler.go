package teamhttp

import (
	"PRService/internal/app"
	"net/http"

	"go.uber.org/zap"
)

type Handler struct {
	*app.Services
	logger *zap.SugaredLogger
}

// NewHandler returns new Handler.
func NewHandler(app *app.Services, logger *zap.SugaredLogger) *Handler {
	return &Handler{app, logger}
}

// POST /team/add
func (h *Handler) AddTeam(w http.ResponseWriter, r *http.Request) {
}

// GET /team/get?team_name=...
func (h *Handler) GetTeam(w http.ResponseWriter, r *http.Request) {

}
