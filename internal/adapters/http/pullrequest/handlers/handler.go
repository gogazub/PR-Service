package pullreqhandler

import (
	pullrequesthttp "PRService/internal/adapters/http/pullrequest"
	"PRService/internal/app"
	"PRService/internal/domain/pullrequest"
	"PRService/pkg/logger"
)

type Handler struct {
	*app.Services
	logger *logger.Logger
}

// NewHandler returns new Handler.
func NewHandler(app *app.Services, logger *logger.Logger) *Handler {
	return &Handler{app, logger}
}

func PRToDTO(pr *pullrequest.PullRequest) pullrequesthttp.PullRequestDTO {

	revs := make([]string, len(pr.Reviewers))
	for i, r := range pr.Reviewers {
		revs[i] = string(r)
	}

	return pullrequesthttp.PullRequestDTO{
		PullRequestID:     string(pr.PullRequestID),
		PullRequestName:   pr.Name,
		AuthorID:          string(pr.Author),
		Status:            pullrequest.StatusToString(pr.Status),
		AssignedReviewers: revs,
	}
}
