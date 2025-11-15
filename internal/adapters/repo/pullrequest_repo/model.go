package pullrequestrepo

import (
	"PRService/internal/domain/pullrequest"
	"time"
)

type PullRequestModel struct {
	PRID     string             `db:"pr_id"`
	Name     string             `db:"name"`
	AuthorID string             `db:"author_id"` // пустая строка == NULL в БД
	Status   pullrequest.Status `db:"status"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}
