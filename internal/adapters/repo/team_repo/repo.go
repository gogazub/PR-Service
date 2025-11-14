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

func (r *TeamRepo) GetByID(ctx context.Context, id domain.TeamID) (*domain.Team, error) {
	return nil, nil
}

func (r *TeamRepo) Create(ctx context.Context, team *domain.Team) error {
	return nil
}

func (r *TeamRepo) GetByName(ctx context.Context, name string) (*domain.Team, error) {
	return nil, nil
}

func (r *TeamRepo) Update(ctx context.Context, team *domain.Team) error {
	return nil
}
