// Package mediainfo provides video file analysis using ffprobe or mediainfo.
// This package extracts detailed metadata from media files that is not available
// from release titles alone (e.g., actual video codec, audio tracks, HDR info).
package mediainfo

import (
	"github.com/kyleaupton/snaggle/backend/internal/model"
)

// Analyzer extracts technical metadata from video files.
type Analyzer struct {
	// ffprobePath is the path to the ffprobe binary
	ffprobePath string
}

// NewAnalyzer creates a new media analyzer.
// By default, it looks for ffprobe in the PATH.
func NewAnalyzer() *Analyzer {
	return &Analyzer{
		ffprobePath: "ffprobe",
	}
}

// WithFFprobePath sets a custom path to the ffprobe binary.
func (a *Analyzer) WithFFprobePath(path string) *Analyzer {
	a.ffprobePath = path
	return a
}

// Analyze extracts technical metadata from a video file.
// Returns nil if the file cannot be analyzed.
//
// TODO: Implement actual ffprobe integration. This is a stub that returns
// empty MediaInfoFields to allow the rest of the system to compile.
func (a *Analyzer) Analyze(filePath string) *model.MediaInfoFields {
	// TODO: Implement ffprobe execution and parsing
	// Example ffprobe command:
	//   ffprobe -v quiet -print_format json -show_streams -show_format <file>
	//
	// For now, return nil to indicate mediainfo is not yet available.
	// The EvaluationContext will gracefully handle nil MediaInfo.
	return nil
}

// AnalyzeWithDefaults is a convenience function that uses the default analyzer settings.
func AnalyzeWithDefaults(filePath string) *model.MediaInfoFields {
	return NewAnalyzer().Analyze(filePath)
}

// VideoCodec represents detected video codecs
type VideoCodec string

const (
	VideoCodecUnknown VideoCodec = "Unknown"
	VideoCodecH264    VideoCodec = "H.264"
	VideoCodecH265    VideoCodec = "H.265"
	VideoCodecAV1     VideoCodec = "AV1"
	VideoCodecVP9     VideoCodec = "VP9"
	VideoCodecMPEG2   VideoCodec = "MPEG-2"
)

// AudioCodec represents detected audio codecs
type AudioCodec string

const (
	AudioCodecUnknown AudioCodec = "Unknown"
	AudioCodecAAC     AudioCodec = "AAC"
	AudioCodecAC3     AudioCodec = "AC3"
	AudioCodecDTS     AudioCodec = "DTS"
	AudioCodecDTSHDMA AudioCodec = "DTS-HD MA"
	AudioCodecTrueHD  AudioCodec = "TrueHD"
	AudioCodecFLAC    AudioCodec = "FLAC"
	AudioCodecOpus    AudioCodec = "Opus"
)

// AudioChannels represents detected audio channel layouts
type AudioChannels string

const (
	AudioChannelsUnknown AudioChannels = "Unknown"
	AudioChannels20      AudioChannels = "2.0"
	AudioChannels51      AudioChannels = "5.1"
	AudioChannels71      AudioChannels = "7.1"
)

// HDRFormat represents detected HDR formats
type HDRFormat string

const (
	HDRNone        HDRFormat = "None"
	HDR10          HDRFormat = "HDR10"
	HDR10Plus      HDRFormat = "HDR10+"
	HDRDolbyVision HDRFormat = "Dolby Vision"
	HDRHLG         HDRFormat = "HLG"
)

// Container represents detected container formats
type Container string

const (
	ContainerUnknown Container = "Unknown"
	ContainerMKV     Container = "MKV"
	ContainerMP4     Container = "MP4"
	ContainerAVI     Container = "AVI"
	ContainerTS      Container = "TS"
)
