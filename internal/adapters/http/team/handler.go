package teamhttp

import (
	"PRService/internal/domain/user"
)



func UsersFromMembers(members []TeamMemberDTO, teamName string) []*user.User{
	users := make([]*user.User, 0, len(members))
	for _, member := range members {
		u := memberToUser(&member, teamName)
		users = append(users, u)
	}
	return users
}

func MembersFromUsers(usrs []*user.User) []TeamMemberDTO {
	members := make([]TeamMemberDTO, 0, len(usrs))
	for _, u := range usrs {
		member := userToMember(u)
		members = append(members, member)
	}
	return members
}

func memberToUser(member *TeamMemberDTO, teamName string) *user.User {
	u := &user.User{
		UserID: user.ID(member.UserID),
		Name: member.Username,
		TeamName: teamName,
		IsActive: member.IsActive,
	}
	return u
}

func userToMember(u *user.User) TeamMemberDTO {
	member := new(TeamMemberDTO)
	member.UserID = string(u.UserID)
	member.Username = u.Name
	member.IsActive = u.IsActive
	return *member
}