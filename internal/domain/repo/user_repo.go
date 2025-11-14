package repo

import (
	"PRService/internal/domain"

	"context"
)

// UserRepo - contract for working with users in DB
type UserRepo interface {

	// Create saves new user
	Create(ctx context.Context, user *domain.User) error

	// GetByID returns user by ID
	GetByID(ctx context.Context, id domain.UserID) (*domain.User, error)

	// UpdateActive updates user's active status
	UpdateActive(ctx context.Context, id domain.UserID, isActive bool) error

	// Delete deletes user by id
	DeleteByID(ctx context.Context, id domain.UserID) error
}
