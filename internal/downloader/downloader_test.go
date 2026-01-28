package downloader

import (
	"strings"
	"testing"

	"github.com/kkdai/youtube/v2"
)

func TestParseVideoQuality(t *testing.T) {
	target, preferLowest, err := parseVideoQuality("720p")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if target != 720 || preferLowest {
		t.Fatalf("expected 720p target, got %d (preferLowest=%v)", target, preferLowest)
	}

	target, preferLowest, err = parseVideoQuality("worst")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if target != 0 || !preferLowest {
		t.Fatalf("expected worst to request lowest quality, got target=%d preferLowest=%v", target, preferLowest)
	}
}

func TestSelectFormatRespectsMaxHeight(t *testing.T) {
	video := &youtube.Video{
		Formats: youtube.FormatList{
			{ItagNo: 18, MimeType: "video/mp4", Width: 640, Height: 360, Bitrate: 800_000, AudioChannels: 2},
			{ItagNo: 22, MimeType: "video/mp4", Width: 1280, Height: 720, Bitrate: 2_000_000, AudioChannels: 2},
			{ItagNo: 247, MimeType: "video/webm", Width: 1280, Height: 720, Bitrate: 1_600_000, AudioChannels: 2},
		},
	}

	format, err := selectFormat(video, Options{Quality: "480p"})
	if err != nil {
		t.Fatalf("selectFormat returned error: %v", err)
	}
	if format.Height != 360 {
		t.Fatalf("expected 360p selection, got %dp (itag %d)", format.Height, format.ItagNo)
	}

	format, err = selectFormat(video, Options{Quality: "720p", Format: "webm"})
	if err != nil {
		t.Fatalf("selectFormat returned error: %v", err)
	}
	if format.ItagNo != 247 {
		t.Fatalf("expected webm format (itag 247), got itag %d", format.ItagNo)
	}
}

func TestSelectAudioBitrate(t *testing.T) {
	video := &youtube.Video{
		Formats: youtube.FormatList{
			{ItagNo: 140, MimeType: "audio/mp4", AudioChannels: 2, Bitrate: 128_000},
			{ItagNo: 251, MimeType: "audio/webm", AudioChannels: 2, Bitrate: 160_000},
		},
	}

	format, err := selectFormat(video, Options{AudioOnly: true, Quality: "128k"})
	if err != nil {
		t.Fatalf("selectFormat returned error: %v", err)
	}
	if format.ItagNo != 140 {
		t.Fatalf("expected itag 140 for 128k target, got %d", format.ItagNo)
	}
}

func TestSelectFormatUnsupportedIncludesHint(t *testing.T) {
	video := &youtube.Video{
		Formats: youtube.FormatList{
			{ItagNo: 18, MimeType: "video/mp4", Width: 640, Height: 360, Bitrate: 800_000, AudioChannels: 2},
		},
	}

	_, err := selectFormat(video, Options{Format: "webm"})
	if err == nil {
		t.Fatalf("expected error for unavailable format")
	}
	if CategoryOf(err) != CategoryUnsupported {
		t.Fatalf("expected unsupported category, got %s", CategoryOf(err))
	}
	if !strings.Contains(err.Error(), "--list-formats") {
		t.Fatalf("expected hint to use --list-formats, got %q", err.Error())
	}
}

func TestValidateInputURL(t *testing.T) {
	_, err := validateInputURL("ftp://example.com/video.mp4")
	if err == nil {
		t.Fatalf("expected error for unsupported scheme")
	}
	if CategoryOf(err) != CategoryInvalidURL {
		t.Fatalf("expected invalid_url category, got %s", CategoryOf(err))
	}
}
