package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgtype"
	dbgen "github.com/kyleaupton/snaggle/backend/internal/db/sqlc"
	"github.com/kyleaupton/snaggle/backend/internal/logger"
	"github.com/kyleaupton/snaggle/backend/internal/model"
	"github.com/kyleaupton/snaggle/backend/internal/policy"
	"github.com/kyleaupton/snaggle/backend/internal/quality"
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

// GetFieldDefinitions returns all available field definitions for policy rules
func (s *PoliciesService) GetFieldDefinitions(ctx context.Context) ([]model.FieldDefinition, error) {
	fields := []model.FieldDefinition{}

	// Candidate fields
	fields = append(fields, []model.FieldDefinition{
		{
			Path:      "candidate.size",
			Label:     "Size",
			Type:      model.FieldTypeNumber,
			ValueType: "int64",
			Operators: []string{"==", "!=", ">", ">=", "<", "<="},
		},
		{
			Path:      "candidate.title",
			Label:     "Title",
			Type:      model.FieldTypeText,
			ValueType: "string",
			Operators: []string{"==", "!=", "contains", "in", "not in"},
		},
		{
			Path:          "candidate.indexer",
			Label:         "Indexer",
			Type:          model.FieldTypeDynamic,
			ValueType:     "string",
			DynamicSource: "/api/v1/indexers/configured",
			Operators:     []string{"==", "!=", "in", "not in"},
		},
		{
			Path:      "candidate.indexer_id",
			Label:     "Indexer ID",
			Type:      model.FieldTypeNumber,
			ValueType: "int64",
			Operators: []string{"==", "!=", ">", ">=", "<", "<="},
		},
		{
			Path:          "candidate.categories",
			Label:         "Categories",
			Type:          model.FieldTypeDynamic,
			ValueType:     "[]string",
			DynamicSource: "", // TODO: Add categories endpoint if needed
			Operators:     []string{"contains", "in", "not in"},
		},
		{
			Path:      "candidate.protocol",
			Label:     "Protocol",
			Type:      model.FieldTypeEnum,
			ValueType: "string",
			EnumValues: []model.EnumValue{
				{Value: "torrent", Label: "Torrent"},
				{Value: "usenet", Label: "Usenet"},
			},
			Operators: []string{"==", "!=", "in", "not in"},
		},
		{
			Path:      "candidate.torrent_seeders",
			Label:     "Torrent Seeders",
			Type:      model.FieldTypeNumber,
			ValueType: "int64",
			Operators: []string{"==", "!=", ">", ">=", "<", "<="},
		},
		{
			Path:      "candidate.torrent_peers",
			Label:     "Torrent Peers",
			Type:      model.FieldTypeNumber,
			ValueType: "int64",
			Operators: []string{"==", "!=", ">", ">=", "<", "<="},
		},
	}...)

	// Quality fields - Resolution
	fields = append(fields, model.FieldDefinition{
		Path:      "quality.resolution",
		Label:     "Resolution",
		Type:      model.FieldTypeEnum,
		ValueType: "string",
		EnumValues: []model.EnumValue{
			{Value: string(quality.ResUnknown), Label: "Unknown"},
			{Value: string(quality.Res480), Label: "480p"},
			{Value: string(quality.Res576), Label: "576p"},
			{Value: string(quality.Res720), Label: "720p"},
			{Value: string(quality.Res1080), Label: "1080p"},
			{Value: string(quality.Res1440), Label: "1440p"},
			{Value: string(quality.Res2160), Label: "2160p"},
			{Value: string(quality.Res4320), Label: "4320p"},
		},
		Operators: []string{"==", "!=", "in", "not in"},
	})

	// Quality fields - Source
	fields = append(fields, model.FieldDefinition{
		Path:      "quality.source",
		Label:     "Source",
		Type:      model.FieldTypeEnum,
		ValueType: "string",
		EnumValues: []model.EnumValue{
			{Value: string(quality.SourceUnknown), Label: "Unknown"},
			{Value: string(quality.SourceCAM), Label: "CAM"},
			{Value: string(quality.SourceTS), Label: "Telesync"},
			{Value: string(quality.SourceTC), Label: "Telecine"},
			{Value: string(quality.SourceSCR), Label: "Screener"},
			{Value: string(quality.SourceDVD), Label: "DVD"},
			{Value: string(quality.SourceDVDRip), Label: "DVD-Rip"},
			{Value: string(quality.SourceHDTV), Label: "HDTV"},
			{Value: string(quality.SourceWEBRip), Label: "WEBRip"},
			{Value: string(quality.SourceWEBDL), Label: "WEB-DL"},
			{Value: string(quality.SourceBluRay), Label: "BluRay"},
			{Value: string(quality.SourceREMUX), Label: "REMUX"},
		},
		Operators: []string{"==", "!=", "in", "not in"},
	})

	// Quality fields - Codec
	fields = append(fields, model.FieldDefinition{
		Path:      "quality.codec",
		Label:     "Video Codec",
		Type:      model.FieldTypeEnum,
		ValueType: "string",
		EnumValues: []model.EnumValue{
			{Value: string(quality.VCUnknown), Label: "Unknown"},
			{Value: string(quality.VCH264), Label: "H.264"},
			{Value: string(quality.VCH265), Label: "H.265"},
			{Value: string(quality.VCAV1), Label: "AV1"},
			{Value: string(quality.VCVP9), Label: "VP9"},
			{Value: string(quality.VCMPEG2), Label: "MPEG-2"},
		},
		Operators: []string{"==", "!=", "in", "not in"},
	})

	// Quality fields - Container
	fields = append(fields, model.FieldDefinition{
		Path:      "quality.container",
		Label:     "Container",
		Type:      model.FieldTypeEnum,
		ValueType: "string",
		EnumValues: []model.EnumValue{
			{Value: string(quality.ContUnknown), Label: "Unknown"},
			{Value: string(quality.ContMKV), Label: "MKV"},
			{Value: string(quality.ContMP4), Label: "MP4"},
			{Value: string(quality.ContAVI), Label: "AVI"},
			{Value: string(quality.ContTS), Label: "TS"},
		},
		Operators: []string{"==", "!=", "in", "not in"},
	})

	// Quality fields - HDR
	fields = append(fields, model.FieldDefinition{
		Path:      "quality.hdr",
		Label:     "HDR Format",
		Type:      model.FieldTypeEnum,
		ValueType: "string",
		EnumValues: []model.EnumValue{
			{Value: string(quality.HDRUnknown), Label: "Unknown"},
			{Value: string(quality.HDRNone), Label: "None"},
			{Value: string(quality.HDR10), Label: "HDR10"},
			{Value: string(quality.HDR10Plus), Label: "HDR10+"},
			{Value: string(quality.HDRDolbyVision), Label: "Dolby Vision"},
			{Value: string(quality.HDRHLG), Label: "HLG"},
		},
		Operators: []string{"==", "!=", "in", "not in"},
	})

	// Quality fields - Bit Depth
	fields = append(fields, model.FieldDefinition{
		Path:      "quality.bit_depth",
		Label:     "Bit Depth",
		Type:      model.FieldTypeEnum,
		ValueType: "string",
		EnumValues: []model.EnumValue{
			{Value: string(quality.BitUnknown), Label: "Unknown"},
			{Value: string(quality.Bit8), Label: "8-bit"},
			{Value: string(quality.Bit10), Label: "10-bit"},
			{Value: string(quality.Bit12), Label: "12-bit"},
		},
		Operators: []string{"==", "!=", "in", "not in"},
	})

	// Quality fields - Audio Codec
	fields = append(fields, model.FieldDefinition{
		Path:      "quality.audio.codec",
		Label:     "Audio Codec",
		Type:      model.FieldTypeEnum,
		ValueType: "string",
		EnumValues: []model.EnumValue{
			{Value: string(quality.ACUnknown), Label: "Unknown"},
			{Value: string(quality.ACAAC), Label: "AAC"},
			{Value: string(quality.ACAC3), Label: "AC3"},
			{Value: string(quality.ACEAC3), Label: "E-AC3"},
			{Value: string(quality.ACDTS), Label: "DTS"},
			{Value: string(quality.ACTrueHD), Label: "TrueHD"},
			{Value: string(quality.ACFLAC), Label: "FLAC"},
			{Value: string(quality.ACMP3), Label: "MP3"},
		},
		Operators: []string{"==", "!=", "in", "not in"},
	})

	// Quality fields - Audio Channels
	fields = append(fields, model.FieldDefinition{
		Path:      "quality.audio.channels",
		Label:     "Audio Channels",
		Type:      model.FieldTypeEnum,
		ValueType: "string",
		EnumValues: []model.EnumValue{
			{Value: string(quality.ChUnknown), Label: "Unknown"},
			{Value: string(quality.Ch20), Label: "2.0"},
			{Value: string(quality.Ch51), Label: "5.1"},
			{Value: string(quality.Ch71), Label: "7.1"},
		},
		Operators: []string{"==", "!=", "in", "not in"},
	})

	// Quality fields - Tier
	fields = append(fields, model.FieldDefinition{
		Path:      "quality.tier",
		Label:     "Quality Tier",
		Type:      model.FieldTypeEnum,
		ValueType: "string",
		EnumValues: []model.EnumValue{
			{Value: string(quality.TierUnknown), Label: "Unknown"},
			{Value: string(quality.TierLow), Label: "Low"},
			{Value: string(quality.TierSD), Label: "SD"},
			{Value: string(quality.TierHD), Label: "HD"},
			{Value: string(quality.TierFullHD), Label: "Full HD"},
			{Value: string(quality.TierUHD), Label: "UHD"},
			{Value: string(quality.TierRemux), Label: "Remux"},
			{Value: string(quality.TierUHDRemux), Label: "UHD Remux"},
		},
		Operators: []string{"==", "!=", "in", "not in"},
	})

	// Quality fields - Boolean fields
	fields = append(fields, []model.FieldDefinition{
		{
			Path:      "quality.is_remux",
			Label:     "Is Remux",
			Type:      model.FieldTypeBoolean,
			ValueType: "bool",
			EnumValues: []model.EnumValue{
				{Value: "true", Label: "True"},
				{Value: "false", Label: "False"},
			},
			Operators: []string{"==", "!="},
		},
		{
			Path:      "quality.is_proper",
			Label:     "Is Proper",
			Type:      model.FieldTypeBoolean,
			ValueType: "bool",
			EnumValues: []model.EnumValue{
				{Value: "true", Label: "True"},
				{Value: "false", Label: "False"},
			},
			Operators: []string{"==", "!="},
		},
		{
			Path:      "quality.is_repack",
			Label:     "Is Repack",
			Type:      model.FieldTypeBoolean,
			ValueType: "bool",
			EnumValues: []model.EnumValue{
				{Value: "true", Label: "True"},
				{Value: "false", Label: "False"},
			},
			Operators: []string{"==", "!="},
		},
		{
			Path:      "quality.is_extended",
			Label:     "Is Extended",
			Type:      model.FieldTypeBoolean,
			ValueType: "bool",
			EnumValues: []model.EnumValue{
				{Value: "true", Label: "True"},
				{Value: "false", Label: "False"},
			},
			Operators: []string{"==", "!="},
		},
	}...)

	// Quality fields - Release Group (text)
	fields = append(fields, model.FieldDefinition{
		Path:      "quality.release_group",
		Label:     "Release Group",
		Type:      model.FieldTypeText,
		ValueType: "string",
		Operators: []string{"==", "!=", "contains", "in", "not in"},
	})

	return fields, nil
}
