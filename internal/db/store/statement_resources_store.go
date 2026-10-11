package store

import (
	"context"
	"rps/internal/db/sqlc"
)

type StatementResourceStore struct {
	*sqlc.Queries
}

func NewStatementResourceStore(q *sqlc.Queries) *StatementResourceStore {
	return &StatementResourceStore{
		q,
	}
}

func (s StatementResourceStore) CreateStatementResourceCustom(ctx context.Context, statementID int64, resource string) error {
	param := sqlc.CreateStatementResourceParams{
		StatementID: statementID,
		Resource:    resource,
	}
	_, err := s.CreateStatementResource(ctx, param)
	if err != nil {
		return err
	}
	return nil
}
