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
