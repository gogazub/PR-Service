package user_usecase

import (
	"PRService/internal/domain/user"
	"context"
)

type Service interface {
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

func (svc *service) CreateUser(ctx context.Context, name string) (*user.User, error) {
	id := "id" // genarate id
	user := user.NewUser(id, name, false)
    // TODO: process error
	err := svc.userRepo.Save(ctx, user)
	return user, err
}

func (svc *service) GetByID(ctx context.Context, userID user.ID) (*user.User, error) {
	return svc.userRepo.GetByID(ctx, userID)

}

func (svc *service) DeleteByID(ctx context.Context, userID user.ID) error {
	return svc.userRepo.DeleteByID(ctx, userID)
}

func (svc *service) UpdateActive(ctx context.Context, cmd UpdateActiveCommand) error {
	return svc.userRepo.UpdateActive(ctx, cmd.UserID, cmd.IsActive)
}
