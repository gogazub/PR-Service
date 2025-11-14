package repo

import (
	"PRService/internal/domain"

	"context"
)

// TeamRepo - contract for working with teams in DB
type TeamRepo interface {

	// Create saves new team
	Create(ctx context.Context, team *domain.Team) error

	// Так как имена команд уникальны - можем искать по ним
	// Но для полноты API добавим оба варианта: и по ID, и по name.
	// GetByName returns team by name
	GetByName(ctx context.Context, name string) (*domain.Team, error)

	// GetByID returns team by id
	GetByID(ctx context.Context, id string) (*domain.Team, error)

	// GetActiveUsersInTeam returns all active users in team
	GetActiveUsersInTeam(ctx context.Context, teamID string) ([]*domain.User, error)

	// Update updates team in DB
	Update(ctx context.Context, team *domain.Team)

	// DeleteByID deletes team by id
	DeleteByID(ctx context.Context, teamID string) error
}
