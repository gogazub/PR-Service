package team_usecase

import (
	"PRService/internal/domain/team"
	"PRService/internal/domain/user"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type Service interface {
	CreateTeam(ctx context.Context, cmd CreateTeamCommand) (*team.Team, error)

	GetByName(ctx context.Context, name string) (*team.Team, error)
	GetByID(ctx context.Context, id string) (*team.Team, error)

	GetActiveUsersInTeam(ctx context.Context, id string) ([]*user.User, error)

	Update(ctx context.Context, cmd UpdateTeamCommand) error

	DeleteByID(ctx context.Context, id string) error
	DeleteByName(ctx context.Context, name string) error
}

type service struct {
	teamRepo team.Repo
}

func New(repo team.Repo) Service {
	return &service{teamRepo: repo}
}

func (svc *service) CreateTeam(ctx context.Context, cmd CreateTeamCommand) (*team.Team, error) {

	id := team.ID(uuid.New().String())

	t := team.NewTeam(
		id,
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

func (svc *service) GetByID(ctx context.Context, id string) (*team.Team, error) {

	t, err := svc.teamRepo.GetByID(ctx, team.ID(id))
	if err != nil {
		return nil, fmt.Errorf("get team by id %s: %w", id, err)
	}

	return t, nil
}

func (svc *service) GetActiveUsersInTeam(ctx context.Context, id string) ([]*user.User, error) {

	users, err := svc.teamRepo.GetActiveUsersInTeam(ctx, team.ID(id))
	if err != nil {
		return nil, fmt.Errorf("get active users in team %s: %w", id, err)
	}

	return users, nil
}

func (svc *service) Update(ctx context.Context, cmd UpdateTeamCommand) error {

	t := &team.Team{
		ID:      team.ID(cmd.TeamID),
		Name:    cmd.Name,
		Members: cmd.Members,
	}

	if err := svc.teamRepo.Update(ctx, t); err != nil {
		return fmt.Errorf("update team (id: %s): %w", cmd.TeamID, err)
	}

	return nil
}

func (svc *service) DeleteByID(ctx context.Context, id string) error {

	if err := svc.teamRepo.DeleteByID(ctx, team.ID(id)); err != nil {
		return fmt.Errorf("delete team by id %s: %w", id, err)
	}

	return nil
}

func (svc *service) DeleteByName(ctx context.Context, name string) error {

	if err := svc.teamRepo.DeleteByName(ctx, name); err != nil {
		return fmt.Errorf("delete team by name %q: %w", name, err)
	}

	return nil
}
