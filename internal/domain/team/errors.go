package team

import "errors"

var ErrTeamNotFound = errors.New("team not found")
var ErrTeamExists = errors.New("team is already exists")
var ErrNoActiveMembers = errors.New("no active members")
