//go:build integration

package downloader

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestIntegrationPublicDownload(t *testing.T) {
	url := os.Getenv("YTDL_TEST_URL")
	if url == "" {
		t.Skip("set YTDL_TEST_URL to run integration download test")
	}

	opts := Options{
		OutputTemplate: filepath.Join(t.TempDir(), "{title}.{ext}"),
		Timeout:        2 * time.Minute,
		Quiet:          true,
	}
	if err := Process(context.Background(), url, opts); err != nil {
		t.Fatalf("download failed: %v", err)
	}
}

func TestIntegrationProgressMultipleBars(t *testing.T) {
	output := captureStderr(t, func() {
		printer := newPrinter(Options{})
		bar1 := newProgressWriter(10, printer, "[1/2] one")
		bar2 := newProgressWriter(10, printer, "[2/2] two")

		_, _ = bar1.Write([]byte("12345"))
		_, _ = bar2.Write([]byte("12345"))
		printer.Log("log line")
		_, _ = bar1.Write([]byte("12345"))
		_, _ = bar2.Write([]byte("12345"))
		bar1.Finish()
		bar2.Finish()
	})

	if !strings.Contains(output, "[1/2] one") || !strings.Contains(output, "[2/2] two") {
		t.Fatalf("expected both bars in output, got %q", output)
	}
	if !strings.Contains(output, "log line") {
		t.Fatalf("expected log line in output, got %q", output)
	}
}
