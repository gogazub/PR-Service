package teamhandlers

import (
	teamhttp "PRService/internal/adapters/http/team"
	"PRService/internal/app"
	"PRService/internal/domain/user"
	"PRService/pkg/logger"
)

type Handler struct {
	*app.Services
	logger *logger.Logger
}

// NewHandler returns new Handler.
func NewHandler(app *app.Services, logger *logger.Logger) *Handler {
	return &Handler{app, logger}
}

// --- helpers ---
func UsersFromMembers(members []teamhttp.TeamMemberDTO, teamName string) []*user.User {
	users := make([]*user.User, 0, len(members))
	for _, member := range members {
		u := memberToUser(&member, teamName)
		users = append(users, u)
	}
	return users
}

func MembersFromUsers(usrs []*user.User) []teamhttp.TeamMemberDTO {
	members := make([]teamhttp.TeamMemberDTO, 0, len(usrs))
	for _, u := range usrs {
		member := userToMember(u)
		members = append(members, member)
	}
	return members
}

func memberToUser(member *teamhttp.TeamMemberDTO, teamName string) *user.User {
	u := &user.User{
		UserID:   user.ID(member.UserID),
		Name:     member.Username,
		TeamName: teamName,
		IsActive: member.IsActive,
	}
	return u
}

func userToMember(u *user.User) teamhttp.TeamMemberDTO {
	member := new(teamhttp.TeamMemberDTO)
	member.UserID = string(u.UserID)
	member.Username = u.Name
	member.IsActive = u.IsActive
	return *member
}
