package pullrequestusecase

import (
	"PRService/internal/domain/pullrequest"
	"PRService/internal/domain/user"
)

type CreatePRCommand struct {
	ID        string
	Name      string
	Author    user.ID
}

type UpdateStatusCommand struct {
	PullRequestID pullrequest.ID
}

type AssignReviewersCommand struct {
	PullRequestID pullrequest.ID
	Reviewers     []user.ID
}

type ReassignReviewerCommand struct {
	PullRequestID pullrequest.ID
	OldReviewerID user.ID
	NewReviewerID user.ID
}
