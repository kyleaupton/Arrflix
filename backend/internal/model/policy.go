package model

import (
	"github.com/google/uuid"
	"github.com/kyleaupton/snaggle/backend/internal/quality"
)

type Plan struct {
	DownloaderID   string // how to download
	LibraryID      string // where to move/hardlink/copy the file to
	NameTemplateID string // how to name the file
}

// TorrentMetadata represents metadata about a torrent for policy evaluation
type CandidateContext struct {
	Candidate DownloadCandidate
	Quality   quality.ParsedQuality
}

type Policy struct {
	ID          uuid.UUID
	Name        string
	Description string
	Enabled     bool
	Priority    int

	Condition Rule
	Actions   []Action
}

type Operator string

const (
	OpEq       Operator = "=="
	OpNe       Operator = "!="
	OpGt       Operator = ">"
	OpGte      Operator = ">="
	OpLt       Operator = "<"
	OpLte      Operator = "<="
	OpContains Operator = "contains"
	OpIn       Operator = "in"
	OpNotIn    Operator = "not in"
	OpAnd      Operator = "and"
	OpOr       Operator = "or"
	OpNot      Operator = "not"
)

type Rule struct {
	ID       uuid.UUID
	Left     string
	Operator Operator
	Right    string
}

type ActionType string

const (
	ActionSetDownloader   ActionType = "set_downloader"
	ActionSetLibrary      ActionType = "set_library"
	ActionSetNameTemplate ActionType = "set_name_template"
	ActionStopProcessing  ActionType = "stop_processing"
)

type Action struct {
	ID    uuid.UUID
	Type  ActionType
	Value string
}

// EvaluationTrace represents the detailed trace of policy evaluation
type EvaluationTrace struct {
	Policies  []PolicyEvaluation `json:"policies"`
	FinalPlan Plan               `json:"finalPlan"`
}

// PolicyEvaluation represents the evaluation result for a single policy
type PolicyEvaluation struct {
	PolicyID          string       `json:"policyId"`
	PolicyName        string       `json:"policyName"`
	Priority          int          `json:"priority"`
	Matched           bool         `json:"matched"`
	RuleEvaluated     *RuleInfo    `json:"ruleEvaluated,omitempty"`
	ActionsApplied    []ActionInfo `json:"actionsApplied"`
	StoppedProcessing bool         `json:"stoppedProcessing"`
}

// RuleInfo represents information about a rule
type RuleInfo struct {
	LeftOperand  string `json:"leftOperand"`
	Operator     string `json:"operator"`
	RightOperand string `json:"rightOperand"`
}

// ActionInfo represents information about an action
type ActionInfo struct {
	Type  string `json:"type"`
	Value string `json:"value"`
	Order int32  `json:"order"`
}
