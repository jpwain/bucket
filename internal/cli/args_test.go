package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseArgsRejectsSameFileByInode(t *testing.T) {
	dir := t.TempDir()
	a := filepath.Join(dir, "a.txt")
	link := filepath.Join(dir, "a-link.txt")
	if err := os.WriteFile(a, []byte("x\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(a, link); err != nil {
		t.Fatal(err)
	}
	if _, err := ParseArgs([]string{a, link}); err == nil {
		t.Fatal("expected same-file rejection")
	}
}

func TestParseArgsVersionFlag(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{"--version", []string{"--version"}},
		{"-v", []string{"-v"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseArgs(tt.args)
			if err != nil {
				t.Fatalf("ParseArgs() error = %v", err)
			}
			if !got.VersionRequested {
				t.Error("expected VersionRequested to be true")
			}
		})
	}
}
