package userrepo

import (
	"PRService/internal/domain/user"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type Repo struct {
	db *sql.DB
}

func New(db *sql.DB) *Repo {
	return &Repo{db}
}

// Save saves new user
func (r *Repo) Save(ctx context.Context, u *user.User) error {
	const q = `INSERT INTO users (user_id, name, is_active)
			   VALUES ($1, $2, $3)`

	_, err := r.db.ExecContext(ctx, q, u.UserID, u.Name, u.IsActive)
	if err != nil {
		return fmt.Errorf("save: user: %s: %w", u.UserID, err)
	}

	return nil
}

// GetByID returns user by ID
func (r *Repo) GetByID(ctx context.Context, id user.ID) (*user.User, error) {
	const q = `
		SELECT user_id, name, is_active
		FROM users
		WHERE user_id = $1
		`
	row := r.db.QueryRowContext(ctx, q, id)

	var u UserModel
	if err := row.Scan(&u.UserID, &u.Name, &u.IsActive); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("get user by id: user: %s: %w", id, user.ErrUserNotFound)
		}
		return nil, fmt.Errorf("get user by id: user: %s: scan: %w", id, err)
	}

	return userModelToDomain(u), nil

}

// UpdateActive updates user's active status
func (r *Repo) UpdateActive(ctx context.Context, id user.ID, isActive bool) error {
	const q = `
		UPDATE users
		SET is_active = $1
		WHERE user_id = $2
	`
	res, err := r.db.ExecContext(ctx, q, isActive, id)
	if err != nil {
		return fmt.Errorf("update active: user: %s: %w", id, err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("update active: user: %s: %w", id, err)
	}
	if rows == 0 {
		return fmt.Errorf("update active: user: %s: %w", id, user.ErrUserNotFound)
	}

	return nil
}

// DeleteByID deletes user by id
func (r *Repo) DeleteByID(ctx context.Context, id user.ID) error {
	const q = `
		DELETE FROM users
		WHERE user_id = $1
	`

	// Можно также как и в Update добавить проверку затронутых строк.
	// Тогда можно будет кидать warning, что юзера не было
	_, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("delete by id: user: %s: %w", id, err)
	}

	return nil
}
