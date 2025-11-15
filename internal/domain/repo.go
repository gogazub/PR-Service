package domain

import (
	"PRService/internal/domain/pullrequest"
	"PRService/internal/domain/team"
	"PRService/internal/domain/user"
)

// Composition Repository for encapsulation UserRepo, TeamRepo, PRRepo details.
type Repository interface {
	UserRepo() user.Repo
	TeamRepo() team.Repo
	PullRequestRepo() pullrequest.Repo
}
