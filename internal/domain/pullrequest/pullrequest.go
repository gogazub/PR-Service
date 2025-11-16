package pullrequest

import "PRService/internal/domain/user"

type ID string

type Status int

const (
	OPEN Status = iota
	MERGED
)

type PullRequest struct {
	PullRequestID ID
	Name          string
	Author        user.ID
	Status        Status
	Reviewers     []user.ID
}

// NewPullRequest returns new PullRequest.
func NewPullRequest(id string, name string, author user.ID, status Status, reviewers []user.ID) *PullRequest {
	pr := &PullRequest{
		PullRequestID: ID(id),
		Name:          name,
		Author:        author,
		Status:        status,
		Reviewers:     reviewers,
	}
	return pr
}

func StatusToString(s Status) string {
	if s == OPEN {
		return "OPEN"
	}
	return "MERGED"
}