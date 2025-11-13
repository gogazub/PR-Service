package domain

type PullRequestID string

type PullRequestStatus int

const (
	OPEN PullRequestStatus = iota
	MERGED
)

type PullRequest struct {
	PullRequestID PullRequestID
	Name          string
	Author        UserID
	Status        PullRequestStatus
	Reviewers     [2]UserID
}

// NewPullRequest returns new PullRequest.
func NewPullRequest(id string, name string, author UserID, status PullRequestStatus, reviewers []UserID) *PullRequest {
	pr := &PullRequest{
		PullRequestID: PullRequestID(id),
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
