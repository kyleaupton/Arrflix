package quality

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestResolve(t *testing.T) {
	p := NewParser()
	res := p.Parse("21.Jump.Street.2012.2160p.UHD.BluRay.REMUX.HEVC.TrueHD.Atmos-GROUP")

	jsonData, _ := json.MarshalIndent(res, "", "  ")
	fmt.Println(string(jsonData))
}
