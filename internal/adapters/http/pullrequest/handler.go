package pullrequesthttp

import (
	"PRService/internal/app"
	"net/http"
)

type Handler struct {
	app.Services
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
