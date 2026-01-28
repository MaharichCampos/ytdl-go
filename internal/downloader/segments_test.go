package downloader

import "testing"

func TestParseHLSManifestDetectsEncryption(t *testing.T) {
	manifestText := `#EXTM3U
#EXT-X-KEY:METHOD=AES-128,URI="https://example.com/key"
#EXTINF:10,
segment1.ts
`
	manifest, err := ParseHLSManifest([]byte(manifestText))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !manifest.Encrypted {
		t.Fatalf("expected encrypted manifest")
	}
	if manifest.KeyMethod != "AES-128" {
		t.Fatalf("expected AES-128, got %q", manifest.KeyMethod)
	}
	if manifest.KeyURI == "" {
		t.Fatalf("expected key URI")
	}
}

func TestDetectDASHDrm(t *testing.T) {
	xml := `<MPD><ContentProtection schemeIdUri="urn:uuid:edef8ba9-79d6-4ace-a3c8-27dcd51d21ed"/></MPD>`
	if ok, _ := DetectDASHDrm([]byte(xml)); !ok {
		t.Fatalf("expected DRM detection for Widevine UUID")
	}
}
