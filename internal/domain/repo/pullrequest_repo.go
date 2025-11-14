package repo

import (
	"PRService/internal/domain"

	"context"
)

// PullRequestRepo - contract for working with users in DB
type PullRequestRepo interface {

	// Create saves new pull request
	Create(ctx context.Context, pr *domain.PullRequest) error

	// GetByID returns pull request by ID
	GetByID(ctx context.Context, prID domain.PullRequestID) (*domain.PullRequest, error)

	// UpdateStatus changes pull request status
	UpdateStatus(ctx context.Context, status domain.PullRequestStatus) error

	// AssignReviewers sets reviewers for pull request
	AssignReviewers(ctx context.Context, prID domain.PullRequestID, reviewers []domain.UserID) error

	// ReassignReviewers replaces one reviewer with another
	ReassignReviewers(ctx context.Context, prID domain.PullRequestID, oldReviewerID, newReviewerID domain.UserID) error

	// ListByUserID returns all user's pull requests
	ListByUserID(ctx context.Context, userID domain.UserID) ([]*domain.PullRequest, error)
}
