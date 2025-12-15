package quality

import (
	"regexp"
	"sort"
	"strings"
	"unicode"
)

type ParsedQuality struct {
	RawTitle   string   `json:"rawTitle"`
	Normalized string   `json:"normalized"`
	Tokens     []string `json:"tokens,omitempty"`

	Resolution Resolution `json:"resolution"`
	Source     Source     `json:"source"`
	Codec      VideoCodec `json:"codec"`
	Container  Container  `json:"container"`

	HDR      HDRFormat `json:"hdr"`
	BitDepth BitDepth  `json:"bitDepth"`

	Audio AudioInfo `json:"audio"`

	IsRemux    bool `json:"isRemux"`
	IsProper   bool `json:"isProper"`
	IsRepack   bool `json:"isRepack"`
	IsExtended bool `json:"isExtended"`
	IsUnrated  bool `json:"isUnrated"`
	IsDubbed   bool `json:"isDubbed"`
	IsSubbed   bool `json:"isSubbed"`

	Edition      string `json:"edition,omitempty"`
	ReleaseGroup string `json:"releaseGroup,omitempty"`

	Confidence Confidence `json:"confidence"`
	Warnings   []string   `json:"warnings,omitempty"`

	Tier QualityTier `json:"tier"`
}

type AudioInfo struct {
	Codec       AudioCodec `json:"codec"`
	Channels    Channels   `json:"channels"`
	Atmos       bool       `json:"atmos"`
	DTSX        bool       `json:"dtsx"`
	BitrateKbps int        `json:"bitrateKbps,omitempty"`
}

type Confidence struct {
	Overall    float32 `json:"overall"`
	Resolution float32 `json:"resolution"`
	Source     float32 `json:"source"`
	Codec      float32 `json:"codec"`
	HDR        float32 `json:"hdr"`
	Audio      float32 `json:"audio"`
}

type Resolution string

const (
	ResUnknown Resolution = "unknown"
	Res480     Resolution = "480p"
	Res576     Resolution = "576p"
	Res720     Resolution = "720p"
	Res1080    Resolution = "1080p"
	Res1440    Resolution = "1440p"
	Res2160    Resolution = "2160p"
	Res4320    Resolution = "4320p"
)

type Source string

const (
	SourceUnknown Source = "unknown"
	SourceCAM     Source = "cam"
	SourceTS      Source = "telesync"
	SourceTC      Source = "telecine"
	SourceSCR     Source = "screener"
	SourceDVD     Source = "dvd"
	SourceDVDRip  Source = "dvd-rip"
	SourceHDTV    Source = "hdtv"
	SourceWEBRip  Source = "webrip"
	SourceWEBDL   Source = "web-dl"
	SourceBluRay  Source = "bluray"
	SourceREMUX   Source = "remux"
)

type VideoCodec string

const (
	VCUnknown VideoCodec = "unknown"
	VCH264    VideoCodec = "h264"
	VCH265    VideoCodec = "h265"
	VCAV1     VideoCodec = "av1"
	VCVP9     VideoCodec = "vp9"
	VCMPEG2   VideoCodec = "mpeg2"
)

type Container string

const (
	ContUnknown Container = "unknown"
	ContMKV     Container = "mkv"
	ContMP4     Container = "mp4"
	ContAVI     Container = "avi"
	ContTS      Container = "ts"
)

type HDRFormat string

const (
	HDRUnknown     HDRFormat = "unknown"
	HDRNone        HDRFormat = "none"
	HDR10          HDRFormat = "hdr10"
	HDR10Plus      HDRFormat = "hdr10+"
	HDRDolbyVision HDRFormat = "dolby_vision"
	HDRHLG         HDRFormat = "hlg"
)

type BitDepth string

const (
	BitUnknown BitDepth = "unknown"
	Bit8       BitDepth = "8"
	Bit10      BitDepth = "10"
	Bit12      BitDepth = "12"
)

type AudioCodec string

const (
	ACUnknown AudioCodec = "unknown"
	ACAAC     AudioCodec = "aac"
	ACAC3     AudioCodec = "ac3"
	ACEAC3    AudioCodec = "eac3"
	ACDTS     AudioCodec = "dts"
	ACTrueHD  AudioCodec = "truehd"
	ACFLAC    AudioCodec = "flac"
	ACMP3     AudioCodec = "mp3"
)

type Channels string

