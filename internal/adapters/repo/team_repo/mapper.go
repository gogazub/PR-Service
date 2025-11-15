package teamrepo

import (
	"PRService/internal/domain/team"
)

func teamModelToDomain(m TeamModel) *team.Team {
	return &team.Team{
		Name:    m.Name,
		Members: m.Members,
	}
}

func teamDomainToModel(t *team.Team) TeamModel {
	return TeamModel{
		Name:    t.Name,
		Members: t.Members,
	}
}
