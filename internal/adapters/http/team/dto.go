package teamhttp

import (
	"PRService/internal/domain/team"
	"PRService/internal/domain/user"
)

type TeamMemberDTO struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type TeamDTO struct {
	TeamName string          `json:"team_name"`
	Members  []TeamMemberDTO `json:"members"`
}

func TeamToDTO(t *team.Team, users []*user.User) TeamDTO {
    dto := TeamDTO{
        TeamName: t.Name,
        Members:  make([]TeamMemberDTO, 0, len(users)),
    }

    for _, u := range users {
        dto.Members = append(dto.Members, TeamMemberDTO{
            UserID:   string(u.UserID),
            Username: u.Name,
            IsActive: u.IsActive,
        })
    }

    return dto
}

func ExtractMemberIDs(members []TeamMemberDTO) []user.ID {
    ids := make([]user.ID, 0, len(members))
    for _, m := range members {
        ids = append(ids, user.ID(m.UserID))
    }
    return ids
}
