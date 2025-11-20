package teamusecase

import (
	"PRService/internal/domain/user"
)

type CreateTeamAndUsersCommand struct {
	Name    string
	Members []*user.User
}

type CreateTeamCommand struct {
	Name    string
	Members []user.ID
}

type UpdateTeamCommand struct {
	Name    string
	Members []user.ID
}
