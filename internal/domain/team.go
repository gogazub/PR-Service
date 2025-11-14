package domain

type TeamID string

type Team struct {
	TeamID  string
	Name    string
	Members []UserID
}

// NewTeam returns new Team.
func NewTeam() *Team {
	return &Team{}
}
