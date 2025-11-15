package pullrequest

import "errors"

var ErrPullRequestNotFound = errors.New("pull request not found")
var ErrPullRequestExists = errors.New("pull request is already exists")
