package team

import "PRService/internal/domain/user"

type ID string

type Team struct {
	Name    string
	Members []user.ID
}

// NewTeam returns new Team.
func NewTeam(name string, members []user.ID) *Team {
	return &Team{Name: name, Members: members}
}
