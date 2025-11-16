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
	app.Services
	logger *zap.SugaredLogger
}

// POST /users/setIsActive
func (h *Handler) SetIsActive(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

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
	json.NewEncoder(w).Encode(resp)
}

// GET /users/getReview?user_id=...
func (h *Handler) GetReview(w http.ResponseWriter, r *http.Request) {

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
