package userrepo

import (
	"PRService/internal/domain/user"
	"context"
	"database/sql"
)

type Repo struct {
	db *sql.DB
}

// Create saves new user
func (r *Repo) Save(ctx context.Context, user *user.User) error {
	// TODO: implement
	return nil
}

// GetByID returns user by ID
func (r *Repo) GetByID(ctx context.Context, id user.ID) (*user.User, error) {
	// TODO: implement
	return nil, nil
}

// UpdateActive updates user's active status
func (r *Repo) UpdateActive(ctx context.Context, id user.ID, isActive bool) error {
	// TODO: implement
	return nil
}

// DeleteByID deletes user by id
func (r *Repo) DeleteByID(ctx context.Context, id user.ID) error {
	// TODO: implement
	return nil
}
