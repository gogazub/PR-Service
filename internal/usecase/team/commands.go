package team_usecase

import (
	"PRService/internal/domain/user"
)

type CreateTeamCommand struct {
	Name    string
	Members []user.ID
}

type UpdateTeamCommand struct {
	TeamID  string
	Name    string
	Members []user.ID
}
