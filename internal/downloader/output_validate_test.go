package downloader

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateOutputMP4(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "file.mp4")
	data := []byte{0x00, 0x00, 0x00, 0x18, 'f', 't', 'y', 'p', 'i', 's', 'o', 'm', 'm', 'o', 'o', 'v'}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("write mp4: %v", err)
	}
	if err := validateOutputFile(path, nil); err != nil {
		t.Fatalf("expected mp4 validation to pass: %v", err)
	}
}

func TestValidateOutputTS(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "file.ts")
	if err := os.WriteFile(path, []byte{0x47, 0x00, 0x00, 0x00}, 0o644); err != nil {
		t.Fatalf("write ts: %v", err)
	}
	if err := validateOutputFile(path, nil); err != nil {
		t.Fatalf("expected ts validation to pass: %v", err)
	}
}

func TestValidateOutputInvalid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "file.mp4")
	if err := os.WriteFile(path, []byte("nope"), 0o644); err != nil {
		t.Fatalf("write invalid: %v", err)
	}
	if err := validateOutputFile(path, nil); err == nil {
		t.Fatalf("expected validation failure")
	}
}
