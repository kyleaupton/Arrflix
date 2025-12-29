package quality

import (
	"testing"
)

func TestQualityMethods(t *testing.T) {
	tests := []struct {
		name       string
		quality    Quality
		source     string
		resolution string
		isRemux    bool
	}{
		{"Bluray1080p", Bluray1080p, "BluRay", "1080p", false},
		{"Bluray1080pRemux", Bluray1080pRemux, "BluRay", "1080p", true},
		{"WEBDL2160p", WEBDL2160p, "WEB-DL", "2160p", false},
		{"HDTV720p", HDTV720p, "HDTV", "720p", false},
		{"SDTV", SDTV, "SDTV", "SD", false},
		{"DVD", DVD, "DVD", "SD", false},
		{"Unknown", Unknown, "Unknown", "Unknown", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.quality.Source(); got != tt.source {
				t.Errorf("Quality.Source() = %v, want %v", got, tt.source)
			}
			if got := tt.quality.Resolution(); got != tt.resolution {
				t.Errorf("Quality.Resolution() = %v, want %v", got, tt.resolution)
			}
			if got := tt.quality.IsRemux(); got != tt.isRemux {
				t.Errorf("Quality.IsRemux() = %v, want %v", got, tt.isRemux)
			}
		})
	}
}

func TestQualityModelMethods(t *testing.T) {
	qm := QualityModel{
		Quality: Bluray1080pRemux,
		Revision: Revision{
			Version:  2,
			IsRepack: true,
		},
	}

	if qm.Source() != "BluRay" {
		t.Errorf("QualityModel.Source() = %v, want %v", qm.Source(), "BluRay")
	}
	if qm.Resolution() != "1080p" {
		t.Errorf("QualityModel.Resolution() = %v, want %v", qm.Resolution(), "1080p")
	}
	if !qm.IsRemux() {
		t.Errorf("QualityModel.IsRemux() = %v, want %v", qm.IsRemux(), true)
	}
	if qm.String() != "Bluray-1080p Remux v2 [REPACK]" {
		t.Errorf("QualityModel.String() = %v, want %v", qm.String(), "Bluray-1080p Remux v2 [REPACK]")
	}
	if qm.Full() != "Bluray-1080p Remux" {
		t.Errorf("QualityModel.Full() = %v, want %v", qm.Full(), "Bluray-1080p Remux")
	}
	if qm.Version() != 2 {
		t.Errorf("QualityModel.Version() = %v, want %v", qm.Version(), 2)
	}
}

func TestFull(t *testing.T) {
	tests := []struct {
		name     string
		quality  Quality
		expected string
	}{
		{"HDTV720p", HDTV720p, "HDTV-720p"},
		{"WEBDL1080p", WEBDL1080p, "WEB-DL-1080p"},
		{"Bluray2160pRemux", Bluray2160pRemux, "Bluray-2160p Remux"},
		{"SDTV", SDTV, "SDTV"},
		{"Unknown", Unknown, "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qm := QualityModel{Quality: tt.quality}
			if got := qm.Full(); got != tt.expected {
				t.Errorf("QualityModel.Full() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestListFields(t *testing.T) {
	fields := ListFields()
	if len(fields) == 0 {
		t.Error("ListFields() returned empty slice")
	}

	// Check that expected fields are present
	fieldNames := make(map[string]bool)
	for _, field := range fields {
		fieldNames[field.Name] = true
	}

	expectedFields := []string{"Full", "Resolution", "Source", "IsRemux", "IsRepack", "Version"}
	for _, expected := range expectedFields {
		if !fieldNames[expected] {
			t.Errorf("ListFields() missing expected field: %s", expected)
		}
	}
}

func TestGetField(t *testing.T) {
	qm := QualityModel{
		Quality: HDTV720p,
		Revision: Revision{
			Version:  2,
			IsRepack: true,
		},
	}

	tests := []struct {
		name     string
		expected interface{}
	}{
		{"Full", "HDTV-720p"},
		{"Resolution", "720p"},
		{"Source", "HDTV"},
		{"IsRemux", false},
		{"IsRepack", true},
		{"Version", 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetField(tt.name, qm)
			if err != nil {
				t.Errorf("GetField(%s) error = %v", tt.name, err)
				return
			}
			if got != tt.expected {
				t.Errorf("GetField(%s) = %v, want %v", tt.name, got, tt.expected)
			}
		})
	}

	// Test unknown field
	_, err := GetField("UnknownField", qm)
	if err == nil {
		t.Error("GetField() with unknown field should return error")
	}
}
