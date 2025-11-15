package teamrepo

import (
	"PRService/internal/domain/team"
	"context"
	"database/sql"
	"os/user"
)

type Repo struct {
	db *sql.DB
}

func NewTeamRepo(db *sql.DB) *Repo {
	return &Repo{db: db}
}

// Save saves new team
func (r *Repo) Save(ctx context.Context, team *team.Team) error {
	// TODO: implement
	return nil
}

// GetByName returns team by name
func (r *Repo) GetByName(ctx context.Context, name string) (*team.Team, error) {
	// TODO: implement
	return nil, nil
}

// GetByID returns team by id
func (r *Repo) GetByID(ctx context.Context, id string) (*team.Team, error) {
	// TODO: implement
	return nil, nil
}

// GetActiveUsersInTeam returns all active users in team
func (r *Repo) GetActiveUsersInTeam(ctx context.Context, teamID string) ([]*user.User, error) {
	// TODO: implement
	return nil, nil
}

// Update modifies an existing team
func (r *Repo) Update(ctx context.Context, team *team.Team) error {
	// TODO: implement
	return nil
}

// DeleteByID deletes team by id
func (r *Repo) DeleteByID(ctx context.Context, teamID string) error {
	// TODO: implement
	return nil
}

// DeleteByName deletes team by name
func (r *Repo) DeleteByName(ctx context.Context, name string) error {
	// TODO: implement
	return nil
}
