package pullrequesthttp

import "time"

type PullRequestDTO struct {
	PullRequestID     string     `json:"pull_request_id"`
	PullRequestName   string     `json:"pull_request_name"`
	AuthorID          string     `json:"author_id"`
	Status            string     `json:"status"`
	AssignedReviewers []string   `json:"assigned_reviewers"`
	CreatedAt         *time.Time `json:"createdAt,omitempty"`
	MergedAt          *time.Time `json:"mergedAt,omitempty"`
}

// POST /pullRequest/create
type CreatePullRequestRequestDTO struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
}

type CreatePullRequestResponseDTO struct {
	PR PullRequestDTO `json:"pr"`
}

// POST /pullRequest/merge
type MergePullRequestRequestDTO struct {
	PullRequestID string `json:"pull_request_id"`
}

type MergePullRequestResponseDTO struct {
	PR PullRequestDTO `json:"pr"`
}

// POST /pullRequest/reassign
type ReassignReviewerRequestDTO struct {
	PullRequestID string `json:"pull_request_id"`
	OldUserID     string `json:"old_user_id"`
}

type ReassignReviewerResponseDTO struct {
	PR         PullRequestDTO `json:"pr"`
	ReplacedBy string         `json:"replaced_by"`
}
