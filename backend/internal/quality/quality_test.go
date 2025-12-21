package quality

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/kyleaupton/snaggle/backend/internal/template"
)

func TestResolve(t *testing.T) {
	p := NewParser()
	res := p.Parse("21.Jump.Street.2012.2160p.UHD.BluRay.REMUX.HEVC.TrueHD.Atmos-GROUP")

	jsonData, _ := json.MarshalIndent(res, "", "  ")
	fmt.Println(string(jsonData))
}

func TestQualityInTemplate(t *testing.T) {
	p := NewParser()
	res := p.Parse("21.Jump.Street.2012.2160p.UHD.BluRay.REMUX.HEVC.TrueHD.Atmos-GROUP")

	ctx := NamingContext{
		Title:   "21 Jump Street",
		Year:    "2012",
		Quality: res,
	}

	tmpl := "{{.Title}} ({{.Year}}) [{{.Quality.Resolution}} {{.Quality.Codec}}]"
	got, err := template.Render(tmpl, ctx)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	want := "21 Jump Street (2012) [2160p h265]"
	if got != want {
		t.Errorf("Render() got = %v, want %v", got, want)
	}
}
