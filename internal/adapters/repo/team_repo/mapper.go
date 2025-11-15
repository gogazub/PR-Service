package teamrepo

import (
	"PRService/internal/domain/team"
)

func teamModelToDomain(m TeamModel) *team.Team {
	return &team.Team{
		ID:      team.ID(m.ID),
		Name:    m.Name,
		Members: m.Members,
	}
}

func teamDomainToModel(t *team.Team) TeamModel {
	return TeamModel{
		ID:   string(t.ID),
		Name: t.Name,
	}
}
