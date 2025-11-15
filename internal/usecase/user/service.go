package user_usecase

import (
	"PRService/internal/domain/user"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type Service interface {
	Save(ctx context.Context, u *user.User)  (*user.User, error)
	CreateUser(ctx context.Context, name string) (*user.User, error)
	GetByID(ctx context.Context, id user.ID) (*user.User, error)
	DeleteByID(ctx context.Context, id user.ID) error
	UpdateActive(ctx context.Context, cmd UpdateActiveCommand) error
}

type service struct {
	userRepo user.Repo
}

func New(repo user.Repo) Service {
	return &service{userRepo: repo}
}

func (svc *service)	Save(ctx context.Context, u *user.User)  (*user.User, error){
	err := svc.userRepo.Save(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return u, nil
}


func (svc *service) CreateUser(ctx context.Context, name string) (*user.User, error) {
	id := uuid.New().String()
	u := user.NewUser(id, name, false)

	err := svc.userRepo.Save(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return u, nil
}

func (svc *service) GetByID(ctx context.Context, userID user.ID) (*user.User, error) {

	u, err := svc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user by id %s: %w", userID, err)

	}
	return u, nil
}

func (svc *service) DeleteByID(ctx context.Context, userID user.ID) error {
	err := svc.userRepo.DeleteByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("delete by id: user id: %s: %w", userID, err)
	}
	return nil
}

func (svc *service) UpdateActive(ctx context.Context, cmd UpdateActiveCommand) error {
	err := svc.userRepo.UpdateActive(ctx, cmd.UserID, cmd.IsActive)
	if err != nil {
		return fmt.Errorf("update active: (id active): (%s, %t): %w", cmd.UserID, cmd.IsActive, err)
	}
	return nil
}
