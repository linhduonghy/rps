package store

import (
	"context"
	"rps/internal/db/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	db                 *pgxpool.Pool
	Policy             *PolicyStore
	Statement          *StatementStore
	StatementAction    *StatementActionStore
	StatementResource  *StatementResourceStore
	StatementPrincipal *StatementPrincipalStore
}

func New(db *pgxpool.Pool) *Store {
	return newStore(db, sqlc.New(db))
}

func newStore(db *pgxpool.Pool, q *sqlc.Queries) *Store {
	return &Store{
		db:                 db,
		Policy:             NewPolicyStore(q),
		Statement:          NewStatementStore(q),
		StatementAction:    NewStatementActionStore(q),
		StatementResource:  NewStatementResourceStore(q),
		StatementPrincipal: NewStatementPrincipalStore(q),
	}
}

func (s *Store) DB() *pgxpool.Pool {
	return s.db
}

func (s *Store) WithTx(
	ctx context.Context,
	fn func(store *Store) error,
) error {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	txStore := newStore(s.db, sqlc.New(tx))

	if err = fn(txStore); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
