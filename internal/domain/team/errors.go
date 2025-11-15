package team

import "errors"

var ErrTeamNotFound = errors.New("team not found")
var ErrTeamExists = errors.New("team is already exists")
