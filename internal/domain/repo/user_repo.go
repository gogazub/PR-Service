package repo

import (
	"PRService/internal/domain"

	"context"
)

// UserRepo - contract for working with users in DB
type UserRepo interface {

	// GetByID returns user by ID
	GetByID(ctx context.Context, id domain.UserID) (*domain.User, error)

	// Create saves new user
	Create(ctx context.Context, user *domain.User) error

	// GetByTeamID returns all active users in team
	GetByTeamID(ctx context.Context, teamID string) ([]*domain.User, error)

	// UpdateActive updates user's active status
	UpdateActive(ctx context.Context, id domain.UserID, isActive bool) error
}
