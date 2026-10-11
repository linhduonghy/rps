package service

import (
	"context"
	"fmt"
	"log/slog"
	"rps/internal/db/helper"
	"rps/internal/db/store"
	"rps/internal/domain"
	apperror "rps/internal/errors"
	"rps/internal/http/request"
)

type policyService struct {
	store *store.Store
	log   *slog.Logger
}

func NewPolicyService(store *store.Store, log *slog.Logger) *policyService {
	return &policyService{
		store: store,
		log:   log,
	}
}

func (s *policyService) CreatePolicy(ctx context.Context, request *request.CreatePolicy) error {
	s.log.Info("create policy", "payload", request)
	return s.store.WithTx(ctx, func(txStore *store.Store) error {
		// validate duplication
		exist, err := txStore.Policy.ExistsPolicyStore(ctx, request.ResourceType, request.ResourceID)
		if err != nil {
			return err
		}
		if exist {
			return apperror.AlreadyExists(fmt.Sprintf("[%s/%s]", request.ResourceType, request.ResourceID))
		}

		policy, err := txStore.Policy.CreatePolicyWithParam(ctx, request.ResourceType, request.ResourceID, request.Version)
		if err != nil {
			return err
		}
		if err = createStatements(ctx, txStore, policy.ID, request.Statements); err != nil {
			return err
		}
		return nil
	})
}

func createStatements(ctx context.Context, txStore *store.Store, policyID int64, statements []request.Statement) error {
	for _, statement := range statements {
		if err := createStatement(ctx, txStore, policyID, statement); err != nil {
			return err
		}
	}
	return nil
}

func createStatement(ctx context.Context, txStore *store.Store, policyID int64, statement request.Statement) error {
	stDB, err := txStore.Statement.CreateStatementCustom(ctx, policyID, statement.SID, statement.Effect)
	if err != nil {
		return err
	}
	//
	if err = createStatementActions(ctx, txStore, stDB.ID, statement.Actions); err != nil {
		return err
	}
	if err = createStatementResources(ctx, txStore, stDB.ID, statement.Resources); err != nil {
		return err
	}
	if err = createStatementPrincipals(ctx, txStore, stDB.ID, statement.Principals); err != nil {
		return err
	}

	return nil
}

func createStatementActions(ctx context.Context, store *store.Store, statementID int64, actions []string) error {
	for _, action := range actions {
		if err := store.StatementAction.CreateStatementActionCustom(ctx, statementID, action); err != nil {
			return err
		}
	}
	return nil
}

func createStatementResources(ctx context.Context, store *store.Store, statementID int64, resources []string) error {
	for _, resource := range resources {
		if err := store.StatementResource.CreateStatementResourceCustom(ctx, statementID, resource); err != nil {
			return err
		}
	}
	return nil
}

func createStatementPrincipals(ctx context.Context, store *store.Store, statementID int64, principals []request.Principal) error {
	for _, principal := range principals {
		if err := store.StatementPrincipal.CreateStatementPrincipalCustom(ctx, statementID, principal.Type, principal.ID); err != nil {
			return err
		}
	}
	return nil
}

func (s *policyService) UpdatePolicy(ctx context.Context, resourceType string, resourceID string, request *request.UpdatePolicy) error {
	s.log.Info("update policy", "payload", request)
	return s.store.WithTx(ctx, func(txStore *store.Store) error {
		// check exist policy
		policyDB, err := s.store.Policy.GetPolicyByResourceTypeAndId(ctx, resourceType, resourceID)
		if err != nil {
			return err
		}
		if policyDB == nil {
			return apperror.NotFound(fmt.Sprintf("[%s/%s]", resourceType, resourceID))
		}

		// update policy metadata
		if err = txStore.Policy.UpdatePolicyStore(ctx, policyDB.ID, request.Version); err != nil {
			return err
		}

		// delete old statements
		if err = txStore.Statement.DeleteStatementsByPolicy(ctx, policyDB.ID); err != nil {
			return err
		}

		// create new statements
		if err = createStatements(ctx, txStore, policyDB.ID, request.Statements); err != nil {
			return err
		}

		return nil
	})
}

func (s *policyService) GetPolicy(
	ctx context.Context,
	resourceType string,
	resourceID string,
) (*domain.Policy, error) {
	policyDb, err := s.store.Policy.GetPolicyByResourceTypeAndId(ctx, resourceType, resourceID)
	if err != nil {
		return nil, err
	}

	policy := &domain.Policy{
		ID:           policyDb.ID,
		ResourceType: policyDb.ResourceType,
		ResourceID:   policyDb.ResourceID,
		Version:      policyDb.Version,
		CreatedAt:    helper.ToTime(policyDb.CreatedAt),
		UpdatedAt:    helper.ToTime(policyDb.UpdatedAt),
	}

	sts, err := s.loadPolicyStatements(ctx, policyDb.ID)
	if err != nil {
		return nil, err
	}
	policy.Statements = sts

	return policy, nil
}

