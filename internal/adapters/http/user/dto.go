package userhttp

import (
	"PRService/internal/domain/pullrequest"
	"PRService/internal/domain/user"
)

type UserDTO struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

type PullRequestShortDTO struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
	Status          string `json:"status"`
}

func UserToDTO(u *user.User) *UserDTO {
	return &UserDTO{
		UserID:   string(u.UserID),
		Username: u.Name,
		TeamName: u.TeamName,
		IsActive: u.IsActive,
	}
}

func PullRequestToShortDTO(pr *pullrequest.PullRequest) *PullRequestShortDTO {
	return &PullRequestShortDTO{
		PullRequestID:   string(pr.PullRequestID),
		PullRequestName: pr.Name,
		AuthorID:        string(pr.Author),
		Status:          pullrequest.StatusToString(pr.Status),
	}
}

