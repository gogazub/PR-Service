package teamrepo

import "time"

type TeamModel struct {
	TeamID   string `db:"team_id"`
	TeamName string `db:"team_name"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}
