package policy

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	dbgen "github.com/kyleaupton/snaggle/backend/internal/db/sqlc"
	"github.com/kyleaupton/snaggle/backend/internal/logger"
	"github.com/kyleaupton/snaggle/backend/internal/model"
	"github.com/kyleaupton/snaggle/backend/internal/quality"
	"github.com/kyleaupton/snaggle/backend/internal/repo"
)

// Engine evaluates policies against torrent metadata to produce plans
type Engine struct {
	repo   *repo.Repository
	logger *logger.Logger
}

func NewEngine(r *repo.Repository, l *logger.Logger) *Engine {
	return &Engine{repo: r, logger: l}
}

// Evaluate evaluates all enabled policies in priority order and returns an EvaluationTrace
func (e *Engine) Evaluate(ctx context.Context, candidate model.DownloadCandidate) (model.EvaluationTrace, error) {
	// print candidate
	jsonData, _ := json.MarshalIndent(candidate, "", "  ")
	e.logger.Debug().Msgf("candidate: %+v", string(jsonData))

	trace := model.EvaluationTrace{
		Policies: []model.PolicyEvaluation{},
		FinalPlan: model.Plan{
			DownloaderID:   "",
			LibraryID:      "",
			NameTemplateID: "",
		},
	}

	policies, err := e.repo.ListPolicies(ctx)
	if err != nil {
		return trace, fmt.Errorf("list policies: %w", err)
	}

	// Filter to only enabled policies
	var enabledPolicies []dbgen.Policy
	for _, p := range policies {
		if p.Enabled {
			enabledPolicies = append(enabledPolicies, p)
		}
	}

	// Evaluate policies in priority order (already sorted DESC by query)
	for _, policy := range enabledPolicies {
		policyEval := model.PolicyEvaluation{
			PolicyID:          policy.ID.String(),
			PolicyName:        policy.Name,
			Priority:          int(policy.Priority),
			Matched:           false,
			ActionsApplied:    []model.ActionInfo{},
			StoppedProcessing: false,
		}

		rule, err := e.repo.GetRuleForPolicy(ctx, policy.ID)
		if err != nil {
			// Policy without a rule doesn't match
			policyEval.RuleEvaluated = nil
			trace.Policies = append(trace.Policies, policyEval)
			continue
		}

		// Store rule info
		policyEval.RuleEvaluated = &model.RuleInfo{
			LeftOperand:  rule.LeftOperand,
			Operator:     rule.Operator,
			RightOperand: rule.RightOperand,
		}

		candidateContext := model.CandidateContext{
			Candidate: candidate,
			Quality:   quality.NewParser().Parse(candidate.Title),
		}

		// Evaluate rule
		matches, err := e.evaluateRule(ctx, rule, candidateContext)
		if err != nil {
			return trace, fmt.Errorf("evaluate rule for policy %s: %w", policy.ID.String(), err)
		}

		if !matches {
			trace.Policies = append(trace.Policies, policyEval)
			continue
		}

		// Policy matched!
		policyEval.Matched = true

		// Get actions for this policy
		actions, err := e.repo.ListActionsForPolicy(ctx, policy.ID)
		if err != nil {
			return trace, fmt.Errorf("list actions for policy %s: %w", policy.ID.String(), err)
		}

		// Apply actions in order and track them
		for _, action := range actions {
			actionInfo := model.ActionInfo{
				Type:  action.Type,
				Value: action.Value,
				Order: action.Order,
			}
			policyEval.ActionsApplied = append(policyEval.ActionsApplied, actionInfo)

			if err := e.applyAction(&trace.FinalPlan, action); err != nil {
				return trace, fmt.Errorf("apply action %s: %w", action.ID.String(), err)
			}

			// Stop processing if stop_processing action
			if action.Type == string(model.ActionStopProcessing) {
				policyEval.StoppedProcessing = true
				trace.Policies = append(trace.Policies, policyEval)
				return trace, nil
			}
		}

		trace.Policies = append(trace.Policies, policyEval)
	}

	// Default logic
	// If we get to this point and there are decisions that are not made, we should
	// attempt to set the missing decisions to the default values. If the user does
	// not have default items set then we should return an error.
	if trace.FinalPlan.DownloaderID == "" {
		downloader, err := e.repo.GetDefaultDownloader(ctx, "torrent")
		if err != nil {
			return trace, fmt.Errorf("get default downloader: %w", err)
		}
		trace.FinalPlan.DownloaderID = downloader.ID.String()
	}

	mediaType, err := candidate.GetMediaType()
	if err != nil {
		return trace, fmt.Errorf("get media type: %w", err)
	}

	if trace.FinalPlan.LibraryID == "" {
		library, err := e.repo.GetDefaultLibrary(ctx, string(mediaType))
		if err != nil {
			return trace, fmt.Errorf("get default library: %w", err)
		}
		trace.FinalPlan.LibraryID = library.ID.String()
	}

	if trace.FinalPlan.NameTemplateID == "" {
		nameTemplate, err := e.repo.GetDefaultNameTemplate(ctx, string(mediaType))
		if err != nil {
			return trace, fmt.Errorf("get default name template: %w", err)
		}
		trace.FinalPlan.NameTemplateID = nameTemplate.ID.String()
	}

	return trace, nil
}

