package repo

import (
	pullrequestrepo "PRService/internal/adapters/repo/pullrequest_repo"
	teamrepo "PRService/internal/adapters/repo/team_repo"
	userrepo "PRService/internal/adapters/repo/user_repo"
)

type Repository struct {
	userRepo        userrepo.UserRepo
	teamRepo        teamrepo.TeamRepo
	pullRequestRepo pullrequestrepo.PullRequestRepo
}

func (r *Repository) UserRepo() userrepo.UserRepo {
	return r.userRepo
}

func (r *Repository) TeamRepo() teamrepo.TeamRepo {
	return r.teamRepo
}
func (r *Repository) PullRequestRepo() pullrequestrepo.PullRequestRepo {
	return r.pullRequestRepo
}
