package pullrequest

import (
	"PRService/internal/domain/user"

	"context"
)

// PullRequestRepo - contract for working with users in DB
type Repo interface {

	// Save saves new pull request
	Save(ctx context.Context, pr *PullRequest) error

	// GetByID returns pull request by ID
	GetByID(ctx context.Context, prID ID) (*PullRequest, error)

	// UpdateStatus changes pull request status
	UpdateStatus(ctx context.Context, prID ID, status Status) error

	// AssignReviewers sets reviewers for pull request
	AssignReviewers(ctx context.Context, prID ID, reviewers []user.ID) error

	// ReassignReviewers replaces one reviewer with another
	ReassignReviewers(ctx context.Context, prID ID, oldReviewerID, newReviewerID user.ID) error

	// ListByUserID returns all user's pull requests
	ListByUserID(ctx context.Context, userID user.ID) ([]*PullRequest, error)
}
