package pullreqhandler

import (
	httperror "PRService/internal/adapters/http/error"
	pullrequesthttp "PRService/internal/adapters/http/pullrequest"
	"PRService/internal/domain/pullrequest"
	pullrequest_usecase "PRService/internal/usecase/pullrequest"
	"encoding/json"
	"errors"
	"net/http"
)

// POST /pullRequest/merge
type MergePullRequestRequest struct {
	PullRequestID string `json:"pull_request_id"`
}

type MergePullRequestResponse struct {
	PR pullrequesthttp.PullRequestDTO `json:"pr"`
}

func (h *Handler) MergePullRequest(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != http.MethodPost {
		httperror.WriteBadRequest(w, "bad method")
		return
	}

	var req MergePullRequestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httperror.WriteBadRequest(w, "invalid json")
		return
	}
	if req.PullRequestID == "" {
		httperror.WriteBadRequest(w, "pull_request_id is required")
		return
	}

	ctx := r.Context()
	cmd := pullrequest_usecase.UpdateStatusCommand{
		PullRequestID: pullrequest.ID(req.PullRequestID),
	}

	pr, err := h.Services.PullRequest.UpdateStatus(ctx, cmd)
	if err != nil {
		if errors.Is(err, pullrequest.ErrPullRequestNotFound) {
			httperror.WriteErrorResponse(
				w,
				http.StatusNotFound,
				httperror.ErrorCodeNotFound,
				"pull request not found",
			)
			return
		}

		h.logger.Error("merge PR failed: update status failed", "error", err)
		httperror.WriteErrorResponse(
			w,
			http.StatusInternalServerError,
			httperror.ErrorCodeInternal,
			"internal error",
		)
		return
	}

	resp := MergePullRequestResponse{
		PR: PRToDTO(pr),
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.Error("merge PR failed: JSON encoding error", "error", err)
	}
}
