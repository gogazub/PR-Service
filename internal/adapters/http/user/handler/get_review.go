package userhandler

import (
	httperror "PRService/internal/adapters/http/error"
	userhttp "PRService/internal/adapters/http/user"
	"PRService/internal/domain/pullrequest"
	"PRService/internal/domain/user"
	"encoding/json"
	"errors"
	"net/http"
)

type GetUserReviewQueryDTO struct {
	UserID string `query:"user_id"`
}

type GetUserReviewResponseDTO struct {
	UserID       string                         `json:"user_id"`
	PullRequests []userhttp.PullRequestShortDTO `json:"pull_requests"`
}

func (h *Handler) GetReview(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != http.MethodGet {
		httperror.WriteBadRequest(w, "bad method")
		return
	}

	// 1. Parse query param
	var req GetUserReviewQueryDTO
	req.UserID = r.URL.Query().Get("user_id")

	if req.UserID == "" {
		httperror.WriteErrorResponse(w, http.StatusBadRequest, httperror.ErrorCodeBadRequest, "missing user_id")
		return
	}

	ctx := r.Context()
	// 2. Domain command
	cmd := req.UserID

	// 3. Call service
	prs, err := h.Services.PullRequest.ListByUserID(ctx, user.ID(cmd))
	if err != nil {
		switch {
		case errors.Is(err, user.ErrUserNotFound):
			httperror.WriteErrorResponse(w, http.StatusNotFound, httperror.ErrorCodeNotFound, "user not found")
		default:
			h.logger.Error("GetReview failed", "error", err)
			httperror.WriteErrorResponse(w, http.StatusInternalServerError, httperror.ErrorCodeInternal, "internal error")
		}
		return
	}

	// 4. Success response
	resp := PullRequestsToReviewResponseDTO(req.UserID, prs)

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.Error("write response failed", "error", err)
	}
}

func PullRequestsToReviewResponseDTO(userID string, prs []*pullrequest.PullRequest) *GetUserReviewResponseDTO {
	resp := new(GetUserReviewResponseDTO)
	resp.UserID = userID
	for _, pr := range prs {
		resp.PullRequests = append(resp.PullRequests, *userhttp.PullRequestToShortDTO(pr))
	}
	return resp
}
