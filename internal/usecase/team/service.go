package team_usecase

import (
	"PRService/internal/domain/team"
	"context"
	"os/user"
)

type Service interface {
	CreateTeam(ctx context.Context, cmd CreateTeamCommand) (*team.Team, error)
	GetByName(ctx context.Context, name string) (*team.Team, error)
	GetByID(ctx context.Context, id team.ID) (*team.Team, error)
	GetActiveUsersInTeam(ctx context.Context, id team.ID) ([]*user.User, error)
	Update(ctx context.Context, cmd UpdateTeamCommand) error
	DeleteByID(ctx context.Context, id team.ID) error
	DeleteByName(ctx context.Context, name string) error
}

type teamService struct {
	teamRepo team.Repo
}

func (svc *teamService) CreateTeam(ctx context.Context, cmd CreateTeamCommand) (*team.Team, error) {
	id := "id" // TODO: id generator
	team := team.NewTeam(team.ID(id), cmd.Name, cmd.Members)
	return team, svc.teamRepo.Save(ctx, team)
}

func (svc *teamService) GetByName(ctx context.Context, name string) (*team.Team, error) {
	return svc.teamRepo.GetByName(ctx, name)
}

func (svc *teamService) GetByID(ctx context.Context, id team.ID) (*team.Team, error) {
	return svc.teamRepo.GetByID(ctx, id)
}

func (svc *teamService) GetActiveUsersInTeam(ctx context.Context, teamID team.ID) ([]*user.User, error) {
	return svc.teamRepo.GetActiveUsersInTeam(ctx, string(teamID))
}

func (svc *teamService) Update(ctx context.Context, cmd UpdateTeamCommand) error {
	team := new(team.Team)
	team.ID = cmd.TeamID
	if cmd.Name != "" {
		team.Name = cmd.Name
	}
	if len(cmd.Members) != 0 {
		team.Members = cmd.Members
	}

	return svc.teamRepo.Update(ctx, team)
}

func (svc *teamService) DeleteByID(ctx context.Context, teamID team.ID) error {
	return svc.teamRepo.DeleteByID(ctx, teamID)
}

func (svc *teamService) DeleteByName(ctx context.Context, teamName string) error {
	return svc.teamRepo.DeleteByName(ctx, teamName)
}
