package team

import (
	"PRService/internal/domain/user"
	"context"
)

// TeamRepo - contract for working with teams in DB
type Repo interface {

	// Save saves new team
	Save(ctx context.Context, team *Team) error

	// Так как имена команд уникальны - можем искать по ним
	// Но для полноты API добавим оба варианта: и по ID, и по name.
	// GetByName returns team by name
	GetByName(ctx context.Context, name string) (*Team, error)

	// GetByID returns team by id
	GetByID(ctx context.Context, id ID) (*Team, error)

	// GetActiveUsersInTeam returns all active users in team
	GetActiveUsersInTeam(ctx context.Context, teamID ID) ([]*user.User, error)

	// Update updates team in DB
	Update(ctx context.Context, team *Team) error

	// DeleteByID deletes team by id
	DeleteByID(ctx context.Context, id ID) error

	// DeleteByName deletes team by name
	DeleteByName(ctx context.Context, name string) error
}
