package store

import (
	"context"
	"rps/internal/db/sqlc"
)

type StatementActionStore struct {
	*sqlc.Queries
}

func NewStatementActionStore(q *sqlc.Queries) *StatementActionStore {
	return &StatementActionStore{
		q,
	}
}

func (s StatementActionStore) CreateStatementActionCustom(ctx context.Context, statementID int64, action string) error {
	param := sqlc.CreateStatementActionParams{
		StatementID: statementID,
		Action:      action,
	}
	if _, err := s.CreateStatementAction(ctx, param); err != nil {
		return err
	}
	return nil
}
