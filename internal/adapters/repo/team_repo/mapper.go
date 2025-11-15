package teamrepo

import (
	"PRService/internal/domain/team"
	"PRService/internal/domain/user"
)

func teamModelToDomain(m TeamModel) *team.Team {
	return &team.Team{
		TeamID:  m.TeamID,
		Name:    m.TeamName,
		Members: []user.ID{},
	}
}

func teamDomainToModel(t *team.Team) TeamModel {
	return TeamModel{
		TeamID:   t.TeamID,
		TeamName: t.Name,
	}
}