const (
	ChUnknown Channels = "unknown"
	Ch20      Channels = "2.0"
	Ch51      Channels = "5.1"
	Ch71      Channels = "7.1"
)

type QualityTier string

const (
	TierUnknown  QualityTier = "unknown"
	TierLow      QualityTier = "low"       // cam/ts/tc/scr
	TierSD       QualityTier = "sd"        // dvd/480/576
	TierHD       QualityTier = "hd"        // 720
	TierFullHD   QualityTier = "full_hd"   // 1080
	TierUHD      QualityTier = "uhd"       // 2160/4320
	TierRemux    QualityTier = "remux"     // remux (non-UHD)
	TierUHDRemux QualityTier = "uhd_remux" // 2160 remux
)

// Parser is stateless; keep compiled regexes here.
type Parser struct {
	reRes        *regexp.Regexp
	reGroup      *regexp.Regexp
	reChan       *regexp.Regexp
	reBitDepth   *regexp.Regexp
	reBitrateK   *regexp.Regexp
	reCleanSplit *regexp.Regexp
}

func NewParser() *Parser {
	return &Parser{
		// 2160p / 1080p / 720p / etc
		reRes: regexp.MustCompile(`\b(480|576|720|1080|1440|2160|4320)p\b`),

		// Release group: last "-GROUP" segment or "[GROUP]" style (best-effort)
		reGroup: regexp.MustCompile(`(?:-|\[)([A-Za-z0-9]{2,})(?:\]?)\s*$`),

		// Channels: 2.0 / 5.1 / 7.1 / 5 1 / 7 1
		reChan: regexp.MustCompile(`\b(2(?:\.| )0|5(?:\.| )1|7(?:\.| )1)\b`),

		// Bit depth: 8bit/10bit/12bit or 8-bit etc
		reBitDepth: regexp.MustCompile(`\b(8|10|12)\s*[- ]?\s*bit\b`),

		// Bitrate in kbps e.g. 768kbps
		reBitrateK: regexp.MustCompile(`\b(\d{3,5})\s*kbps\b`),

		// split on non-alnum, keep dots/underscores/hyphens as separators
		reCleanSplit: regexp.MustCompile(`[^\p{L}\p{N}]+`),
	}
}

func (p *Parser) Parse(title string) ParsedQuality {
	q := ParsedQuality{
		RawTitle:   title,
		Normalized: normalizeTitle(title),
		Resolution: ResUnknown,
		Source:     SourceUnknown,
		Codec:      VCUnknown,
		Container:  ContUnknown,
		HDR:        HDRUnknown,
		BitDepth:   BitUnknown,
		Audio: AudioInfo{
			Codec:    ACUnknown,
			Channels: ChUnknown,
		},
		Tier: TierUnknown,
	}

	tokens := tokenize(p.reCleanSplit, q.Normalized)
	q.Tokens = tokens

	// --- Primary: Resolution ---
	if m := p.reRes.FindStringSubmatch(q.Normalized); len(m) == 2 {
		q.Resolution = mapResolution(m[1])
		q.Confidence.Resolution = 1.0
	}

	// --- Source / Traits ---
	applySourceAndTraits(&q, tokens)

	// --- Video codec ---
	applyVideoCodec(&q, tokens)

	// --- Container ---
	applyContainer(&q, tokens)

	// --- HDR + bit depth ---
	applyHDR(&q, tokens)
	applyBitDepth(p, &q)

	// --- Audio ---
	applyAudioCodec(&q, tokens)
	applyChannels(p, &q)
	applyAudioFlags(&q, tokens)
	applyAudioBitrate(p, &q)

	// --- Edition / Cut flags (optional but useful) ---
	applyEdition(&q, tokens)

	// --- Group (best effort) ---
	q.ReleaseGroup = extractGroup(p.reGroup, title)

	// --- Derived tier + validations ---
	q.Tier = deriveTier(&q)
	addWarnings(&q)

	// --- Confidence ---
	q.Confidence.Source = scoreUnknown(q.Source != SourceUnknown)
	q.Confidence.Codec = scoreUnknown(q.Codec != VCUnknown)
	q.Confidence.HDR = scoreHDR(q.HDR, q.BitDepth)
	q.Confidence.Audio = scoreAudio(q.Audio)

	q.Confidence.Overall = clamp01(0.30*q.Confidence.Resolution +
		0.25*q.Confidence.Source +
		0.20*q.Confidence.Codec +
		0.15*q.Confidence.Audio +
		0.10*q.Confidence.HDR)

	return q
}

