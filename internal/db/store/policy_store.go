package store

import (
	"context"
	"errors"
	"fmt"
	"rps/internal/db/sqlc"
	apperror "rps/internal/errors"

	"github.com/jackc/pgx/v5"
)

type PolicyStore struct {
	*sqlc.Queries
}

func NewPolicyStore(queries *sqlc.Queries) *PolicyStore {
	return &PolicyStore{
		queries,
	}
}

func (p *PolicyStore) GetPolicyByResourceTypeAndId(ctx context.Context, resourceType string, resourceID string) (*sqlc.Policy, error) {
	params := sqlc.GetPolicyByResourceParams{ResourceType: resourceType, ResourceID: resourceID}
	policyDb, err := p.GetPolicyByResource(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.NotFound(fmt.Sprintf("policy [%s/%s]", resourceType, resourceID))
		}
		return nil, err
	}
	return &policyDb, nil
}

func (p *PolicyStore) CreatePolicyWithParam(ctx context.Context, resourceType string, resourceID string, version string) (*sqlc.Policy, error) {
	params := sqlc.CreatePolicyParams{ResourceType: resourceType, ResourceID: resourceID, Version: version}
	policyDb, err := p.CreatePolicy(ctx, params)
	if err != nil {
		return nil, err
	}
	return &policyDb, nil
}

func (p *PolicyStore) ExistsPolicyStore(ctx context.Context, resourceType string, resourceID string) (bool, error) {
	return p.ExistsPolicy(ctx, sqlc.ExistsPolicyParams{ResourceType: resourceType, ResourceID: resourceID})
}

func (p *PolicyStore) UpdatePolicyStore(ctx context.Context, id int64, version string) error {
	_, err := p.UpdatePolicy(ctx, sqlc.UpdatePolicyParams{
		ID:      id,
		Version: version,
	})
	return err
}

func (p *PolicyStore) DeletePolicyByResourceCustom(ctx context.Context, resourceType string, resourceID string) error {
	if err := p.DeletePolicyByResource(ctx, sqlc.DeletePolicyByResourceParams{
		ResourceType: resourceType,
		ResourceID:   resourceID,
	}); err != nil {
		return err
	}
	return nil
}

func (p *PolicyStore) ListPoliciesWithParam(ctx context.Context, resourceType string, resourceID string, offset int32, limit int32) ([]sqlc.Policy, error) {
	return p.ListPolicies(ctx, sqlc.ListPoliciesParams{
		ResourceType: resourceType,
		ResourceID:   resourceID,
		Off:          offset,
		Lim:          limit,
	})
}

func (p *PolicyStore) CountPoliciesWithParam(ctx context.Context, resourceType string, resourceID string) (int64, error) {
	return p.CountPolicies(ctx, sqlc.CountPoliciesParams{
		ResourceType: resourceType,
		ResourceID:   resourceID,
	})
}
