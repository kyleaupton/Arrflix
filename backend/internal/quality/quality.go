package quality

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Resolution string

const (
	ResUnknown Resolution = "Unknown"
	ResSD      Resolution = "SD"
	Res480p    Resolution = "480p"
	Res576p    Resolution = "576p"
	Res720p    Resolution = "720p"
	Res1080p   Resolution = "1080p"
	Res1440p   Resolution = "1440p"
	Res2160p   Resolution = "2160p"
	Res4320p   Resolution = "4320p"
)

type Source string

const (
	SourceUnknown Source = "Unknown"
	SourceSDTV    Source = "SDTV"
	SourceCAM     Source = "CAM"
	SourceTS      Source = "Telesync"
	SourceTC      Source = "Telecine"
	SourceSCR     Source = "Screener"
	SourceDVD     Source = "DVD"
	SourceDVDRip  Source = "DVD-Rip"
	SourceHDTV    Source = "HDTV"
	SourceWEBRip  Source = "WEBRip"
	SourceWEBDL   Source = "WEB-DL"
	SourceBluRay  Source = "BluRay"
	SourceREMUX   Source = "REMUX"
	SourceRAWHD   Source = "Raw-HD"
)

// Quality represents the Sonarr Quality ID
type Quality int

const (
	Unknown          Quality = 0
	SDTV             Quality = 1
	DVD              Quality = 2
	WEBDL1080p       Quality = 3
	HDTV720p         Quality = 4
	WEBDL720p        Quality = 5
	Bluray720p       Quality = 6
	Bluray1080p      Quality = 7
	WEBDL480p        Quality = 8
	HDTV1080p        Quality = 9
	RAWHD            Quality = 10
	WEBRip480p       Quality = 12
	Bluray480p       Quality = 13
	WEBRip720p       Quality = 14
	WEBRip1080p      Quality = 15
	HDTV2160p        Quality = 16
	WEBRip2160p      Quality = 17
	WEBDL2160p       Quality = 18
	Bluray2160p      Quality = 19
	Bluray1080pRemux Quality = 20
	Bluray2160pRemux Quality = 21
	Bluray576p       Quality = 22
)

func (q Quality) String() string {
	switch q {
	case Unknown:
		return "Unknown"
	case SDTV:
		return "SDTV"
	case DVD:
		return "DVD"
	case WEBDL1080p:
		return "WEB-DL-1080p"
	case HDTV720p:
		return "HDTV-720p"
	case WEBDL720p:
		return "WEB-DL-720p"
	case Bluray720p:
		return "Bluray-720p"
	case Bluray1080p:
		return "Bluray-1080p"
	case WEBDL480p:
		return "WEB-DL-480p"
	case HDTV1080p:
		return "HDTV-1080p"
	case RAWHD:
		return "Raw-HD"
	case WEBRip480p:
		return "WEBRip-480p"
	case Bluray480p:
		return "Bluray-480p"
	case WEBRip720p:
		return "WEBRip-720p"
	case WEBRip1080p:
		return "WEBRip-1080p"
	case HDTV2160p:
		return "HDTV-2160p"
	case WEBRip2160p:
		return "WEBRip-2160p"
	case WEBDL2160p:
		return "WEB-DL-2160p"
	case Bluray2160p:
		return "Bluray-2160p"
	case Bluray1080pRemux:
		return "Bluray-1080p Remux"
	case Bluray2160pRemux:
		return "Bluray-2160p Remux"
	case Bluray576p:
		return "Bluray-576p"
	default:
		return "Unknown"
	}
}

func (q Quality) Source() string {
	switch q {
	case SDTV:
		return string(SourceSDTV)
	case DVD:
		return string(SourceDVD)
	case WEBDL1080p, WEBDL720p, WEBDL480p, WEBDL2160p:
		return string(SourceWEBDL)
	case HDTV720p, HDTV1080p, HDTV2160p:
		return string(SourceHDTV)
	case Bluray720p, Bluray1080p, Bluray480p, Bluray2160p, Bluray576p, Bluray1080pRemux, Bluray2160pRemux:
		return string(SourceBluRay)
	case WEBRip480p, WEBRip720p, WEBRip1080p, WEBRip2160p:
		return string(SourceWEBRip)
	case RAWHD:
		return string(SourceRAWHD)
	default:
		return string(SourceUnknown)
	}
}

func (q Quality) Resolution() string {
	switch q {
	case WEBDL1080p, HDTV1080p, Bluray1080p, WEBRip1080p, Bluray1080pRemux:
		return string(Res1080p)
	case HDTV720p, WEBDL720p, Bluray720p, WEBRip720p:
		return string(Res720p)
	case WEBDL480p, Bluray480p, WEBRip480p:
		return string(Res480p)
	case HDTV2160p, WEBRip2160p, WEBDL2160p, Bluray2160p, Bluray2160pRemux:
		return string(Res2160p)
	case Bluray576p:
		return string(Res576p)
	case SDTV, DVD:
		return string(ResSD)
	default:
		return string(ResUnknown)
	}
}

