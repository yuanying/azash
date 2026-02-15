package normalize

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNormalizeTree(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()

	in := `<root b="2" a="1"><id>550e8400-e29b-41d4-a716-446655440000</id><date>2026-02-15T10:20:30Z</date><x>  a   b  </x></root>`
	if err := os.WriteFile(filepath.Join(src, "test.xml"), []byte(in), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := NormalizeTree(src, dst); err != nil {
		t.Fatalf("NormalizeTree error: %v", err)
	}
	outb, err := os.ReadFile(filepath.Join(dst, "test.xml"))
	if err != nil {
		t.Fatal(err)
	}
	out := string(outb)
	if !strings.Contains(out, "UUID_NORMALIZED") {
		t.Fatalf("uuid not normalized: %s", out)
	}
	if !strings.Contains(out, "DATE_NORMALIZED") {
		t.Fatalf("date not normalized: %s", out)
	}
}
