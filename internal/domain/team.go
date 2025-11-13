package domain

type Team struct {
	Name    string
	Members []UserID
}

// NewTeam returns new Team.
func NewTeam() *Team {
	return &Team{}
}
