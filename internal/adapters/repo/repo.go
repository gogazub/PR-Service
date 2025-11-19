package repo

import (
	pullrequestrepo "PRService/internal/adapters/repo/pullrequest"
	teamrepo "PRService/internal/adapters/repo/team"
	userrepo "PRService/internal/adapters/repo/user"
)

type Repository struct {
	userRepo        userrepo.Repo
	teamRepo        teamrepo.Repo
	pullRequestRepo pullrequestrepo.Repo
}

func (r *Repository) UserRepo() userrepo.Repo {
	return r.userRepo
}

func (r *Repository) TeamRepo() teamrepo.Repo {
	return r.teamRepo
}
func (r *Repository) PullRequestRepo() pullrequestrepo.Repo {
	return r.pullRequestRepo
}