// -------------------------
// Normalization / Tokenizing
// -------------------------

func normalizeTitle(s string) string {
	s = strings.TrimSpace(s)
	// Lowercase for matching; keep original in RawTitle.
	s = strings.ToLower(s)
	// Make common separators uniform-ish.
	s = strings.ReplaceAll(s, "_", " ")
	return s
}

func tokenize(splitter *regexp.Regexp, normalized string) []string {
	raw := splitter.Split(normalized, -1)
	out := make([]string, 0, len(raw))
	for _, t := range raw {
		t = strings.TrimSpace(t)
		if t == "" {
			continue
		}
		out = append(out, t)
	}
	return out
}

// -------------------------
// Resolution
// -------------------------

func mapResolution(n string) Resolution {
	switch n {
	case "480":
		return Res480
	case "576":
		return Res576
	case "720":
		return Res720
	case "1080":
		return Res1080
	case "1440":
		return Res1440
	case "2160":
		return Res2160
	case "4320":
		return Res4320
	default:
		return ResUnknown
	}
}

// -------------------------
// Source / Traits
// -------------------------

func applySourceAndTraits(q *ParsedQuality, tokens []string) {
	has := makeSet(tokens)

	// Traits
	q.IsRemux = has["remux"]
	q.IsProper = has["proper"]
	q.IsRepack = has["repack"]
	q.IsExtended = has["extended"] || has["ext"] || has["extendedcut"] || (has["extended"] && has["cut"])
	q.IsUnrated = has["unrated"]
	q.IsDubbed = has["dubbed"] || has["dub"]
	q.IsSubbed = has["subbed"] || has["subs"] || has["sub"]

	// Source scoring priority: Remux > BluRay > WEB-DL > WEBRip > HDTV > DVDRip/DVD > Screener > TS/TC/CAM
	// (You can tweak this list to match your preferences.)
	switch {
	case q.IsRemux:
		q.Source = SourceREMUX
	case has["bluray"] || has["blu-ray"] || has["bdrip"] || has["brrip"] || has["bdremux"]:
		q.Source = SourceBluRay
	case has["webdl"] || has["web-dl"] || (has["web"] && has["dl"]) || has["webdlrip"]:
		q.Source = SourceWEBDL
	case has["webrip"] || (has["web"] && has["rip"]):
		q.Source = SourceWEBRip
	case has["hdtv"]:
		q.Source = SourceHDTV
	case has["dvdrip"] || has["dvd-rip"]:
		q.Source = SourceDVDRip
	case has["dvd"]:
		q.Source = SourceDVD
	case has["scr"] || has["screener"]:
		q.Source = SourceSCR
	case has["ts"] || has["telesync"]:
		q.Source = SourceTS
	case has["tc"] || has["telecine"]:
		q.Source = SourceTC
	case has["cam"]:
		q.Source = SourceCAM
	default:
		// leave unknown
	}
}

// -------------------------
// Video Codec
// -------------------------

func applyVideoCodec(q *ParsedQuality, tokens []string) {
	has := makeSet(tokens)

	// Prefer explicit modern markers.
	switch {
	case has["av1"]:
		q.Codec = VCAV1
	case has["hevc"] || has["h265"] || has["x265"] || has["h.265"] || has["x.265"]:
		q.Codec = VCH265
	case has["h264"] || has["x264"] || has["h.264"] || has["x.264"] || has["avc"]:
		q.Codec = VCH264
	case has["vp9"]:
		q.Codec = VCVP9
	case has["mpeg2"] || has["mpeg-2"]:
		q.Codec = VCMPEG2
	default:
		// unknown
	}
}

// -------------------------
// Container
// -------------------------

func applyContainer(q *ParsedQuality, tokens []string) {
	has := makeSet(tokens)
	switch {
	case has["mkv"]:
		q.Container = ContMKV
	case has["mp4"]:
		q.Container = ContMP4
	case has["avi"]:
		q.Container = ContAVI
	case has["ts"]:
		// beware: "ts" can be telesync too. If source is TS, keep it.
		// container TS is common for raw streams; only set if we didn't already identify telesync source.
		if q.Source != SourceTS {
			q.Container = ContTS
		}
	}
}