// evaluateRule evaluates a rule against torrent metadata
func (e *Engine) evaluateRule(ctx context.Context, rule dbgen.Rule, candidateContext model.CandidateContext) (bool, error) {
	operator := model.Operator(rule.Operator)

	// Handle logical operators (and, or, not) which reference other rules
	switch operator {
	case model.OpAnd:
		return e.evaluateLogicalAnd(ctx, rule, candidateContext)
	case model.OpOr:
		return e.evaluateLogicalOr(ctx, rule, candidateContext)
	case model.OpNot:
		return e.evaluateLogicalNot(ctx, rule, candidateContext)
	}

	// Evaluate left operand
	leftVal, err := e.getValue(rule.LeftOperand, candidateContext)
	if err != nil {
		return false, fmt.Errorf("get left value: %w", err)
	}

	// Evaluate right operand
	rightVal, err := e.getValue(rule.RightOperand, candidateContext)
	if err != nil {
		return false, fmt.Errorf("get right value: %w", err)
	}

	// Compare based on operator
	return e.compare(leftVal, operator, rightVal)
}

// evaluateLogicalAnd evaluates an AND rule (left and right are rule UUIDs)
func (e *Engine) evaluateLogicalAnd(ctx context.Context, rule dbgen.Rule, candidateContext model.CandidateContext) (bool, error) {
	leftRule, err := e.getRuleByID(ctx, rule.LeftOperand)
	if err != nil {
		return false, err
	}
	rightRule, err := e.getRuleByID(ctx, rule.RightOperand)
	if err != nil {
		return false, err
	}

	leftResult, err := e.evaluateRule(ctx, leftRule, candidateContext)
	if err != nil {
		return false, err
	}
	rightResult, err := e.evaluateRule(ctx, rightRule, candidateContext)
	if err != nil {
		return false, err
	}

	return leftResult && rightResult, nil
}

// evaluateLogicalOr evaluates an OR rule (left and right are rule UUIDs)
func (e *Engine) evaluateLogicalOr(ctx context.Context, rule dbgen.Rule, candidateContext model.CandidateContext) (bool, error) {
	leftRule, err := e.getRuleByID(ctx, rule.LeftOperand)
	if err != nil {
		return false, err
	}
	rightRule, err := e.getRuleByID(ctx, rule.RightOperand)
	if err != nil {
		return false, err
	}

	leftResult, err := e.evaluateRule(ctx, leftRule, candidateContext)
	if err != nil {
		return false, err
	}
	rightResult, err := e.evaluateRule(ctx, rightRule, candidateContext)
	if err != nil {
		return false, err
	}

	return leftResult || rightResult, nil
}

// evaluateLogicalNot evaluates a NOT rule (right is a rule UUID)
func (e *Engine) evaluateLogicalNot(ctx context.Context, rule dbgen.Rule, candidateContext model.CandidateContext) (bool, error) {
	rightRule, err := e.getRuleByID(ctx, rule.RightOperand)
	if err != nil {
		return false, err
	}

	result, err := e.evaluateRule(ctx, rightRule, candidateContext)
	if err != nil {
		return false, err
	}

	return !result, nil
}

