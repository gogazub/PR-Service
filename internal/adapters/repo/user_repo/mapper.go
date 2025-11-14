package userrepo

import "PRService/internal/domain"

func userModelToDomain(m UserModel) *domain.User {
	return &domain.User{
		UserID:   domain.UserID(m.UserID),
		Name:     m.Name,
		IsActive: m.IsActive,
	}
}

func userDomainToModel(u *domain.User) UserModel {
	return UserModel{
		UserID:   string(u.UserID),
		Name:     u.Name,
		IsActive: u.IsActive,
	}
}
