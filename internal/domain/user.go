package domain

type UserID string

type User struct {
	UserID   UserID
	Name     string
	IsActive bool
}

// NewUser returns new User.
func NewUser(user string, name string, isActive bool) *User {
	return &User{UserID: UserID(user), Name: name, IsActive: isActive}
}
