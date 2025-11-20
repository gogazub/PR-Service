package pullreqhandler

import (
	httperror "PRService/internal/adapters/http/error"
	pullrequesthttp "PRService/internal/adapters/http/pullrequest"
	"PRService/internal/domain/pullrequest"
	"PRService/internal/domain/user"
	pullrequest_usecase "PRService/internal/usecase/pullrequest"
	"encoding/json"
	"errors"
	"net/http"
)

// POST /pullRequest/reassign
type ReassignReviewerRequest struct {
	PullRequestID string `json:"pull_request_id"`
	OldUserID     string `json:"old_user_id"`
}

type ReassignReviewerResponse struct {
	PR         pullrequesthttp.PullRequestDTO `json:"pr"`
	ReplacedBy string                         `json:"replaced_by"`
}

func (h *Handler) ReassignReviewer(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != http.MethodPost {
		httperror.WriteBadRequest(w, "bad method")
		return
	}

	var req ReassignReviewerRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httperror.WriteBadRequest(w, "bad json")
		return
	}

	cmd := pullrequest_usecase.ReassignReviewerCommand{
		PullRequestID: pullrequest.ID(req.PullRequestID),
		OldReviewerID: user.ID(req.OldUserID),
	}

	ctx := r.Context()

	pr, newReviewerID, err := h.Services.ReassignReviewer(ctx, cmd)
	if err != nil {
		if errors.Is(err, pullrequest.ErrNoAuthor) {
			httperror.WriteErrorResponse(w, http.StatusNotFound, httperror.ErrorCodeBadRequest, "no author")
			return
		}
		if errors.Is(err, pullrequest.ErrNoCandidate) {
			httperror.WriteErrorResponse(w, http.StatusConflict, httperror.ErrorCodeNoCandidate, "no candidate")
			return
		}
		if errors.Is(err, pullrequest.ErrPullRequestExists) {
			httperror.WriteErrorResponse(w, http.StatusNotFound, httperror.ErrorCodeNotFound, "pr not found")
			return
		}
		if errors.Is(err, pullrequest.ErrMerged) {
			httperror.WriteErrorResponse(w, http.StatusConflict, httperror.ErrorCodePRMerged, "pr merged")
			return
		}
		if errors.Is(err, pullrequest.ErrReviewerNotAssigned) {
			httperror.WriteErrorResponse(
				w,
				http.StatusConflict,
				httperror.ErrorCodeNotAssigned,
				err.Error(),
			)
			return
		}
		if errors.Is(err, pullrequest.ErrPullRequestNotFound) {
			httperror.WriteErrorResponse(
				w,
				http.StatusNotFound,
				httperror.ErrorCodeNotFound,
				err.Error(),
			)
			return
		}

		h.logger.Error("reassign failed", "error", err)
		httperror.WriteErrorResponse(w, http.StatusInternalServerError, httperror.ErrorCodeInternal, "internal error")
		return
	}

	resp := ReassignReviewerResponse{
		PR:         PRToDTO(pr),
		ReplacedBy: string(newReviewerID),
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.Error("reassign failed: JSON encoding error", "error", err)
	}
}
