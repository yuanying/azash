package input

import (
	"errors"
	"io"
	"path/filepath"
	"strings"
	"testing"
)

func testdataPath(name string) string {
	return filepath.Join("..", "..", "testdata", "compat", "input", name)
}

func TestZipSource(t *testing.T) {
	t.Run("aozora_127_ruby.zip has 1 text entry", func(t *testing.T) {
		src, err := Open(testdataPath("aozora_127_ruby.zip"))
		if err != nil {
			t.Fatalf("Open error: %v", err)
		}
		defer src.Close()

		if src.Format() != FormatZip {
			t.Errorf("Format() = %v, want %v", src.Format(), FormatZip)
		}

		entries := src.TextEntries()
		if len(entries) != 1 {
			t.Fatalf("TextEntries() len = %d, want 1", len(entries))
		}
		if entries[0].Name != "rashomon.txt" {
			t.Errorf("TextEntries()[0].Name = %q, want %q", entries[0].Name, "rashomon.txt")
		}

		if src.ImageOnly() {
			t.Error("ImageOnly() = true, want false")
		}
	})

	t.Run("aozora_127_ruby.zip OpenText reads content", func(t *testing.T) {
		src, err := Open(testdataPath("aozora_127_ruby.zip"))
		if err != nil {
			t.Fatalf("Open error: %v", err)
		}
		defer src.Close()

		rc, entry, err := src.OpenText(0)
		if err != nil {
			t.Fatalf("OpenText(0) error: %v", err)
		}
		defer rc.Close()

		if entry.Name != "rashomon.txt" {
			t.Errorf("entry.Name = %q, want %q", entry.Name, "rashomon.txt")
		}

		data, err := io.ReadAll(rc)
		if err != nil {
			t.Fatalf("ReadAll error: %v", err)
		}
		if len(data) == 0 {
			t.Error("OpenText returned empty content")
		}
	})

	t.Run("aozora_127_ruby.zip ArchiveTextParentPath", func(t *testing.T) {
		src, err := Open(testdataPath("aozora_127_ruby.zip"))
		if err != nil {
			t.Fatalf("Open error: %v", err)
		}
		defer src.Close()

		path := src.ArchiveTextParentPath(0)
		// rashomon.txt はルート直下なので空文字列
		if path != "" {
			t.Errorf("ArchiveTextParentPath(0) = %q, want empty", path)
		}
	})

	t.Run("test_png.zip is imageOnly", func(t *testing.T) {
		src, err := Open(testdataPath("test_png.zip"))
		if err != nil {
			t.Fatalf("Open error: %v", err)
		}
		defer src.Close()

		if !src.ImageOnly() {
			t.Error("ImageOnly() = false, want true")
		}

		entries := src.TextEntries()
		if len(entries) != 0 {
			t.Errorf("TextEntries() len = %d, want 0", len(entries))
		}

		images := src.ImageEntries()
		if len(images) != 11 {
			t.Errorf("ImageEntries() len = %d, want 11", len(images))
		}
	})

	t.Run("test_png.zip OpenText returns ErrNoTextEntry", func(t *testing.T) {
		src, err := Open(testdataPath("test_png.zip"))
		if err != nil {
			t.Fatalf("Open error: %v", err)
		}
		defer src.Close()

		_, _, err = src.OpenText(0)
		if !errors.Is(err, ErrNoTextEntry) {
			t.Errorf("OpenText(0) error = %v, want %v", err, ErrNoTextEntry)
		}
	})

	t.Run("test_png.zip image entries have correct paths", func(t *testing.T) {
		src, err := Open(testdataPath("test_png.zip"))
		if err != nil {
			t.Fatalf("Open error: %v", err)
		}
		defer src.Close()

		images := src.ImageEntries()
		for _, img := range images {
			if !strings.HasPrefix(img.Name, "test_png/") {
				t.Errorf("ImageEntry.Name = %q, should start with 'test_png/'", img.Name)
			}
			if !strings.HasSuffix(img.Name, ".png") {
				t.Errorf("ImageEntry.Name = %q, should end with '.png'", img.Name)
			}
		}
	})

	t.Run("OpenText out of range", func(t *testing.T) {
		src, err := Open(testdataPath("aozora_127_ruby.zip"))
		if err != nil {
			t.Fatalf("Open error: %v", err)
		}
		defer src.Close()

		_, _, err = src.OpenText(1)
		if !errors.Is(err, ErrTextIndexOutOfRange) {
			t.Errorf("OpenText(1) error = %v, want %v", err, ErrTextIndexOutOfRange)
		}
	})

	t.Run("nonexistent zip returns error", func(t *testing.T) {
		_, err := Open(testdataPath("nonexistent.zip"))
		if err == nil {
			t.Error("Open(nonexistent.zip) expected error, got nil")
		}
	})
}
