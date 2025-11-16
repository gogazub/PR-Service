package userhandler

import (
	httperror "PRService/internal/adapters/http/error"
	userhttp "PRService/internal/adapters/http/user"
	"PRService/internal/domain/user"
	user_usecase "PRService/internal/usecase/user"
	"encoding/json"
	"errors"
	"net/http"
)

type SetIsActiveRequestDTO struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

type SetIsActiveResponseDTO struct {
	User userhttp.UserDTO `json:"user"`
}

func (h *Handler) SetIsActive(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != http.MethodPost {
		httperror.WriteBadRequest(w, "bad method")
		return
	}

	// 1. Parse JSON
	var req SetIsActiveRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httperror.WriteErrorResponse(w, http.StatusBadRequest, httperror.ErrorCodeBadRequest, "invalid json")
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
			httperror.WriteErrorResponse(w, http.StatusNotFound, httperror.ErrorCodeNotFound, err.Error())
		default:
			h.logger.Error("UpdateActive failed", "error", err)
			httperror.WriteErrorResponse(w, http.StatusInternalServerError, httperror.ErrorCodeInternal, "internal error")
		}
		return
	}

	// 4. Success response
	resp := userhttp.UserToDTO(u)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.Error("write response failed", "error", err)
	}
}
