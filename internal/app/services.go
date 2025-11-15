package app

import (
	pullrequest_usecase "PRService/internal/usecase/pullrequest"
	team_usecase "PRService/internal/usecase/team"
	user_usecase "PRService/internal/usecase/user"
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
