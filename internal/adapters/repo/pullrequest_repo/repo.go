package pullrequestrepo

import (
	"PRService/internal/domain/pullrequest"
	"PRService/internal/domain/user"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type Repo struct {
	db *sql.DB
}

func NewPullRequestRepo(db *sql.DB) *Repo {
	return &Repo{db: db}
}

// Save saves new pull request
func (r *Repo) Save(ctx context.Context, pr *pullrequest.PullRequest) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("save pull request: begin tx: %w", err)
	}

	defer func() {
		_ = tx.Rollback()
	}()

	// 1. Save (id name author status)
	const qPR = `
		INSERT INTO pull_requests (pr_id, name, author_id, status)
		VALUES ($1, $2, $3, $4)
	`

	var authorID any
	if pr.Author != "" {
		authorID = pr.Author
	} else {
		authorID = nil // NULL в БД
	}

	if _, err := tx.ExecContext(ctx, qPR,
		pr.PullRequestID,
		pr.Name,
		authorID,
		pr.Status,
	); err != nil {
		return fmt.Errorf("save pull request: insert pull_requests: %w", err)
	}

	// 2. save reviewers
	const qReviewer = `
		INSERT INTO pr_reviewers (pr_id, user_id)
		VALUES ($1, $2)
	`

	for _, reviewerID := range pr.Reviewers {
		if reviewerID == "" {
			continue
		}
		if _, err := tx.ExecContext(ctx, qReviewer, pr.PullRequestID, reviewerID); err != nil {
			return fmt.Errorf("save pull request: insert reviewer %s: %w", reviewerID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("save pull request: commit: %w", err)
	}

	return nil
}

// GetByID returns pull request by ID
func (r *Repo) GetByID(ctx context.Context, prID pullrequest.ID) (*pullrequest.PullRequest, error) {

	const q = `
		SELECT pr_id, name, author_id, status
		FROM pull_requests
		WHERE pr_id = $1
	`

	var (
		m        PullRequestModel
		authorID sql.NullString
	)

	row := r.db.QueryRowContext(ctx, q, prID)
	if err := row.Scan(&m.PRID, &m.Name, &authorID, &m.Status); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pullrequest.ErrPullRequestNotFound
		}
		return nil, fmt.Errorf("get pull request by id: scan: %w", err)
	}

	if authorID.Valid {
		m.AuthorID = authorID.String
	}

	reviewers, err := r.loadReviewers(ctx, nil, prID)
	if err != nil {
		return nil, fmt.Errorf("get pull request by id: load reviewers: %w", err)
	}

	return pullRequestModelToDomain(m, reviewers), nil
}

// UpdateStatus changes pull request status
func (r *Repo) UpdateStatus(ctx context.Context, prID pullrequest.ID, status pullrequest.Status) error {
	const q = `
		UPDATE pull_requests
		SET status = $1
		WHERE pr_id = $2
	`

	res, err := r.db.ExecContext(ctx, q, status, prID)
	if err != nil {
		return fmt.Errorf("update status: exec: %w", err)
	}

	// 1. check if the row was in the db
	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("update status: rows affected: %w", err)
	}
	if affected == 0 {
		return pullrequest.ErrPullRequestNotFound
	}

	return nil
}

