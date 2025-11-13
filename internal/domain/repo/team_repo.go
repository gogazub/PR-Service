package repo

import (
	"PRService/internal/domain"

	"context"
)

// TeamRepo - contract for working with teams in DB
type TeamRepo interface {

	// GetByName returns team by name
	GetByName(ctx context.Context, name string) (*domain.Team, error)

	// Create saves new team
	Create(ctx context.Context, team *domain.Team) error
}
