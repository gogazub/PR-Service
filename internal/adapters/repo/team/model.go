package teamrepo

import (
	"PRService/internal/domain/user"
	"time"
)

type TeamModel struct {
	Name    string `db:"team_name"`
	Members []user.ID

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}
