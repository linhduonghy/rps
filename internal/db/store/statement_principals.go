package store

import (
	"context"
	"rps/internal/db/sqlc"
)

type StatementPrincipalStore struct {
	*sqlc.Queries
}

func NewStatementPrincipalStore(q *sqlc.Queries) *StatementPrincipalStore {
	return &StatementPrincipalStore{
		q,
	}
}

func (s StatementPrincipalStore) CreateStatementPrincipalCustom(ctx context.Context, statementID int64, principalType string,
	principalID string) error {
	param := sqlc.CreateStatementPrincipalParams{
		StatementID:   statementID,
		PrincipalType: principalType,
		PrincipalID:   principalID,
	}
	_, err := s.CreateStatementPrincipal(ctx, param)
	if err != nil {
		return err
	}
	return nil
}
