package template

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

// Package template provides functionality for rendering name templates.
// This is a minimal implementation intended for filesystem-safe naming.

// Render renders a template string with the provided variables.
func Render(tmplStr string, data any) (string, error) {
	funcMap := template.FuncMap{
		"sanitize": sanitizeValue,
		"clean":    cleanQualityValue,
	}

	t, err := template.New("naming").Funcs(funcMap).Parse(tmplStr)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func sanitizeValue(s string) string {
	// Avoid path traversal and illegal/annoying path chars across common filesystems.
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

// CleanTitle sanitizes a title for use in filenames
// This is the same logic as sanitizeValue but exported for use in context building
func CleanTitle(s string) string {
	return sanitizeValue(s)
}

// cleanQualityValue returns an empty string if the value is "unknown", otherwise returns the value sanitized.
func cleanQualityValue(v any) string {
	s := strings.TrimSpace(strings.ToLower(fmt.Sprintf("%v", v)))
	if s == "unknown" || s == "" {
		return ""
	}
	return sanitizeValue(s)
}
