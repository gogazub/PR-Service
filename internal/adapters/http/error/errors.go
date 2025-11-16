package httperror

import (
	"encoding/json"
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
	// TODO: process error
	json.NewEncoder(w).Encode(ErrorResponseDTO{
		Error: ErrorDTO{
			Code:    ErrorCodeNotFound,
			Message: msg,
		},
	})
}

func WriteErrorResponse(w http.ResponseWriter, status int, code ErrorCode, msg string) {
	w.WriteHeader(status)
	// process error
	json.NewEncoder(w).Encode(ErrorResponseDTO{
		Error: ErrorDTO{
			Code:    code,
			Message: msg,
		},
	})
}