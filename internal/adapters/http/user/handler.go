package userhttp

import (
	httperror "PRService/internal/adapters/http/error"
	"PRService/internal/app"
	"PRService/internal/domain/user"
	user_usecase "PRService/internal/usecase/user"
	"encoding/json"
	"errors"
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

// POST /users/setIsActive
func (h *Handler) SetIsActive(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != http.MethodPost {
		writeBadRequest(w, "bad method")
	}

	// 1. Parse JSON
	var req SetIsActiveRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, httperror.ErrorCodeBadRequest, "invalid json")
		return
	}

	ctx := r.Context()

	// 2. Domain command
	cmd := user_usecase.UpdateActiveCommand{
		UserID:   user.ID(req.UserID),
		IsActive: req.IsActive,
	}

	// 3. Call service
	u, err := h.Services.User.UpdateActive(ctx, cmd)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrUserNotFound):
			writeErrorResponse(w, http.StatusNotFound, httperror.ErrorCodeNotFound, err.Error())
		default:
			h.logger.Error("UpdateActive failed", "error", err)
			writeErrorResponse(w, http.StatusInternalServerError, httperror.ErrorCodeInternal, "internal error")
		}
		return
	}

	// 4. Success response
	resp := UserToDTO(u)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.Error("write response failed", "error", err)
	}
}

// GET /users/getReview?user_id=...
func (h *Handler) GetReview(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != http.MethodGet {
		writeErrorResponse(w, http.StatusBadRequest, httperror.ErrorCodeBadRequest, "bad method")
		return
	}

	// 1. Parse query param
	var req GetUserReviewQueryDTO
	req.UserID = r.URL.Query().Get("user_id")

	if req.UserID == "" {
		writeErrorResponse(w, http.StatusBadRequest, httperror.ErrorCodeBadRequest, "missing user_id")
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
			writeErrorResponse(w, http.StatusNotFound, httperror.ErrorCodeNotFound, err.Error())
		default:
			h.logger.Error("GetReview failed", "error", err)
			writeErrorResponse(w, http.StatusInternalServerError, httperror.ErrorCodeInternal, "internal error")
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


// Так как в спецификации нет подходящей ошибки invalid json,
// то вынуждено проставим ErrorCodeNotFound и пояснения в message
func writeBadRequest(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusBadRequest)
	// TODO: process error
	json.NewEncoder(w).Encode(httperror.ErrorResponseDTO{
		Error: httperror.ErrorDTO{
			Code:    httperror.ErrorCodeNotFound,
			Message: msg,
		},
	})
}

func writeErrorResponse(w http.ResponseWriter, status int, code httperror.ErrorCode, msg string) {
	w.WriteHeader(status)
	// process error
	json.NewEncoder(w).Encode(httperror.ErrorResponseDTO{
		Error: httperror.ErrorDTO{
			Code:    code,
			Message: msg,
		},
	})
}
