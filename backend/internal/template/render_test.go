package template

import (
	"testing"

	"github.com/kyleaupton/snaggle/backend/internal/quality"
)

func TestRender(t *testing.T) {
	q := quality.ParseQuality("21.Jump.Street.2012.2160p.UHD.BluRay.REMUX.HEVC.TrueHD.Atmos-GROUP")

	// Create a namespaced context structure matching EvaluationContext.ToTemplateData()
	media := map[string]any{
		"Title":      "21 Jump Street",
		"CleanTitle": "21 Jump Street",
		"Year":       2012,
	}

	ctx := map[string]any{
		"Media":   media,
		"Quality": q,
	}

	tests := []struct {
		name     string
		template string
		want     string
	}{
		{
			name:     "Simple title and year",
			template: "{{.Media.Title}} ({{.Media.Year}})",
			want:     "21 Jump Street (2012)",
		},
		{
			name:     "With quality resolution",
			template: "{{.Media.Title}} ({{.Media.Year}}) [{{.Quality.Resolution}}]",
			want:     "21 Jump Street (2012) [2160p]",
		},
		{
			name:     "With clean title",
			template: "{{.Media.CleanTitle}} ({{.Media.Year}})",
			want:     "21 Jump Street (2012)",
		},
		{
			name:     "With clean function for unknown",
			template: "{{.Media.Title}} ({{.Media.Year}}) [{{clean .Quality.Resolution}}]",
			want:     "21 Jump Street (2012) [2160p]",
		},
		{
			name:     "With sanitize function",
			template: "{{sanitize .Media.Title}}",
			want:     "21 Jump Street",
		},
		{
			name:     "Full quality string",
			template: "{{.Media.CleanTitle}} ({{.Media.Year}}) [{{.Quality.Full}}]",
			want:     "21 Jump Street (2012) [Bluray-2160p Remux]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Render(tt.template, ctx)
			if err != nil {
				t.Fatalf("Render() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("Render() got = %v, want %v", got, tt.want)
			}
		})
	}
}
