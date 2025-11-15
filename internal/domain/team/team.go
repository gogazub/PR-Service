package team

import "PRService/internal/domain/user"

type ID string

type Team struct {
	ID      ID
	Name    string
	Members []user.ID
}

// NewTeam returns new Team.
func NewTeam(id ID, name string, members []user.ID) *Team {
	return &Team{ID: id, Name: name, Members: members}
}
