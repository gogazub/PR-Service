package pullrequestrepo

import (
	"PRService/internal/domain"
	"time"
)

type PullRequestModel struct {
	PRID     string                   `db:"pr_id"`
	Name     string                   `db:"name"`
	AuthorID string                   `db:"author_id"`
	Status   domain.PullRequestStatus `db:"status"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}
