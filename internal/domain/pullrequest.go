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
	Status        PullRequestStatus
	Reviewers     [2]UserID
}

// NewPullRequest returns new PullRequest.
func NewPullRequest(id string, name string, status PullRequestStatus, reviewers []UserID) *PullRequest {
	pr := &PullRequest{
		PullRequestID: PullRequestID(id),
		Name:          name,
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