package userrepo

import (
	"database/sql"
)

type UserRepo struct {
	db *sql.Conn
}

// NewUserRepo returns new UserRepo.
func NewUserRepo(db *sql.Conn) *UserRepo {
	return &UserRepo{db: db}
}
/*
func (r *UserRepo) GetByID(ctx context.Context, id domain.UserID) (*domain.User, error) {

}

func (r *UserRepo) Create(ctx context.Context, user *domain.User) error {

}

func (r *UserRepo) GetByTeamID(ctx context.Context, teamID string) ([]*domain.User, error) {

}

func (r *UserRepo) UpdateActive(ctx context.Context, id domain.UserID, isActive bool) error {

}
*/
