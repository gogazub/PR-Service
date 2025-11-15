package pullrequest_usecase

import (
	"PRService/internal/domain/pullrequest"
	"PRService/internal/domain/user"
)

type CreatePRCommand struct {
	Name      string
	Author    user.ID
	Reviewers []user.ID
}

type UpdateStatusCommand struct {
	PullRequestID pullrequest.ID
	Status        pullrequest.Status
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
