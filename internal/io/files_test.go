package io

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSerializePreservesTrailingNewline(t *testing.T) {
	got := Serialize([]string{"a", "b"}, "\n", true)
	if got != "a\nb\n" {
		t.Fatalf("got %q", got)
	}
}

func TestAtomicWriteFile(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "x.txt")
	if err := AtomicWriteFile(p, []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}
	bs, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	if string(bs) != "hello" {
		t.Fatalf("got %q", string(bs))
	}
}