// getRuleByID gets a rule by UUID string
func (e *Engine) getRuleByID(ctx context.Context, ruleIDStr string) (dbgen.Rule, error) {
	ruleID, err := uuid.Parse(ruleIDStr)
	if err != nil {
		return dbgen.Rule{}, fmt.Errorf("invalid rule UUID: %w", err)
	}

	// We need to get the rule, but we only have policy_id in our queries
	// For now, we'll need to get all rules and find the one we need
	// This is inefficient but works for now - can be optimized later
	policies, err := e.repo.ListPolicies(ctx)
	if err != nil {
		return dbgen.Rule{}, err
	}

	for _, policy := range policies {
		rule, err := e.repo.GetRuleForPolicy(ctx, policy.ID)
		if err != nil {
			continue
		}
		ruleUUID := pgtype.UUID{Bytes: ruleID, Valid: true}
		if rule.ID == ruleUUID {
			return rule, nil
		}
	}

	return dbgen.Rule{}, fmt.Errorf("rule not found: %s", ruleIDStr)
}

// getValue gets a value from metadata or returns the literal value
func (e *Engine) getValue(operand string, candidateContext model.CandidateContext) (interface{}, error) {
	// Check if it's a field reference (candidate.*)
	if strings.HasPrefix(operand, "candidate.") {
		field := strings.TrimPrefix(operand, "candidate.")
		return e.getCandidateField(field, candidateContext)
	}

	// Check if it's a field reference (quality.*)
	if strings.HasPrefix(operand, "quality.") {
		field := strings.TrimPrefix(operand, "quality.")
		return e.getQualityField(field, candidateContext)
	}

	// Backward compatibility: support torrent.* (deprecated)
	if strings.HasPrefix(operand, "torrent.") {
		e.logger.Warn().Str("field", operand).Msg("Using deprecated torrent.* field, please migrate to candidate.*")
		field := strings.TrimPrefix(operand, "torrent.")
		return e.getCandidateField(field, candidateContext)
	}

	// Try to parse as number
	if num, err := strconv.ParseInt(operand, 10, 64); err == nil {
		return num, nil
	}
	if num, err := strconv.ParseFloat(operand, 64); err == nil {
		return num, nil
	}

	// Return as string
	return operand, nil
}

// getCandidateField gets a value from the candidate
func (e *Engine) getCandidateField(field string, candidateContext model.CandidateContext) (interface{}, error) {
	candidate := candidateContext.Candidate
	switch field {
	case "size":
		return int64(candidate.Size), nil
	case "title":
		return candidate.Title, nil
	case "indexer":
		return candidate.Indexer, nil
	case "indexer_id":
		return int64(candidate.IndexerID), nil
	case "categories":
		return candidate.Categories, nil
	case "protocol":
		return candidate.Protocol, nil
	case "torrent_seeders":
		if candidate.Protocol == "torrent" {
			return int64(candidate.Seeders), nil
		}
		return nil, fmt.Errorf("field torrent_seeders only available for torrent protocol")
	case "torrent_peers":
		if candidate.Protocol == "torrent" {
			return int64(candidate.Peers), nil
		}
		return nil, fmt.Errorf("field torrent_peers only available for torrent protocol")
	// Backward compatibility mappings
	case "seeders":
		return int64(candidate.Seeders), nil
	case "peers":
		return int64(candidate.Peers), nil
	case "tracker":
		return candidate.Indexer, nil
	case "tracker_id":
		return int64(candidate.IndexerID), nil
	default:
		return nil, fmt.Errorf("unknown candidate field: %s", field)
	}
}

// getQualityField gets a value from the parsed quality
func (e *Engine) getQualityField(field string, candidateContext model.CandidateContext) (interface{}, error) {
	quality := candidateContext.Quality
	
	// Handle nested audio fields
	if strings.HasPrefix(field, "audio.") {
		audioField := strings.TrimPrefix(field, "audio.")
		switch audioField {
		case "codec":
			return string(quality.Audio.Codec), nil
		case "channels":
			return string(quality.Audio.Channels), nil
		default:
			return nil, fmt.Errorf("unknown quality.audio field: %s", audioField)
		}
	}
	
	switch field {
	case "resolution":
		return string(quality.Resolution), nil
	case "source":
		return string(quality.Source), nil
	case "codec":
		return string(quality.Codec), nil
	case "container":
		return string(quality.Container), nil
	case "hdr":
		return string(quality.HDR), nil
	case "bit_depth":
		return string(quality.BitDepth), nil
	case "tier":
		return string(quality.Tier), nil
	case "is_remux":
		return quality.IsRemux, nil
	case "is_proper":
		return quality.IsProper, nil
	case "is_repack":
		return quality.IsRepack, nil
	case "is_extended":
		return quality.IsExtended, nil
	case "release_group":
		return quality.ReleaseGroup, nil
	default:
		return nil, fmt.Errorf("unknown quality field: %s", field)
	}
}

