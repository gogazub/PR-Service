package pullrequestrepo

import (
	"PRService/internal/domain"
	"context"
	"database/sql"
)

type PullRequestRepo struct {
	db *sql.DB
}

func NewPullRequestRepo(db *sql.DB) *PullRequestRepo {
	return &PullRequestRepo{db: db}
}

// Create saves new pull request
func (r *PullRequestRepo) Create(ctx context.Context, pr *domain.PullRequest) error {
	// TODO: implement
	return nil
}

// GetByID returns pull request by ID
func (r *PullRequestRepo) GetByID(ctx context.Context, prID domain.PullRequestID) (*domain.PullRequest, error) {
	// TODO: implement
	return nil, nil
}

// UpdateStatus changes pull request status
func (r *PullRequestRepo) UpdateStatus(ctx context.Context, status domain.PullRequestStatus) error {
	// TODO: implement
	return nil
}

// AssignReviewers sets reviewers for pull request
func (r *PullRequestRepo) AssignReviewers(ctx context.Context, prID domain.PullRequestID, reviewers []domain.UserID) error {
	// TODO: implement
	return nil
}

// ReassignReviewers replaces one reviewer with another
func (r *PullRequestRepo) ReassignReviewers(ctx context.Context, prID domain.PullRequestID, oldReviewerID, newReviewerID domain.UserID) error {
	// TODO: implement
	return nil
}

// ListByUserID returns all user's pull requests
func (r *PullRequestRepo) ListByUserID(ctx context.Context, userID domain.UserID) ([]*domain.PullRequest, error) {
	// TODO: implement
	return nil, nil
}