// -------------------------
// HDR / BitDepth
// -------------------------

func applyHDR(q *ParsedQuality, tokens []string) {
	has := makeSet(tokens)

	// Dolby Vision markers
	if has["dv"] || has["dovi"] || has["dolbyvision"] || (has["dolby"] && has["vision"]) {
		q.HDR = HDRDolbyVision
		return
	}

	// HDR10+ markers
	if has["hdr10+"] || has["hdr10plus"] || (has["hdr10"] && has["plus"]) {
		q.HDR = HDR10Plus
		return
	}

	// HDR10 markers
	if has["hdr10"] {
		q.HDR = HDR10
		return
	}

	// HLG markers
	if has["hlg"] {
		q.HDR = HDRHLG
		return
	}

	// “HDR” alone is ambiguous; treat as unknown-ish rather than HDR10.
	if has["hdr"] {
		q.HDR = HDRUnknown
		return
	}

	// If we see explicit SDR marker, set none.
	if has["sdr"] {
		q.HDR = HDRNone
		return
	}

	// Leave unknown
}

func applyBitDepth(p *Parser, q *ParsedQuality) {
	if m := p.reBitDepth.FindStringSubmatch(q.Normalized); len(m) == 2 {
		switch m[1] {
		case "8":
			q.BitDepth = Bit8
		case "10":
			q.BitDepth = Bit10
		case "12":
			q.BitDepth = Bit12
		}
		return
	}

	// Heuristic: HEVC + HDR often implies 10-bit, but don’t over-assert.
	// Keep unknown unless you want a “best guess” mode.
}

// -------------------------
// Audio
// -------------------------

func applyAudioCodec(q *ParsedQuality, tokens []string) {
	has := makeSet(tokens)

	// Prefer TrueHD / DTS / EAC3 markers first (more specific).
	switch {
	case has["truehd"] || (has["true"] && has["hd"]):
		q.Audio.Codec = ACTrueHD
	case has["eac3"] || has["e-ac3"] || has["ddp"] || has["dd+"]:
		q.Audio.Codec = ACEAC3
	case has["ac3"] || has["dd"] || has["dolbydigital"]:
		q.Audio.Codec = ACAC3
	case has["dtsx"]:
		q.Audio.Codec = ACDTS
		q.Audio.DTSX = true
	case has["dts"] || has["dtshd"] || has["dts-hd"]:
		q.Audio.Codec = ACDTS
	case has["aac"]:
		q.Audio.Codec = ACAAC
	case has["flac"]:
		q.Audio.Codec = ACFLAC
	case has["mp3"]:
		q.Audio.Codec = ACMP3
	default:
		// unknown
	}
}

func applyChannels(p *Parser, q *ParsedQuality) {
	if m := p.reChan.FindStringSubmatch(q.Normalized); len(m) == 2 {
		ch := strings.ReplaceAll(m[1], " ", ".")
		switch ch {
		case "2.0":
			q.Audio.Channels = Ch20
		case "5.1":
			q.Audio.Channels = Ch51
		case "7.1":
			q.Audio.Channels = Ch71
		}
	}
}

func applyAudioFlags(q *ParsedQuality, tokens []string) {
	has := makeSet(tokens)
	if has["atmos"] {
		q.Audio.Atmos = true
	}
	if has["dtsx"] {
		q.Audio.DTSX = true
		// if codec unknown but DTSX present, assume DTS family
		if q.Audio.Codec == ACUnknown {
			q.Audio.Codec = ACDTS
		}
	}
}

func applyAudioBitrate(p *Parser, q *ParsedQuality) {
	if m := p.reBitrateK.FindStringSubmatch(q.Normalized); len(m) == 2 {
		q.Audio.BitrateKbps = atoiSafe(m[1])
	}
}

// -------------------------
// Edition / Cut
// -------------------------

func applyEdition(q *ParsedQuality, tokens []string) {
	has := makeSet(tokens)

	// very small starter set; expand later
	switch {
	case has["criterion"]:
		q.Edition = "criterion"
	case has["imax"]:
		q.Edition = "imax"
	case has["directorscut"] || (has["director"] && has["cut"]) || has["dc"]:
		q.Edition = "directors_cut"
	case has["extended"] && q.Edition == "":
		q.Edition = "extended"
	}
}

// -------------------------
// Tier + Warnings
// -------------------------

