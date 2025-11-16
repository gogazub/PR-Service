package app

import (
	"PRService/internal/domain/team"
	"PRService/internal/domain/user"
	pullrequest_usecase "PRService/internal/usecase/pullrequest"
	team_usecase "PRService/internal/usecase/team"
	user_usecase "PRService/internal/usecase/user"
	"context"
	"fmt"
)

// Агрегированные в одну сущность сервисы. Точка взаимодействия с приложение.
type Services struct {
	User        user_usecase.Service
	Team        team_usecase.Service
	PullRequest pullrequest_usecase.Service
}

func NewServices(user user_usecase.Service, team team_usecase.Service, pr pullrequest_usecase.Service) *Services {
	return &Services{User: user, Team: team, PullRequest: pr}
}

func (svc *Services) CreateTeam(ctx context.Context, cmd team_usecase.CreateTeamAndUsersCommand) (*team.Team, []*user.User, error) {
	// TODO: создать здесь транзакцию. Либо создаем и team и все users, либо ничего 
	
	// 1. Make cmd
	members := cmd.Members
	ids := make([]user.ID, 0, len(members))
	for _, u := range members {
		ids = append(ids, u.UserID)
	}
	createTeamCmd := team_usecase.CreateTeamCommand{
		Name: cmd.Name,
		Members: ids,
	}

	// 2. Save team
	t, err := svc.Team.CreateTeam(ctx, createTeamCmd)
	if err != nil {
		return nil, nil, fmt.Errorf("service: create team: team: %q: %w", createTeamCmd.Name, err)
	}

	// 3. Save users
	users := make([]*user.User, 0, len(members))
	for _, u := range members {
		_, err := svc.User.CreateUser(ctx, u)
		if err != nil {
			return nil, nil, fmt.Errorf("service: create team: create user: user id: %s: %w", u.UserID, err)
		}
		users = append(users, u)
	}

	return t, users, nil
	
}