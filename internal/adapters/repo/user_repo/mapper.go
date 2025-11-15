package userrepo

import (
	"PRService/internal/domain/user"
)

func userModelToDomain(m UserModel) *user.User {
	return &user.User{
		UserID:   user.ID(m.UserID),
		Name:     m.Name,
		IsActive: m.IsActive,
	}
}

func userDomainToModel(u *user.User) UserModel {
	return UserModel{
		UserID:   string(u.UserID),
		Name:     u.Name,
		IsActive: u.IsActive,
	}
}
