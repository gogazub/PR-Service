package userusecase

import "PRService/internal/domain/user"

type UpdateActiveCommand struct {
	UserID   user.ID
	IsActive bool
}
