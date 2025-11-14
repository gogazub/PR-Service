package teamrepo

import (
	"PRService/internal/domain"
	"context"
	"database/sql"
)

type TeamRepo struct {
	db *sql.DB
}

func NewTeamRepo(db *sql.DB) *TeamRepo {
	return &TeamRepo{db: db}
}

// Create saves new team
func (r *TeamRepo) Create(ctx context.Context, team *domain.Team) error {
	// TODO: implement
	return nil
}

// GetByName returns team by name
func (r *TeamRepo) GetByName(ctx context.Context, name string) (*domain.Team, error) {
	// TODO: implement
	return nil, nil
}

// GetByID returns team by id
func (r *TeamRepo) GetByID(ctx context.Context, id string) (*domain.Team, error) {
	// TODO: implement
	return nil, nil
}

// GetActiveUsersInTeam returns all active users in team
func (r *TeamRepo) GetActiveUsersInTeam(ctx context.Context, teamID string) ([]*domain.User, error) {
	// TODO: implement
	return nil, nil
}

// Update modifies an existing team
func (r *TeamRepo) Update(ctx context.Context, team *domain.Team) {
	// TODO: implement
}

// DeleteByID deletes team by id
func (r *TeamRepo) DeleteByID(ctx context.Context, teamID string) error {
	// TODO: implement
	return nil
}
