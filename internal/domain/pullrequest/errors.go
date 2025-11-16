package pullrequest

import "errors"

var ErrPullRequestNotFound = errors.New("pull request not found")
var ErrPullRequestExists = errors.New("pull request is already exists")
var ErrNoAuthor = errors.New("no author")
var ErrReviewerNotAssigned = errors.New("reviewer not assigned")
var ErrNoCandidate = errors.New("no candidate")
var ErrMerged = errors.New("merged")
