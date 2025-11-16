package pullrequestrepo

import (
	"PRService/internal/domain/pullrequest"
	"PRService/internal/domain/user"
)

func pullRequestModelToDomain(m PullRequestModel, reviewers []user.ID) *pullrequest.PullRequest {
	var author user.ID
	if m.AuthorID != "" {
		author = user.ID(m.AuthorID)
	}

	var reviewersArr []user.ID
	for i := 0; i < len(reviewersArr) && i < len(reviewers); i++ {
		reviewersArr = append(reviewersArr, reviewers[i])
	}

	return &pullrequest.PullRequest{
		PullRequestID: pullrequest.ID(m.PRID),
		Name:          m.Name,
		Author:        author,
		Status:        pullrequest.StringToStatus(m.Status),
		Reviewers:     reviewersArr,
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
		Status:   pullrequest.StatusToString(pr.Status),
	}
}
