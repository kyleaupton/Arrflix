package template

import (
	"testing"

	"github.com/kyleaupton/snaggle/backend/internal/quality"
)

func TestRender(t *testing.T) {
	q := quality.ParseQuality("21.Jump.Street.2012.2160p.UHD.BluRay.REMUX.HEVC.TrueHD.Atmos-GROUP")

	ctx := map[string]any{
		"Title":   "21 Jump Street",
		"Year":    "2012",
		"Quality": q,
	}

	tests := []struct {
		name     string
		template string
		want     string
	}{
		{
			name:     "Simple title and year",
			template: "{{.Title}} ({{.Year}})",
			want:     "21 Jump Street (2012)",
		},
		{
			name:     "With quality resolution",
			template: "{{.Title}} ({{.Year}}) [{{.Quality.Resolution}}]",
			want:     "21 Jump Street (2012) [2160p]",
		},
		// {
		// 	name:     "With quality resolution and codec",
		// 	template: "{{.Title}} ({{.Year}}) [{{.Quality.Resolution}} {{.Quality.Codec}}]",
		// 	want:     "21 Jump Street (2012) [2160p h265]",
		// },
		{
			name:     "With clean function for unknown",
			template: "{{.Title}} ({{.Year}}) [{{clean .Quality.Resolution}}]",
			want:     "21 Jump Street (2012) [2160p]",
		},
		// {
		// 	name:     "With unknown value using clean",
		// 	template: "{{.Title}} [{{clean .Quality.Edition}}]",
		// 	want:     "21 Jump Street []",
		// },
		{
			name:     "With sanitize function",
			template: "{{sanitize .Title}}",
			want:     "21 Jump Street",
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
