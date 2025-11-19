package userrepo

import "time"

type UserModel struct {
	UserID   string `db:"user_id"`
	Name     string `db:"user_name"`
	TeamName string `db:"team_name"`
	IsActive bool   `db:"is_active"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}
