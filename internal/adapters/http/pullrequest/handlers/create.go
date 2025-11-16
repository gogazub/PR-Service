package pullreqhandler

import (
	httperror "PRService/internal/adapters/http/error"
	pullrequesthttp "PRService/internal/adapters/http/pullrequest"
	"PRService/internal/domain/pullrequest"
	"PRService/internal/domain/team"
	"PRService/internal/domain/user"
	pullrequest_usecase "PRService/internal/usecase/pullrequest"
	"encoding/json"
	"errors"
	"net/http"
)

// POST /pullRequest/create

type CreatePullRequestRequest struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
}

type CreatePullRequestResponse struct {
	PR pullrequesthttp.PullRequestDTO `json:"pr"`
}

func (h *Handler) CreatePullRequest(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != http.MethodPost {
		httperror.WriteBadRequest(w, "bad method")
		return
	}

	var req CreatePullRequestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httperror.WriteBadRequest(w, "invalid json")
		return
	}

	ctx := r.Context()
	cmd := pullrequest_usecase.CreatePRCommand{
		ID:     req.PullRequestID,
		Name:   req.PullRequestName,
		Author: user.ID(req.AuthorID),
	}
	pr, err := h.Services.CreatePR(ctx, cmd)
	if err != nil {
		if errors.Is(err, pullrequest.ErrPullRequestExists) {
			httperror.WriteErrorResponse(w, http.StatusConflict, httperror.ErrorCodePRExists, "pr is already exists")
			return
		}
		if errors.Is(err, user.ErrUserNotFound) {
			httperror.WriteErrorResponse(w, http.StatusNotFound, httperror.ErrorCodeNotFound, "user not found")
			return
		}
		if errors.Is(err, team.ErrTeamNotFound) {
			httperror.WriteErrorResponse(w, http.StatusNotFound, httperror.ErrorCodeNotFound, "team not found")
			return
		}
		h.logger.Errorf("create pr: %w",err)
		httperror.WriteErrorResponse(w, http.StatusInternalServerError, httperror.ErrorCodeInternal, "internal error")
		return
	}
	
	
	resp := CreatePullRequestResponse{
		PR : PRToDTO(pr),
	}
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.Errorf("create pr: json encode: %w", err)
	}

}
