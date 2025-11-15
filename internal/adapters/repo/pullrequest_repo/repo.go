package pullrequestrepo

import (
	"PRService/internal/domain/pullrequest"
	"PRService/internal/domain/user"
	"context"
	"database/sql"
)

type Repo struct {
	db *sql.DB
}

func NewPullRequestRepo(db *sql.DB) *Repo {
	return &Repo{db: db}
}

// Save saves new pull request
func (r *Repo) Save(ctx context.Context, pr *pullrequest.PullRequest) error {
	// TODO: implement
	return nil
}

// GetByID returns pull request by ID
func (r *Repo) GetByID(ctx context.Context, prID pullrequest.ID) (*pullrequest.PullRequest, error) {
	// TODO: implement
	return nil, nil
}

// UpdateStatus changes pull request status
func (r *Repo) UpdateStatus(ctx context.Context, prID pullrequest.ID, status pullrequest.Status) error {
	// TODO: implement
	return nil
}

// AssignReviewers sets reviewers for pull request
func (r *Repo) AssignReviewers(ctx context.Context, prID pullrequest.ID, reviewers []user.ID) error {
	// TODO: implement
	return nil
}

// ReassignReviewers replaces one reviewer with another
func (r *Repo) ReassignReviewers(ctx context.Context, prID pullrequest.ID, oldReviewerID, newReviewerID user.ID) error {
	// TODO: implement
	return nil
}

// ListByUserID returns all user's pull requests
func (r *Repo) ListByUserID(ctx context.Context, userID user.ID) ([]*pullrequest.PullRequest, error) {
	// TODO: implement
	return nil, nil
}
