package io

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"bucket/internal/domain"
)

func LoadBucket(path string) (domain.Bucket, error) {
	bs, err := os.ReadFile(path)
	if err != nil {
		return domain.Bucket{}, fmt.Errorf("read %s: %w", path, err)
	}
	nl := detectNewline(bs)
	hadTrailing := hasTrailingNewline(bs)
	text := string(bs)
	lines := splitLines(text)

	return domain.Bucket{
		Path:                path,
		Name:                filepath.Base(path),
		Lines:               lines,
		SavedBaselineText:   text,
		SavedBaselineCount:  len(lines),
		Newline:             nl,
		BaselineHadTrailing: hadTrailing,
		Cursor:              0,
	}, nil
}

func detectNewline(bs []byte) string {
	if bytes.Contains(bs, []byte("\r\n")) {
		return "\r\n"
	}
	return "\n"
}

func hasTrailingNewline(bs []byte) bool {
	return bytes.HasSuffix(bs, []byte("\n"))
}

func splitLines(text string) []string {
	if text == "" {
		return []string{}
	}
	text = strings.ReplaceAll(text, "\r\n", "\n")
	hasTrailing := strings.HasSuffix(text, "\n")
	text = strings.TrimSuffix(text, "\n")
	if text == "" {
		if hasTrailing {
			return []string{""}
		}
		return []string{}
	}
	return strings.Split(text, "\n")
}

func Serialize(lines []string, nl string, hadTrailing bool) string {
	if nl == "" {
		nl = "\n"
	}
	text := strings.Join(lines, nl)
	if hadTrailing {
		text += nl
	}
	return text
}

func IsDirty(b domain.Bucket) bool {
	return Serialize(b.Lines, b.Newline, b.BaselineHadTrailing) != b.SavedBaselineText
}

func SaveDirty(b *domain.Bucket) error {
	next := Serialize(b.Lines, b.Newline, b.BaselineHadTrailing)
	if next == b.SavedBaselineText {
		return nil
	}
	if err := AtomicWriteFile(b.Path, []byte(next), 0o644); err != nil {
		return err
	}
	b.SavedBaselineText = next
	b.SavedBaselineCount = len(b.Lines)
	return nil
}

func AtomicWriteFile(path string, data []byte, mode os.FileMode) error {
	dir := filepath.Dir(path)
	f, err := os.CreateTemp(dir, ".bucket-*")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	tmp := f.Name()
	defer os.Remove(tmp)

	if _, err := f.Write(data); err != nil {
		f.Close()
		return fmt.Errorf("write temp file: %w", err)
	}
	if err := f.Sync(); err != nil {
		f.Close()
		return fmt.Errorf("sync temp file: %w", err)
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("close temp file: %w", err)
	}
	if err := os.Chmod(tmp, mode); err != nil {
		return fmt.Errorf("chmod temp file: %w", err)
	}
	if err := os.Rename(tmp, path); err != nil {
		return fmt.Errorf("rename temp file: %w", err)
	}
	return nil
}
