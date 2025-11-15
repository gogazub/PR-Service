package team

import (
	"PRService/internal/domain/user"
	"context"
)

// TeamRepo - contract for working with teams in DB
type Repo interface {

	// Save saves new team
	Save(ctx context.Context, team *Team) error

	// GetByName returns team by name
	GetByName(ctx context.Context, name string) (*Team, error)

	// GetActiveUsersInTeam returns all active users in team
	GetActiveUsersInTeam(ctx context.Context, name string) ([]*user.User, error)

	// Update updates team in DB
	Update(ctx context.Context, team *Team) error

	// DeleteByName deletes team by name
	DeleteByName(ctx context.Context, name string) error
}
