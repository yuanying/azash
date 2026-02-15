package samples

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadCSV(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "samples.csv")
	data := "id,title,source_url,downloaded_at,input_file,expected_encoding,tags,notes\n" +
		"sample-001,Title,https://example.com,2026-02-14,a.zip,AUTO,\"ruby,chuki\",note\n"
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}

	s, err := LoadCSV(path)
	if err != nil {
		t.Fatalf("LoadCSV error: %v", err)
	}
	if len(s) != 1 {
		t.Fatalf("expected 1 sample, got %d", len(s))
	}
	if s[0].ID != "sample-001" {
		t.Fatalf("unexpected id: %s", s[0].ID)
	}
	if len(s[0].Tags) != 2 {
		t.Fatalf("unexpected tags: %#v", s[0].Tags)
	}
}

func TestLoadCSV_DuplicateID(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "samples.csv")
	data := "id,title,source_url,downloaded_at,input_file,expected_encoding,tags,notes\n" +
		"sample-001,Title,https://example.com,2026-02-14,a.zip,AUTO,ruby,note\n" +
		"sample-001,Title2,https://example.com,2026-02-14,b.zip,AUTO,ruby,note\n"
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadCSV(path); err == nil {
		t.Fatal("expected duplicate id error")
	}
}