// Build Statements block
func (s *policyService) loadPolicyStatements(ctx context.Context, policyId int64) ([]domain.Statement, error) {
	statements, err := s.store.Statement.ListStatementsByPolicy(ctx, policyId)
	if err != nil {
		return nil, err
	}

	// TODO : compute in parallel to get data
	actionsByStatement, err := s.buildActionsByStatementMap(ctx, policyId)
	if err != nil {
		return nil, err
	}
	resourcesByStatement, err := s.buildResourcesByStatementMap(ctx, policyId)
	if err != nil {
		return nil, err
	}
	principalsByStatement, err := s.buildPrincipalsByStatementMap(ctx, policyId)
	if err != nil {
		return nil, err
	}

	sts := make([]domain.Statement, 0, len(statements))
	for _, statement := range statements {
		sts = append(sts, domain.Statement{
			ID:         statement.ID,
			SID:        statement.Sid.String,
			Effect:     domain.ToEffect(statement.Effect),
			Actions:    actionsByStatement[statement.ID],
			Resources:  resourcesByStatement[statement.ID],
			Principals: principalsByStatement[statement.ID],
			CreatedAt:  helper.ToTime(statement.CreatedAt),
		})
	}
	return sts, nil
}

func (s *policyService) buildActionsByStatementMap(ctx context.Context, policyId int64) (map[int64][]string, error) {
	actions, err := s.store.StatementAction.ListStatementActionsByPolicy(ctx, policyId)
	if err != nil {
		return nil, err
	}
	actionsByStatement := map[int64][]string{}
	for _, action := range actions {
		actionsByStatement[action.StatementID] = append(actionsByStatement[action.StatementID], action.Action)
	}
	return actionsByStatement, nil
}

func (s *policyService) buildResourcesByStatementMap(ctx context.Context, policyId int64) (map[int64][]string, error) {
	resources, err := s.store.StatementResource.ListStatementResourcesByPolicy(ctx, policyId)
	if err != nil {
		return nil, err
	}
	resourcesByStatement := map[int64][]string{}
	for _, resource := range resources {
		resourcesByStatement[resource.StatementID] = append(resourcesByStatement[resource.StatementID], resource.Resource)
	}
	return resourcesByStatement, nil
}

func (s *policyService) buildPrincipalsByStatementMap(ctx context.Context, policyId int64) (map[int64][]domain.Principal, error) {
	principals, err := s.store.StatementPrincipal.ListStatementPrincipalsByPolicy(ctx, policyId)
	if err != nil {
		return nil, err
	}
	principalsByStatement := map[int64][]domain.Principal{}
	for _, principal := range principals {
		principalsByStatement[principal.StatementID] = append(principalsByStatement[principal.StatementID], domain.Principal{
			Type: principal.PrincipalType,
			ID:   principal.PrincipalID,
		})
	}
	return principalsByStatement, nil
}

func (s *policyService) ListPolicies(ctx context.Context, req *request.ListPolicies) ([]*domain.Policy, int64, error) {
	policiesDB, err := s.store.Policy.ListPoliciesWithParam(ctx, req.ResourceType, req.ResourceID, req.GetOffset(), req.GetLimit())
	if err != nil {
		return nil, 0, err
	}

	count, err := s.store.Policy.CountPoliciesWithParam(ctx, req.ResourceType, req.ResourceID)
	if err != nil {
		return nil, 0, err
	}

	policies := make([]*domain.Policy, len(policiesDB))
	for i, policy := range policiesDB {
		policies[i] = &domain.Policy{
			ID:           policy.ID,
			ResourceType: policy.ResourceType,
			ResourceID:   policy.ResourceID,
			Version:      policy.Version,
			CreatedAt:    helper.ToTime(policy.CreatedAt),
			UpdatedAt:    helper.ToTime(policy.UpdatedAt),
		}
	}
	return policies, count, nil
}

func (s *policyService) DeletePolicy(
	ctx context.Context,
	resourceType string,
	resourceID string,
) error {
	if err := s.store.Policy.DeletePolicyByResourceCustom(ctx, resourceType, resourceID); err != nil {
		return err
	}
	return nil
}
