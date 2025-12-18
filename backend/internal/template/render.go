package template

import "strings"

// Package template provides functionality for rendering name templates.
// This is a minimal implementation intended for filesystem-safe naming.

// Render renders a template string with the provided variables.
func Render(template string, vars map[string]string) (string, error) {
	out := template
	for k, v := range vars {
		v = sanitizeValue(v)
		out = strings.ReplaceAll(out, "{"+k+"}", v)
		out = strings.ReplaceAll(out, "{{"+k+"}}", v)
	}
	return out, nil
}

func sanitizeValue(s string) string {
	// Avoid path traversal and illegal/annoying path chars across common filesystems.
	// Keep this conservative; callers can add more normalization later.
	r := strings.NewReplacer(
		"/", "-",
		"\\", "-",
		":", "-",
		"*", "",
		"?", "",
		"\"", "",
		"<", "",
		">", "",
		"|", "",
	)
	s = r.Replace(s)
	s = strings.TrimSpace(s)
	return s
}