func (q Quality) IsRemux() bool {
	return q == Bluray1080pRemux || q == Bluray2160pRemux
}

type Revision struct {
	Version  int
	Real     int
	IsRepack bool
}

type QualityModel struct {
	Quality  Quality
	Revision Revision
}

func (qm QualityModel) String() string {
	res := qm.Quality.String()
	if qm.Revision.Version > 1 {
		res += fmt.Sprintf(" v%d", qm.Revision.Version)
	}
	if qm.Revision.IsRepack {
		res += " [REPACK]"
	}
	return res
}

func (qm QualityModel) Source() string {
	return qm.Quality.Source()
}

func (qm QualityModel) Resolution() string {
	return qm.Quality.Resolution()
}

func (qm QualityModel) IsRemux() bool {
	return qm.Quality.IsRemux()
}

// Full returns the Radarr/Sonarr-style quality tag (e.g., "HDTV-720p", "WEBDL-1080p")
// without revision information. This matches what name templates expect: {{.Quality.Full}}
func (qm QualityModel) Full() string {
	return qm.Quality.String()
}

// Version returns the revision version number
func (qm QualityModel) Version() int {
	return qm.Revision.Version
}

// FieldInfo represents metadata about an available quality field
type FieldInfo struct {
	Name        string                         // Field name (e.g., "Full", "Resolution")
	Type        string                         // Field type: "string", "bool", "int"
	Description string                         // Human-readable description
	Accessor    func(QualityModel) interface{} // Function to get the field value
}

// QualityFields is the registry of all available quality fields
var QualityFields = []FieldInfo{
	{
		Name:        "Full",
		Type:        "string",
		Description: "Full quality tag (e.g., HDTV-720p, WEBDL-1080p)",
		Accessor:    func(qm QualityModel) interface{} { return qm.Full() },
	},
	{
		Name:        "Resolution",
		Type:        "string",
		Description: "Resolution value (e.g., 720p, 1080p, 2160p)",
		Accessor:    func(qm QualityModel) interface{} { return qm.Resolution() },
	},
	{
		Name:        "Source",
		Type:        "string",
		Description: "Source type (e.g., HDTV, WEB-DL, BluRay)",
		Accessor:    func(qm QualityModel) interface{} { return qm.Source() },
	},
	{
		Name:        "IsRemux",
		Type:        "bool",
		Description: "Whether the quality is a remux",
		Accessor:    func(qm QualityModel) interface{} { return qm.IsRemux() },
	},
	{
		Name:        "IsRepack",
		Type:        "bool",
		Description: "Whether the release is a repack",
		Accessor:    func(qm QualityModel) interface{} { return qm.Revision.IsRepack },
	},
	{
		Name:        "Version",
		Type:        "int",
		Description: "Revision version number",
		Accessor:    func(qm QualityModel) interface{} { return qm.Version() },
	},
}

// GetField retrieves a field value by name from a QualityModel
func GetField(name string, qm QualityModel) (interface{}, error) {
	for _, field := range QualityFields {
		if field.Name == name {
			return field.Accessor(qm), nil
		}
	}
	return nil, fmt.Errorf("unknown quality field: %s", name)
}

// ListFields returns all available quality fields
func ListFields() []FieldInfo {
	return QualityFields
}

var (
	ResolutionRegex = regexp.MustCompile(`(?i)\b(?:(?P<R360p>360p)|(?P<R480p>480p|480i|640x480|848x480)|(?P<R540p>540p)|(?P<R576p>576p)|(?P<R720p>720p|1280x720|960p)|(?P<R1080p>1080p|1920x1080|1440p|FHD|1080i|4kto1080p)|(?P<R2160p>2160p|3840x2160|4k[-_. ](?:UHD|HEVC|BD|H265)|(?:UHD|HEVC|BD|H265)[-_. ]4k))\b`)

	// Simplified sources for Go
	BlurayRegex = regexp.MustCompile(`(?i)\b(BluRay|Blu-Ray|HD-?DVD|BDMux|BD)\b`)
	WebDlRegex  = regexp.MustCompile(`(?i)\b(WEB[-_. ]DL(?:mux)?|WEBDL|AmazonHD|AmazonSD|iTunesHD|MaxdomeHD|NetflixU?HD|WebHD|HBOMaxHD|DisneyHD|[. ]WEB[. ](?:[xh][ .]?26[45]|AVC|HEVC|DDP?5[. ]1)|[. ]WEB$|(?:720|1080|2160)p[-. ]WEB[-. ]|[-. ]WEB[-. ](?:720|1080|2160)p|\b\s\/\sWEB\s\/\s\b|(?:AMZN|NF|DP)[. -]WEB[. -])`)
	WebRipRegex = regexp.MustCompile(`(?i)\b(WebRip|Web-Rip|WEBMux)\b`)
	HdtvRegex   = regexp.MustCompile(`(?i)\b(HDTV)\b`)
	DvdRegex    = regexp.MustCompile(`(?i)\b(DVD|DVDRip|NTSC|PAL|xvidvd)\b`)

	ProperRegex  = regexp.MustCompile(`(?i)\bproper\b`)
	RepackRegex  = regexp.MustCompile(`(?i)\b(repack\d?|rerip\d?)\b`)
	VersionRegex = regexp.MustCompile(`(?i)\d[-._ ]?v(?P<version>\d)[-._ ]|\[v(?P<version>\d)\]|repack(?P<version>\d)|rerip(?P<version>\d)|(?:480|576|720|1080|2160)p[._ ]v(?P<version>\d)`)

	RemuxRegex = regexp.MustCompile(`(?i)(?:[_. ]|\d{4}p-|\bHybrid-)(?P<remux>(?:(BD|UHD)[-_. ]?)?Remux)\b|(?P<remux>(?:(BD|UHD)[-_. ]?)?Remux[_. ]\d{4}p)`)
)

