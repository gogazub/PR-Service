package teamrepo

import "PRService/internal/domain"

func teamModelToDomain(m TeamModel) *domain.Team {
	return &domain.Team{
		TeamID:  m.TeamID,
		Name:    m.TeamName,
		Members: []domain.UserID{},
	}
}

func teamDomainToModel(t *domain.Team) TeamModel {
	return TeamModel{
		TeamID:   t.TeamID,
		TeamName: t.Name,
	}
}
