package transactor

import (
	"context"
	"database/sql"
	"fmt"
)

type txKey struct{}

// Transactor управляет транзакциями
type Transactor struct {
	db *sql.DB
}

func NewTransactor(db *sql.DB) *Transactor {
	return &Transactor{db: db}
}

// WithinTransaction выполняет функцию tFunc внутри транзакции.
func (t *Transactor) WithinTransaction(ctx context.Context, tFunc func(ctx context.Context) error) error {
	// Проверяем, есть ли tx внутри ctx
	if _, ok := ctx.Value(txKey{}).(*sql.Tx); ok {
		return tFunc(ctx)
	}

	tx, err := t.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	txCtx := context.WithValue(ctx, txKey{}, tx)

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	if err := tFunc(txCtx); err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

// SQLQueryable интерфейс, объединяющий *sql.DB и *sql.Tx
type SQLQueryable interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

// ExtractReq достает *sql.DB/*sql.Tx
func (t *Transactor) ExtractReq(ctx context.Context) SQLQueryable {
	if tx, ok := ctx.Value(txKey{}).(*sql.Tx); ok {
		return tx
	}
	return t.db
}
