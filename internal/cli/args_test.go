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
