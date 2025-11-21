package httperror

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrorCode string

const (
	ErrorCodeTeamExists  ErrorCode = "TEAM_EXISTS"
	ErrorCodePRExists    ErrorCode = "PR_EXISTS"
	ErrorCodePRMerged    ErrorCode = "PR_MERGED"
	ErrorCodeNotAssigned ErrorCode = "NOT_ASSIGNED"
	ErrorCodeNoCandidate ErrorCode = "NO_CANDIDATE"
	ErrorCodeNotFound    ErrorCode = "NOT_FOUND"
	ErrorCodeBadRequest  ErrorCode = "IVALID_JSON"
	ErrorCodeInternal    ErrorCode = "INTERNAL"
)

type ErrorDTO struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

type ErrorResponseDTO struct {
	Error ErrorDTO `json:"error"`
}

// Так как в спецификации нет подходящей ошибки invalid json,
// то вынуждено проставим ErrorCodeNotFound и пояснения в message
func WriteBadRequest(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusBadRequest)
	err := json.NewEncoder(w).Encode(ErrorResponseDTO{
		Error: ErrorDTO{
			Code:    ErrorCodeNotFound,
			Message: msg,
		},
	})
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		fmt.Printf("json encoding error: %v", err)
	}
}

func WriteErrorResponse(w http.ResponseWriter, status int, code ErrorCode, msg string) {
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(ErrorResponseDTO{
		Error: ErrorDTO{
			Code:    code,
			Message: msg,
		},
	})
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		fmt.Printf("json encoding error: %v", err)
	}
}
