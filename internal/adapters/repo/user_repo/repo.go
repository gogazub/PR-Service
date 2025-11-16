package userrepo

import (
	"PRService/internal/domain/user"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
)

type Repo struct {
	db *sql.DB
}

func New(db *sql.DB) *Repo {
	return &Repo{db}
}

// Save saves new user. if exists -> update
func (r *Repo) Save(ctx context.Context, u *user.User) error {
	const q = `INSERT INTO users (user_id, user_name, team_name, is_active)
				VALUES ($1, $2, $3, $4)
				ON CONFLICT (user_id) DO UPDATE
				SET user_name = EXCLUDED.user_name,
					team_name = EXCLUDED.team_name,
					is_active = EXCLUDED.is_active;
			`

	_, err := r.db.ExecContext(ctx, q, u.UserID, u.Name, u.TeamName, u.IsActive)
	if err != nil {
		return fmt.Errorf("user repo: save: user: %s: %w", u.UserID, err)
	}

	return nil
}

// GetByID returns user by ID
func (r *Repo) GetByID(ctx context.Context, id user.ID) (*user.User, error) {
	const q = `
		SELECT user_id, user_name, team_name, is_active
		FROM users
		WHERE user_id = $1
		`
	row := r.db.QueryRowContext(ctx, q, id)

	var u UserModel
	if err := row.Scan(&u.UserID, &u.Name, &u.TeamName, &u.IsActive); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("get user by id: user: %s: %w", id, user.ErrUserNotFound)
		}
		return nil, fmt.Errorf("get user by id: user: %s: scan: %w", id, err)
	}

	return userModelToDomain(&u), nil

}

// GetByIDs returns users by IDs
func (r *Repo) GetByIDs(ctx context.Context, ids []user.ID) ([]*user.User, error) {
	const q = `
		SELECT user_id, user_name, team_name, is_active
		FROM users
		WHERE user_id = ANY($1)
	`

	if len(ids) == 0 {
		return []*user.User{}, nil
	}

	strIDs := make([]string, len(ids))
	for i, id := range ids {
		strIDs[i] = string(id)
	}

	rows, err := r.db.QueryContext(ctx, q, pq.Array(strIDs))
	if err != nil {
		return nil, fmt.Errorf("get users by ids: query: %w", err)
	}
	defer rows.Close()

	var users []*user.User
	for rows.Next() {
		u := new(user.User)
		if err := rows.Scan(&u.UserID, &u.Name, &u.TeamName, &u.IsActive); err != nil {
			return nil, fmt.Errorf("get users by ids: scan: %w", err)
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("get users by ids: rows: %w", err)
	}

	return users, nil
}

// UpdateActive updates user's active status
func (r *Repo) UpdateActive(ctx context.Context, id user.ID, isActive bool) (*user.User, error) {
	const q = `
        UPDATE users
        SET is_active = $1
        WHERE user_id = $2
        RETURNING user_name
    `

	var u user.User
	err := r.db.QueryRowContext(ctx, q, isActive, id).Scan(
		&u.Name,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("update active: user: %s: %w", id, user.ErrUserNotFound)
		}
		return nil, fmt.Errorf("update active: user: %s: %w", id, err)
	}

	return &u, nil
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
