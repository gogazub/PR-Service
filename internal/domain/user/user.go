package user

type ID string

type User struct {
	UserID   ID
	Name     string
	TeamName string
	IsActive bool
}

// NewUser returns new User.
func NewUser(id string, name string, teamName string, isActive bool) *User {
	return &User{UserID: ID(id), Name: name, TeamName: teamName, IsActive: isActive}
}
