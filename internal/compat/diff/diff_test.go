package diff

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCompareTrees(t *testing.T) {
	ja := t.TempDir()
	goa := t.TempDir()

	if err := os.WriteFile(filepath.Join(ja, "a.xml"), []byte(`<r><x a="1">hello</x></r>`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(goa, "a.xml"), []byte(`<r><x a="2">hello2</x></r>`), 0o644); err != nil {
		t.Fatal(err)
	}

	s, fds, cats, err := CompareTrees(ja, goa, Options{})
	if err != nil {
		t.Fatalf("CompareTrees error: %v", err)
	}
	if s.FilesCompared != 1 || s.FilesChanged != 1 {
		t.Fatalf("unexpected summary: %#v", s)
	}
	if len(fds) != 1 || !fds[0].Changed {
		t.Fatalf("unexpected file diffs: %#v", fds)
	}
	if len(cats) == 0 {
		t.Fatalf("expected non-empty categories")
	}
}