func ParseResolution(name string) int {
	match := ResolutionRegex.FindStringSubmatch(name)
	if match == nil {
		return 0
	}

	for i, groupName := range ResolutionRegex.SubexpNames() {
		if i != 0 && groupName != "" && match[i] != "" {
			switch groupName {
			case "R360p":
				return 360
			case "R480p":
				return 480
			case "R540p":
				return 540
			case "R576p":
				return 576
			case "R720p":
				return 720
			case "R1080p":
				return 1080
			case "R2160p":
				return 2160
			}
		}
	}
	return 0
}

func ParseQuality(name string) QualityModel {
	normalizedName := strings.ReplaceAll(name, "_", " ")
	normalizedName = strings.TrimSpace(normalizedName)

	result := QualityModel{Quality: Unknown}

	// Parse Revision
	if vMatch := VersionRegex.FindStringSubmatch(normalizedName); vMatch != nil {
		for i, groupName := range VersionRegex.SubexpNames() {
			if groupName == "version" && vMatch[i] != "" {
				v, _ := strconv.Atoi(vMatch[i])
				result.Revision.Version = v
			}
		}
	}

	if ProperRegex.MatchString(normalizedName) {
		if result.Revision.Version < 2 {
			result.Revision.Version = 2
		} else {
			result.Revision.Version++
		}
	}

	if RepackRegex.MatchString(normalizedName) {
		result.Revision.IsRepack = true
		if result.Revision.Version < 2 {
			result.Revision.Version = 2
		} else {
			result.Revision.Version++
		}
	}

	resolution := ParseResolution(normalizedName)
	remuxMatch := RemuxRegex.MatchString(normalizedName)

	// Source Detection
	if BlurayRegex.MatchString(normalizedName) {
		switch resolution {
		case 2160:
			if remuxMatch {
				result.Quality = Bluray2160pRemux
			} else {
				result.Quality = Bluray2160p
			}
		case 1080:
			if remuxMatch {
				result.Quality = Bluray1080pRemux
			} else {
				result.Quality = Bluray1080p
			}
		case 720:
			result.Quality = Bluray720p
		case 576:
			result.Quality = Bluray576p
		case 480, 360, 540:
			result.Quality = Bluray480p
		default:
			if remuxMatch {
				result.Quality = Bluray1080pRemux
			} else {
				result.Quality = Bluray720p
			}
		}
		return result
	}

	if WebDlRegex.MatchString(normalizedName) && !WebRipRegex.MatchString(normalizedName) {
		switch resolution {
		case 2160:
			result.Quality = WEBDL2160p
		case 1080:
			result.Quality = WEBDL1080p
		case 720:
			result.Quality = WEBDL720p
		default:
			result.Quality = WEBDL480p
		}
		return result
	}

	if WebRipRegex.MatchString(normalizedName) {
		switch resolution {
		case 2160:
			result.Quality = WEBRip2160p
		case 1080:
			result.Quality = WEBRip1080p
		case 720:
			result.Quality = WEBRip720p
		default:
			result.Quality = WEBRip480p
		}
		return result
	}

	if HdtvRegex.MatchString(normalizedName) {
		switch resolution {
		case 2160:
			result.Quality = HDTV2160p
		case 1080:
			result.Quality = HDTV1080p
		case 720:
			result.Quality = HDTV720p
		default:
			result.Quality = SDTV
		}
		return result
	}

	if DvdRegex.MatchString(normalizedName) {
		result.Quality = DVD
		return result
	}

	// Fallback to resolution if source unknown
	if resolution != 0 {
		switch resolution {
		case 2160:
			result.Quality = HDTV2160p
		case 1080:
			result.Quality = HDTV1080p
		case 720:
			result.Quality = HDTV720p
		default:
			result.Quality = SDTV
		}
	}

	return result
}
