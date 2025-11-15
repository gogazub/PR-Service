package team_usecase

import (
	"PRService/internal/domain/team"
	"PRService/internal/domain/user"
)

type CreateTeamCommand struct {
	Name    string
	Members []user.ID
}

type UpdateTeamCommand struct {
	TeamID  team.ID
	Name    string
	Members []user.ID
}
