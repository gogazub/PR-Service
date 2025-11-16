package user

import (
	"context"
)

// UserRepo - contract for working with users in DB
type Repo interface {

	// Save saves new user
	Save(ctx context.Context, user *User) error

	// GetByID returns user by ID
	GetByID(ctx context.Context, id ID) (*User, error)

	// GetByIDs returns users by IDs
	GetByIDs(ctx context.Context, id []ID) ([]*User, error)

	// UpdateActive updates user's active status
	UpdateActive(ctx context.Context, id ID, isActive bool) (*User, error)

	// Delete deletes user by id
	DeleteByID(ctx context.Context, id ID) error
}
