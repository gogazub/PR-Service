package pullrequestrepo

import (
	"PRService/internal/domain/pullrequest"
	"PRService/internal/domain/user"
)

func pullRequestModelToDomain(m PullRequestModel) *pullrequest.PullRequest {
	var author user.ID
	if m.AuthorID != "" {
		author = user.ID(m.AuthorID)
	}

	return &pullrequest.PullRequest{
		PullRequestID: pullrequest.ID(m.PRID),
		Name:          m.Name,
		Author:        author,
		Status:        m.Status,
		Reviewers:     [2]user.ID{},
	}
}

func pullRequestDomainToModel(pr *pullrequest.PullRequest) PullRequestModel {
	// Может быть такое, что со временем автор pr будет удален.
	authorID := ""
	if pr.Author != "" {
		authorID = string(pr.Author)
	}

	return PullRequestModel{
		PRID:     string(pr.PullRequestID),
		Name:     pr.Name,
		AuthorID: authorID,
		Status:   pr.Status,
	}
}
