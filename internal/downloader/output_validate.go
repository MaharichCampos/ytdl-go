package downloader

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kkdai/youtube/v2"
)

func validateOutputFile(path string, format *youtube.Format) error {
	info, err := os.Stat(path)
	if err != nil {
		return wrapCategory(CategoryFilesystem, fmt.Errorf("stat output: %w", err))
	}
	if info.Size() == 0 {
		return wrapCategory(CategoryUnsupported, fmt.Errorf("output file is empty"))
	}

	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".mp4", ".m4v", ".mov", ".m4a", ".m4s":
		return validateMP4(path)
	case ".webm", ".mkv":
		return validateEBML(path)
	case ".ts":
		return validateMPEGTS(path)
	case ".mp3":
		return validateMP3(path)
	default:
		if format != nil && strings.Contains(strings.ToLower(format.MimeType), "mp4") {
			return validateMP4(path)
		}
		return nil
	}
}

func readHeader(path string, size int) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	buf := make([]byte, size)
	n, err := file.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf[:n], nil
}

func validateMP4(path string) error {
	header, err := readHeader(path, 12)
	if err != nil {
		return wrapCategory(CategoryFilesystem, fmt.Errorf("read mp4 header: %w", err))
	}
	if len(header) < 8 || string(header[4:8]) != "ftyp" {
		return wrapCategory(CategoryUnsupported, fmt.Errorf("invalid mp4 header"))
	}
	body, err := readHeader(path, 1024*1024)
	if err != nil {
		return wrapCategory(CategoryFilesystem, fmt.Errorf("read mp4 body: %w", err))
	}
	if !bytes.Contains(body, []byte("moov")) && !bytes.Contains(body, []byte("moof")) {
		return wrapCategory(CategoryUnsupported, fmt.Errorf("missing moov/moof atom"))
	}
	return nil
}

func validateEBML(path string) error {
	header, err := readHeader(path, 4)
	if err != nil {
		return wrapCategory(CategoryFilesystem, fmt.Errorf("read ebml header: %w", err))
	}
	if len(header) < 4 || binary.BigEndian.Uint32(header) != 0x1A45DFA3 {
		return wrapCategory(CategoryUnsupported, fmt.Errorf("invalid webm header"))
	}
	return nil
}

func validateMPEGTS(path string) error {
	header, err := readHeader(path, 189)
	if err != nil {
		return wrapCategory(CategoryFilesystem, fmt.Errorf("read ts header: %w", err))
	}
	if len(header) < 1 || header[0] != 0x47 {
		return wrapCategory(CategoryUnsupported, fmt.Errorf("invalid transport stream header"))
	}
	if len(header) >= 189 && header[188] != 0x47 {
		return wrapCategory(CategoryUnsupported, fmt.Errorf("invalid transport stream sync"))
	}
	return nil
}

func validateMP3(path string) error {
	header, err := readHeader(path, 3)
	if err != nil {
		return wrapCategory(CategoryFilesystem, fmt.Errorf("read mp3 header: %w", err))
	}
	if len(header) < 3 {
		return wrapCategory(CategoryUnsupported, fmt.Errorf("invalid mp3 header"))
	}
	if string(header) == "ID3" || header[0] == 0xFF {
		return nil
	}
	return wrapCategory(CategoryUnsupported, fmt.Errorf("invalid mp3 header"))
}
