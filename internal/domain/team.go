package domain

type Team struct {
	TeamID  string
	Name    string
	Members []UserID
}

// NewTeam returns new Team.
func NewTeam() *Team {
	return &Team{}
}
