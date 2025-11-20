package teamusecase

import (
	"PRService/internal/domain/team"
	"PRService/internal/domain/user"
	"context"
	"fmt"
)

type Service interface {
	Save(ctx context.Context, t *team.Team) (*team.Team, error)
	CreateTeam(ctx context.Context, cmd CreateTeamCommand) (*team.Team, error)

	GetByName(ctx context.Context, name string) (*team.Team, error)

	GetActiveUsersInTeam(ctx context.Context, name string) ([]*user.User, error)

	Update(ctx context.Context, cmd UpdateTeamCommand) error

	DeleteByName(ctx context.Context, name string) error
}

type service struct {
	teamRepo team.Repo
}

func New(repo team.Repo) Service {
	return &service{teamRepo: repo}
}

func (svc *service) Save(ctx context.Context, t *team.Team) (*team.Team, error) {
	err := svc.teamRepo.Save(ctx, t)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (svc *service) CreateTeam(ctx context.Context, cmd CreateTeamCommand) (*team.Team, error) {

	if t, _:= svc.teamRepo.GetByName(ctx, cmd.Name); t != nil {
		return nil, team.ErrTeamExists
	}

	t := team.NewTeam(
		cmd.Name,
		cmd.Members,
	)

	if err := svc.teamRepo.Save(ctx, t); err != nil {
		return nil, fmt.Errorf("create team: %w", err)
	}

	return t, nil
}

func (svc *service) GetByName(ctx context.Context, name string) (*team.Team, error) {

	t, err := svc.teamRepo.GetByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("get team by name %q: %w", name, err)
	}

	return t, nil
}

func (svc *service) GetActiveUsersInTeam(ctx context.Context, name string) ([]*user.User, error) {

	users, err := svc.teamRepo.GetActiveUsersInTeam(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("get active users in team %q: %w", name, err)
	}

	return users, nil
}

func (svc *service) Update(ctx context.Context, cmd UpdateTeamCommand) error {

	t := &team.Team{
		Name:    cmd.Name,
		Members: cmd.Members,
	}

	if err := svc.teamRepo.Update(ctx, t); err != nil {
		return fmt.Errorf("update team (name: %s): %w", cmd.Name, err)
	}

	return nil
}

func (svc *service) DeleteByName(ctx context.Context, name string) error {

	if err := svc.teamRepo.DeleteByName(ctx, name); err != nil {
		return fmt.Errorf("delete team by name %q: %w", name, err)
	}

	return nil
}
