package pullrequesthttp

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

// POST /pullRequest/create
func (h *Handler) CreatePullRequest(w http.ResponseWriter, r *http.Request) {
	
}

// POST /pullRequest/merge
func (h *Handler) MergePullRequest(w http.ResponseWriter, r *http.Request) {
	
}

// POST /pullRequest/reassign
func (h *Handler) ReassignReviewer(w http.ResponseWriter, r *http.Request) {
	
}
