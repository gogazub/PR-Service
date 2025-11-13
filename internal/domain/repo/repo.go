package repo

// Composition Repository for encapsulation UserRepo, TeamRepo, PRRepo details.
type Repository interface {
	UserRepo() UserRepo
	TeamRepo() TeamRepo
	PullRequestRepo() PullRequestRepo
}
