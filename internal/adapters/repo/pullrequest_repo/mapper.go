package pullrequestrepo

import "PRService/internal/domain"

func pullRequestModelToDomain(m PullRequestModel) *domain.PullRequest {
	var author domain.UserID
	if m.AuthorID != "" {
		author = domain.UserID(m.AuthorID)
	}

	return &domain.PullRequest{
		PullRequestID: domain.PullRequestID(m.PRID),
		Name:          m.Name,
		Author:        author,
		Status:        m.Status,
		Reviewers:     [2]domain.UserID{},
	}
}

func pullRequestDomainToModel(pr *domain.PullRequest) PullRequestModel {
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