func deriveTier(q *ParsedQuality) QualityTier {
	// Remux tiers
	if q.Source == SourceREMUX || q.IsRemux {
		if q.Resolution == Res2160 || q.Resolution == Res4320 {
			return TierUHDRemux
		}
		return TierRemux
	}

	// Low-quality sources
	if q.Source == SourceCAM || q.Source == SourceTS || q.Source == SourceTC || q.Source == SourceSCR {
		return TierLow
	}

	// Resolution-based (with SD catch)
	switch q.Resolution {
	case Res480, Res576:
		return TierSD
	case Res720:
		return TierHD
	case Res1080, Res1440:
		return TierFullHD
	case Res2160, Res4320:
		return TierUHD
	default:
		// fallback by source if no res
		switch q.Source {
		case SourceDVD, SourceDVDRip:
			return TierSD
		case SourceHDTV, SourceWEBRip, SourceWEBDL, SourceBluRay:
			return TierUnknown
		default:
			return TierUnknown
		}
	}
}

func addWarnings(q *ParsedQuality) {
	var w []string

	// If DV/HDR but codec unknown, titles might be missing HEVC marker.
	if (q.HDR == HDRDolbyVision || q.HDR == HDR10 || q.HDR == HDR10Plus || q.HDR == HDRHLG) && q.Codec == VCUnknown {
		w = append(w, "HDR/DV detected but video codec not detected (common if title omits HEVC/x265).")
	}

	// DV usually rides on HEVC in the wild; warn if explicitly H264.
	if q.HDR == HDRDolbyVision && q.Codec == VCH264 {
		w = append(w, "Dolby Vision with H.264 is unusual; verify title/metadata.")
	}

	// REMUX without BluRay-ish source hints (often still OK, but warn).
	if q.IsRemux && q.Source == SourceREMUX && q.Resolution == ResUnknown {
		w = append(w, "REMUX detected but no resolution detected; verify title parsing.")
	}

	// TS ambiguity: "ts" can be telesync or container.
	if q.Source == SourceTS && q.Container == ContTS {
		w = append(w, "TS detected as both source and container; source likely telesync.")
	}

	q.Warnings = w
}

// -------------------------
// Helpers
// -------------------------

func makeSet(tokens []string) map[string]bool {
	m := make(map[string]bool, len(tokens))
	for _, t := range tokens {
		if t == "" {
			continue
		}
		m[t] = true
		// also store a “collapsed” token (no punctuation) for combos like "hdr10+"
		m[collapseToken(t)] = true
	}
	return m
}

func collapseToken(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || r == '+' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func extractGroup(re *regexp.Regexp, rawTitle string) string {
	// best effort: look for "-GROUP" or "[GROUP]" at the end
	m := re.FindStringSubmatch(strings.TrimSpace(rawTitle))
	if len(m) == 2 {
		return m[1]
	}
	return ""
}

func scoreUnknown(ok bool) float32 {
	if ok {
		return 1.0
	}
	return 0.0
}

func scoreHDR(h HDRFormat, bd BitDepth) float32 {
	// HDR is optional; unknown shouldn't heavily penalize overall score.
	switch h {
	case HDRDolbyVision, HDR10Plus, HDR10, HDRHLG:
		if bd != BitUnknown {
			return 1.0
		}
		return 0.8
	case HDRNone:
		return 0.7
	case HDRUnknown:
		return 0.6
	default:
		return 0.6
	}
}

func scoreAudio(a AudioInfo) float32 {
	s := float32(0.0)
	if a.Codec != ACUnknown {
		s += 0.6
	}
	if a.Channels != ChUnknown {
		s += 0.3
	}
	if a.Atmos || a.DTSX {
		s += 0.1
	}
	return clamp01(s)
}

func clamp01(x float32) float32 {
	if x < 0 {
		return 0
	}
	if x > 1 {
		return 1
	}
	return x
}

func atoiSafe(s string) int {
	n := 0
	for _, r := range s {
		if r < '0' || r > '9' {
			break
		}
		n = n*10 + int(r-'0')
	}
	return n
}

// Utility: give you stable deterministic display tokens if you want.
func UniqueSortedTokens(tokens []string) []string {
	m := map[string]struct{}{}
	for _, t := range tokens {
		if t == "" {
			continue
		}
		m[t] = struct{}{}
	}
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
