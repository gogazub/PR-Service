package teamrepo

import (
	"PRService/internal/domain/team"
	"PRService/internal/domain/user"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type Repo struct {
	db *sql.DB
}

func New(db *sql.DB) *Repo {
	return &Repo{db: db}
}

// Save saves new team
func (r *Repo) Save(ctx context.Context, t *team.Team) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("save team: begin tx: %w", err)
	}

	defer func() {
		_ = tx.Rollback()
	}()

	const qTeam = `
		INSERT INTO teams team_name
		VALUES $1
	`

	if _, err := tx.ExecContext(ctx, qTeam, t.Name); err != nil {
		return fmt.Errorf("save team: insert team: %w", err)
	}

	const qMembers = `
        INSERT INTO user_teams (user_id, team_name)
        VALUES ($1, $2)
    `

	for _, userID := range t.Members {
		if _, err := tx.ExecContext(ctx, qMembers, userID, t.Name); err != nil {
			return fmt.Errorf("save team: insert member %s: %w", userID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("save team: commit: %w", err)
	}

	return nil
}

// GetByName returns team by name
func (r *Repo) GetByName(ctx context.Context, name string) (*team.Team, error) {
	const qTeam = `
        SELECT team_id, team_name
        FROM teams
        WHERE team_name = $1
    `

	var tm TeamModel

	row := r.db.QueryRowContext(ctx, qTeam, name)
	if err := row.Scan(&tm.Name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, team.ErrTeamNotFound
		}
		return nil, fmt.Errorf("get by name: scan: %w", err)
	}

	members, err := r.loadMembers(ctx, nil, tm.Name)
	if err != nil {
		return nil, fmt.Errorf("get by name: team: %s: load members: %w", tm.Name, err)
	}
	tm.Members = members

	return teamModelToDomain(tm), nil
}

// GetActiveUsersInTeam returns all active users in team
func (r *Repo) GetActiveUsersInTeam(ctx context.Context, teamName string) ([]*user.User, error) {
	const q = `
		SELECT u.name, u.is_active
		FROM users u
		INNER JOIN user_teams ut ON ut.user_id = u.user_id
		WHERE ut.team_name = $1
		  AND u.is_active = TRUE
	`

	rows, err := r.db.QueryContext(ctx, q, teamName)
	if err != nil {
		return nil, fmt.Errorf("get active users in team: query: %w", err)
	}
	defer rows.Close()

	var usersList []*user.User

	for rows.Next() {
		u := &user.User{}
		if err := rows.Scan(&u.UserID, &u.Name, &u.IsActive); err != nil {
			return nil, fmt.Errorf("get active users in team: scan: %w", err)
		}
		usersList = append(usersList, u)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("get active users in team: rows err: %w", err)
	}

	return usersList, nil
}

// Update modifies an existing team
func (r *Repo) Update(ctx context.Context, t *team.Team) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("update team: begin tx: %w", err)
	}

	defer func() {
		_ = tx.Rollback()
	}()

	const qUpdateTeam = `
		UPDATE teams
		SET team_name = $1
		WHERE team_name = $2
	`

	res, err := tx.ExecContext(ctx, qUpdateTeam, t.Name, t.Name)
	if err != nil {
		return fmt.Errorf("update team: update team: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("update team: rows affected: %w", err)
	}
	if affected == 0 {
		return team.ErrTeamNotFound
	}

	const qDeleteMembers = `
		DELETE FROM user_teams
		WHERE team_name = $1
	`
	if _, err := tx.ExecContext(ctx, qDeleteMembers, t.Name); err != nil {
		return fmt.Errorf("update team: delete members: %w", err)
	}

	const qInsertMember = `
        INSERT INTO user_teams (user_name, team_name)
        VALUES ($1, $2)
    `
	for _, userID := range t.Members {
		if _, err := tx.ExecContext(ctx, qInsertMember, userID, t.Name); err != nil {
			return fmt.Errorf("update team: insert member %s: %w", userID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("update team: commit: %w", err)
	}

	return nil
}

// DeleteByName deletes team by name
func (r *Repo) DeleteByName(ctx context.Context, name string) error {
	const q = `
		DELETE FROM teams
		WHERE team_name = $1
	`

	res, err := r.db.ExecContext(ctx, q, name)
	if err != nil {
		return fmt.Errorf("delete by name: exec: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete by name: rows affected: %w", err)
	}
	if affected == 0 {
		return team.ErrTeamNotFound
	}

	return nil
}

func (r *Repo) loadMembers(ctx context.Context, tx *sql.Tx, teamName string) ([]user.ID, error) {
	const q = `
		SELECT user_id
		FROM user_teams
		WHERE team_name = $1
	`

	var (
		rows *sql.Rows
		err  error
	)

	if tx == nil {
		rows, err = r.db.QueryContext(ctx, q, teamName)
	} else {
		rows, err = tx.QueryContext(ctx, q, teamName)
	}

	if err != nil {
		return nil, fmt.Errorf("load members: query: %w", err)
	}
	defer rows.Close()

	var members []user.ID

	for rows.Next() {
		var uid user.ID
		if err := rows.Scan(&uid); err != nil {
			return nil, fmt.Errorf("load members: scan: %w", err)
		}
		members = append(members, uid)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("load members: rows err: %w", err)
	}

	return members, nil
}