// AssignReviewers sets reviewers for pull request
func (r *Repo) AssignReviewers(ctx context.Context, prID pullrequest.ID, reviewers []user.ID) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("assign reviewers: begin tx: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	// 1. Check if the pr exists
	const qCheckPR = `
		SELECT 1
		FROM pull_requests
		WHERE pr_id = $1
	`
	var exists int
	if err := tx.QueryRowContext(ctx, qCheckPR, prID).Scan(&exists); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return pullrequest.ErrPullRequestNotFound
		}
		return fmt.Errorf("assign reviewers: check pr exists: %w", err)
	}

	// 2. delete old reviewers
	const qDelete = `
		DELETE FROM pr_reviewers
		WHERE pr_id = $1
	`
	if _, err := tx.ExecContext(ctx, qDelete, prID); err != nil {
		return fmt.Errorf("assign reviewers: delete old reviewers: %w", err)
	}

	// 3. insert new reviewers
	const qInsert = `
		INSERT INTO pr_reviewers (pr_id, user_id)
		VALUES ($1, $2)
	`

	for _, reviewer := range reviewers {
		if reviewer == "" {
			continue
		}
		if _, err := tx.ExecContext(ctx, qInsert, prID, reviewer); err != nil {
			return fmt.Errorf("assign reviewers: insert reviewer %s: %w", reviewer, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("assign reviewers: commit: %w", err)
	}

	return nil
}

// ReassignReviewers replaces one reviewer with another
func (r *Repo) ReassignReviewers(ctx context.Context, prID pullrequest.ID, oldReviewerID, newReviewerID user.ID) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("reassign reviewers: begin tx: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	// 1. Check if the pr exists
	const qCheckPR = `
		SELECT 1
		FROM pull_requests
		WHERE pr_id = $1
	`
	var exists int
	if err := tx.QueryRowContext(ctx, qCheckPR, prID).Scan(&exists); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return pullrequest.ErrPullRequestNotFound
		}
		return fmt.Errorf("reassign reviewers: check pr exists: %w", err)
	}

	// 2. delete old reviewer
	const qDelete = `
		DELETE FROM pr_reviewers
		WHERE pr_id = $1 AND user_id = $2
	`
	if _, err := tx.ExecContext(ctx, qDelete, prID, oldReviewerID); err != nil {
		return fmt.Errorf("reassign reviewers: delete old reviewer: %w", err)
	}

	// 3. isnert new reviewer
	const qInsert = `
		INSERT INTO pr_reviewers (pr_id, user_id)
		VALUES ($1, $2)
		ON CONFLICT (pr_id, user_id) DO NOTHING
	`
	if _, err := tx.ExecContext(ctx, qInsert, prID, newReviewerID); err != nil {
		return fmt.Errorf("reassign reviewers: insert new reviewer: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("reassign reviewers: commit: %w", err)
	}

	return nil
}

// ListByUserID returns all user's pull requests
// (author or reviewer)
func (r *Repo) ListByUserID(ctx context.Context, userID user.ID) ([]*pullrequest.PullRequest, error) {
	const q = `
		SELECT DISTINCT pr.pr_id, pr.name, pr.author_id, pr.status
		FROM pull_requests pr
		LEFT JOIN pr_reviewers prr ON prr.pr_id = pr.pr_id
		WHERE pr.author_id = $1
		   OR prr.user_id = $1
	`

	rows, err := r.db.QueryContext(ctx, q, userID)
	if err != nil {
		return nil, fmt.Errorf("list by user id: query: %w", err)
	}
	defer rows.Close()

	var (
		result []*pullrequest.PullRequest
	)

	for rows.Next() {
		var (
			m        PullRequestModel
			authorID sql.NullString
		)

		if err := rows.Scan(&m.PRID, &m.Name, &authorID, &m.Status); err != nil {
			return nil, fmt.Errorf("list by user id: scan: %w", err)
		}
		if authorID.Valid {
			m.AuthorID = authorID.String
		}

		reviewers, err := r.loadReviewers(ctx, nil, pullrequest.ID(m.PRID))
		if err != nil {
			return nil, fmt.Errorf("list by user id: load reviewers for pr %s: %w", m.PRID, err)
		}

		result = append(result, pullRequestModelToDomain(m, reviewers))
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list by user id: rows err: %w", err)
	}

	return result, nil
}

func (r *Repo) loadReviewers(ctx context.Context, tx *sql.Tx, prID pullrequest.ID) ([]user.ID, error) {
	const q = `
		SELECT user_id
		FROM pr_reviewers
		WHERE pr_id = $1
	`

	var (
		rows *sql.Rows
		err  error
	)

	if tx == nil {
		rows, err = r.db.QueryContext(ctx, q, prID)
	} else {
		rows, err = tx.QueryContext(ctx, q, prID)
	}
	if err != nil {
		return nil, fmt.Errorf("load reviewers: query: %w", err)
	}
	defer rows.Close()

	var reviewers []user.ID
	for rows.Next() {
		var uid user.ID
		if err := rows.Scan(&uid); err != nil {
			return nil, fmt.Errorf("load reviewers: scan: %w", err)
		}
		reviewers = append(reviewers, uid)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("load reviewers: rows err: %w", err)
	}

	return reviewers, nil
}
