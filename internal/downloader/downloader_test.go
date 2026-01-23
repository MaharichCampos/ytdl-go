package downloader

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/kkdai/youtube/v2"
)

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{name: "valid https", input: "https://www.youtube.com/watch?v=dQw4w9WgXcQ", wantErr: false},
		{name: "valid http", input: "http://example.com/video.mp4", wantErr: false},
		{name: "missing scheme", input: "www.youtube.com/watch?v=dQw4w9WgXcQ", wantErr: true},
		{name: "missing host", input: "https:///video", wantErr: true},
		{name: "unsupported scheme", input: "ftp://example.com/video", wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateURL(test.input)
			if test.wantErr && err == nil {
				t.Fatalf("expected error for %q", test.input)
			}
			if !test.wantErr && err != nil {
				t.Fatalf("unexpected error for %q: %v", test.input, err)
			}
		})
	}
}

func TestWriteFormats(t *testing.T) {
	video := &youtube.Video{
		Formats: []youtube.Format{
			{
				ItagNo:        22,
				Quality:       "hd720",
				QualityLabel:  "720p",
				MimeType:      "video/mp4",
				Bitrate:       2000000,
				AudioChannels: 2,
				Width:         1280,
				Height:        720,
				ContentLength: 1234,
			},
		},
	}

	var buf bytes.Buffer
	if err := writeFormats(&buf, video); err != nil {
		t.Fatalf("writeFormats returned error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "itag") {
		t.Fatalf("expected header in output, got: %q", output)
	}
	if !strings.Contains(output, "720p") {
		t.Fatalf("expected format row in output, got: %q", output)
	}
}

func TestWrapAccessError(t *testing.T) {
	tests := []struct {
		name       string
		input      error
		wantNil    bool
		wantPrefix string
	}{
		{
			name:    "nil error returns nil",
			input:   nil,
			wantNil: true,
		},
		{
			name:       "private video error gets wrapped",
			input:      errors.New("This video is private"),
			wantPrefix: "restricted access:",
		},
		{
			name:       "members only error gets wrapped",
			input:      errors.New("Members only content"),
			wantPrefix: "restricted access:",
		},
		{
			name:       "age-restricted error gets wrapped",
			input:      errors.New("This video is age-restricted"),
			wantPrefix: "restricted access:",
		},
		{
			name:       "non-restricted error is unchanged",
			input:      errors.New("network timeout"),
			wantPrefix: "",
		},
		{
			name:       "generic error is unchanged",
			input:      errors.New("something went wrong"),
			wantPrefix: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := wrapAccessError(test.input)
			if test.wantNil {
				if result != nil {
					t.Fatalf("expected nil, got %v", result)
				}
				return
			}
			if result == nil {
				t.Fatalf("expected error, got nil")
			}
			if test.wantPrefix != "" {
				if !strings.HasPrefix(result.Error(), test.wantPrefix) {
					t.Fatalf("expected error to start with %q, got %q", test.wantPrefix, result.Error())
				}
			} else {
				// Non-restricted errors should be unchanged
				if result.Error() != test.input.Error() {
					t.Fatalf("expected error %q to be unchanged, got %q", test.input.Error(), result.Error())
				}
			}
		})
	}
}
