package httperror

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