// compare compares two values based on operator
func (e *Engine) compare(left interface{}, operator model.Operator, right interface{}) (bool, error) {
	switch operator {
	case model.OpEq:
		return e.equals(left, right), nil
	case model.OpNe:
		return !e.equals(left, right), nil
	case model.OpGt:
		return e.greaterThan(left, right)
	case model.OpGte:
		gt, err := e.greaterThan(left, right)
		if err != nil {
			return false, err
		}
		eq := e.equals(left, right)
		return gt || eq, nil
	case model.OpLt:
		return e.lessThan(left, right)
	case model.OpLte:
		lt, err := e.lessThan(left, right)
		if err != nil {
			return false, err
		}
		eq := e.equals(left, right)
		return lt || eq, nil
	case model.OpContains:
		return e.contains(left, right)
	case model.OpIn:
		return e.in(left, right)
	case model.OpNotIn:
		result, err := e.in(left, right)
		return !result, err
	default:
		return false, fmt.Errorf("unsupported operator: %s", operator)
	}
}

func (e *Engine) equals(left, right interface{}) bool {
	return fmt.Sprintf("%v", left) == fmt.Sprintf("%v", right)
}

func (e *Engine) greaterThan(left, right interface{}) (bool, error) {
	leftNum, rightNum, err := e.toNumbers(left, right)
	if err != nil {
		return false, err
	}
	return leftNum > rightNum, nil
}

func (e *Engine) lessThan(left, right interface{}) (bool, error) {
	leftNum, rightNum, err := e.toNumbers(left, right)
	if err != nil {
		return false, err
	}
	return leftNum < rightNum, nil
}

func (e *Engine) toNumbers(left, right interface{}) (float64, float64, error) {
	leftNum, ok := e.toFloat64(left)
	if !ok {
		return 0, 0, fmt.Errorf("left operand is not a number: %v", left)
	}
	rightNum, ok := e.toFloat64(right)
	if !ok {
		return 0, 0, fmt.Errorf("right operand is not a number: %v", right)
	}
	return leftNum, rightNum, nil
}

func (e *Engine) toFloat64(v interface{}) (float64, bool) {
	switch val := v.(type) {
	case int64:
		return float64(val), true
	case float64:
		return val, true
	case int:
		return float64(val), true
	case uint:
		return float64(val), true
	case uint64:
		return float64(val), true
	default:
		return 0, false
	}
}

func (e *Engine) contains(left, right interface{}) (bool, error) {
	leftStr := fmt.Sprintf("%v", left)
	rightStr := fmt.Sprintf("%v", right)
	return strings.Contains(leftStr, rightStr), nil
}

func (e *Engine) in(left, right interface{}) (bool, error) {
	leftStr := fmt.Sprintf("%v", left)

	// Right should be a comma-separated list or array
	rightStr := fmt.Sprintf("%v", right)
	values := strings.Split(rightStr, ",")
	for _, v := range values {
		if strings.TrimSpace(v) == leftStr {
			return true, nil
		}
	}

	// Also check if right is a slice
	if categories, ok := right.([]string); ok {
		for _, cat := range categories {
			if cat == leftStr {
				return true, nil
			}
		}
	}

	return false, nil
}

// applyAction applies an action to the plan
func (e *Engine) applyAction(plan *model.Plan, action dbgen.Action) error {
	actionType := model.ActionType(action.Type)

	switch actionType {
	case model.ActionSetDownloader:
		plan.DownloaderID = action.Value
	case model.ActionSetLibrary:
		plan.LibraryID = action.Value
	case model.ActionSetNameTemplate:
		plan.NameTemplateID = action.Value
	case model.ActionStopProcessing:
		// Handled in Evaluate
	default:
		return fmt.Errorf("unknown action type: %s", actionType)
	}

	return nil
}
