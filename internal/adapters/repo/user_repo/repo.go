package userrepo

import (
	"PRService/internal/domain"
	"context"
	"database/sql"
)

type UserRepo struct {
	db *sql.DB
}

// NewUserRepo returns new UserRepo.
func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) GetByID(ctx context.Context, id domain.UserID) (*domain.User, error) {
	return nil, nil
}

func (r *UserRepo) Create(ctx context.Context, user *domain.User) error {
	return nil
}

func (r *UserRepo) GetByTeamID(ctx context.Context, teamID string) ([]*domain.User, error) {
	return nil, nil
}

func (r *UserRepo) UpdateActive(ctx context.Context, id domain.UserID, isActive bool) error {
	return nil
}
