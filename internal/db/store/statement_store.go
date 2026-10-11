package store

import (
	"context"
	"rps/internal/db/helper"
	"rps/internal/db/sqlc"
)

type StatementStore struct {
	*sqlc.Queries
}

func NewStatementStore(q *sqlc.Queries) *StatementStore {
	return &StatementStore{
		q,
	}
}

func (s *StatementStore) CreateStatementCustom(ctx context.Context, policyID int64, sid string, effect string) (*sqlc.Statement, error) {
	var param = sqlc.CreateStatementParams{
		PolicyID: policyID,
		Sid:      helper.ToPgText(sid),
		Effect:   effect,
	}
	st, err := s.CreateStatement(ctx, param)
	if err != nil {
		return nil, err
	}
	return &st, nil
}
