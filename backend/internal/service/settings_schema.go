package service

type SettingType string

const (
	SettingText SettingType = "text"
	SettingBool SettingType = "bool"
	SettingInt  SettingType = "int"
	SettingJSON SettingType = "json"
)

type SettingSpec struct {
	Key     string
	Type    SettingType
	Default any
}

// Registry enumerates all supported settings and their types/defaults.
// Extend this map as your application grows.
var Registry = map[string]SettingSpec{
	"site.title":            {Key: "site.title", Type: SettingText, Default: "Snaggle"},
	"auth.allow_signups":    {Key: "auth.allow_signups", Type: SettingBool, Default: false},
	"requests.max_per_user": {Key: "requests.max_per_user", Type: SettingInt, Default: int64(5)},
}
