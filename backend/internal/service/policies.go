package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgtype"
	dbgen "github.com/kyleaupton/snaggle/backend/internal/db/sqlc"
	"github.com/kyleaupton/snaggle/backend/internal/logger"
	"github.com/kyleaupton/snaggle/backend/internal/model"
	"github.com/kyleaupton/snaggle/backend/internal/policy"
	"github.com/kyleaupton/snaggle/backend/internal/repo"
)

type PoliciesService struct {
	repo   *repo.Repository
	engine *policy.Engine
}

func NewPoliciesService(r *repo.Repository, logg *logger.Logger) *PoliciesService {
	return &PoliciesService{
		repo:   r,
		engine: policy.NewEngine(r, logg),
	}
}

func (s *PoliciesService) List(ctx context.Context) ([]dbgen.Policy, error) {
	return s.repo.ListPolicies(ctx)
}

func (s *PoliciesService) Get(ctx context.Context, id pgtype.UUID) (dbgen.Policy, error) {
	return s.repo.GetPolicy(ctx, id)
}

func (s *PoliciesService) Create(ctx context.Context, name string, description *string, enabled bool, priority int32) (dbgen.Policy, error) {
	if name == "" {
		return dbgen.Policy{}, errors.New("name required")
	}
	return s.repo.CreatePolicy(ctx, name, description, enabled, priority)
}

func (s *PoliciesService) Update(ctx context.Context, id pgtype.UUID, name string, description *string, enabled bool, priority int32) (dbgen.Policy, error) {
	if name == "" {
		return dbgen.Policy{}, errors.New("name required")
	}
	return s.repo.UpdatePolicy(ctx, id, name, description, enabled, priority)
}

func (s *PoliciesService) Delete(ctx context.Context, id pgtype.UUID) error {
	return s.repo.DeletePolicy(ctx, id)
}

func (s *PoliciesService) GetRule(ctx context.Context, policyID pgtype.UUID) (dbgen.Rule, error) {
	return s.repo.GetRuleForPolicy(ctx, policyID)
}

func (s *PoliciesService) CreateRule(ctx context.Context, policyID pgtype.UUID, leftOperand, operator, rightOperand string) (dbgen.Rule, error) {
	// Validate operator
	validOps := []string{"==", "!=", ">", ">=", "<", "<=", "contains", "in", "not in", "and", "or", "not"}
	valid := false
	for _, op := range validOps {
		if operator == op {
			valid = true
			break
		}
	}
	if !valid {
		return dbgen.Rule{}, errors.New("invalid operator")
	}

	return s.repo.CreateRule(ctx, policyID, leftOperand, operator, rightOperand)
}

func (s *PoliciesService) UpdateRule(ctx context.Context, id pgtype.UUID, leftOperand, operator, rightOperand string) (dbgen.Rule, error) {
	// Validate operator
	validOps := []string{"==", "!=", ">", ">=", "<", "<=", "contains", "in", "not in", "and", "or", "not"}
	valid := false
	for _, op := range validOps {
		if operator == op {
			valid = true
			break
		}
	}
	if !valid {
		return dbgen.Rule{}, errors.New("invalid operator")
	}

	return s.repo.UpdateRule(ctx, id, leftOperand, operator, rightOperand)
}

func (s *PoliciesService) DeleteRule(ctx context.Context, id pgtype.UUID) error {
	return s.repo.DeleteRule(ctx, id)
}

func (s *PoliciesService) ListActions(ctx context.Context, policyID pgtype.UUID) ([]dbgen.Action, error) {
	return s.repo.ListActionsForPolicy(ctx, policyID)
}

func (s *PoliciesService) GetAction(ctx context.Context, id pgtype.UUID) (dbgen.Action, error) {
	return s.repo.GetAction(ctx, id)
}

func (s *PoliciesService) CreateAction(ctx context.Context, policyID pgtype.UUID, actionType, value string, order int32) (dbgen.Action, error) {
	// Validate action type
	validTypes := []string{"set_downloader", "set_library", "set_name_template", "stop_processing"}
	valid := false
	for _, t := range validTypes {
		if actionType == t {
			valid = true
			break
		}
	}
	if !valid {
		return dbgen.Action{}, errors.New("invalid action type")
	}

	if value == "" && actionType != "stop_processing" {
		return dbgen.Action{}, errors.New("value required for action type")
	}

	return s.repo.CreateAction(ctx, policyID, actionType, value, order)
}

func (s *PoliciesService) UpdateAction(ctx context.Context, id pgtype.UUID, actionType, value string, order int32) (dbgen.Action, error) {
	// Validate action type
	validTypes := []string{"set_downloader", "set_library", "set_name_template", "stop_processing"}
	valid := false
	for _, t := range validTypes {
		if actionType == t {
			valid = true
			break
		}
	}
	if !valid {
		return dbgen.Action{}, errors.New("invalid action type")
	}

	if value == "" && actionType != "stop_processing" {
		return dbgen.Action{}, errors.New("value required for action type")
	}

	return s.repo.UpdateAction(ctx, id, actionType, value, order)
}

func (s *PoliciesService) DeleteAction(ctx context.Context, id pgtype.UUID) error {
	return s.repo.DeleteAction(ctx, id)
}

// Evaluate evaluates policies against torrent metadata and returns an EvaluationTrace
func (s *PoliciesService) Evaluate(ctx context.Context, candidate model.DownloadCandidate) (model.EvaluationTrace, error) {
	return s.engine.Evaluate(ctx, candidate)
}
