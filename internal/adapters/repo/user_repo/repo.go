package userrepo

import (
	"PRService/internal/domain"
	"context"
	"database/sql"
)

type UserRepo struct {
	db *sql.DB
}

// Create saves new user
func (r *UserRepo) Create(ctx context.Context, user *domain.User) error {
	// TODO: implement
	return nil
}

// GetByID returns user by ID
func (r *UserRepo) GetByID(ctx context.Context, id domain.UserID) (*domain.User, error) {
	// TODO: implement
	return nil, nil
}

// UpdateActive updates user's active status
func (r *UserRepo) UpdateActive(ctx context.Context, id domain.UserID, isActive bool) error {
	// TODO: implement
	return nil
}

// DeleteByID deletes user by id
func (r *UserRepo) DeleteByID(ctx context.Context, id domain.UserID) error {
	// TODO: implement
	return nil
}
