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

func (r *PullRequestRepo) GetByID(ctx context.Context, id domain.PullRequestID) (*domain.PullRequest, error) {
	return nil, nil
}

func (r *PullRequestRepo) Create(ctx context.Context, pr *domain.PullRequest) error {
	return nil
}

func (r *PullRequestRepo) Update(ctx context.Context, pr *domain.PullRequest) error {
	return nil
}

func (r *PullRequestRepo) GetByAuthorID(ctx context.Context, authorID domain.UserID) ([]*domain.PullRequest, error) {
	return nil, nil
}

func (r *PullRequestRepo) UpdateStatus(ctx context.Context, id domain.PullRequestID, status domain.PullRequestStatus) error {
	return nil
}

func (r *PullRequestRepo) GetByStatus(ctx context.Context, status domain.PullRequestStatus) ([]*domain.PullRequest, error) {
	return nil, nil
}
