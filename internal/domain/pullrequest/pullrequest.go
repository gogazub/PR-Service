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
	Reviewers     [2]user.ID
}

// NewPullRequest returns new PullRequest.
func NewPullRequest(id string, name string, author user.ID, status Status, reviewers []user.ID) *PullRequest {
	pr := &PullRequest{
		PullRequestID: ID(id),
		Name:          name,
		Author:        author,
		Status:        status,
	}

	// Take only first two elements
	for i, val := range reviewers {
		if i > 1 {
			break
		}
		pr.Reviewers[i] = val
	}
	return pr
}

func StatusToString(s Status) string {
	if s == OPEN {
		return "OPEN"
	}
	return "MERGED"
}