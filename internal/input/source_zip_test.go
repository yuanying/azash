package input

import (
	"archive/zip"
	"errors"
	"io"
	"os"
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
		defer func() { _ = src.Close() }()

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
		defer func() { _ = src.Close() }()

		rc, entry, err := src.OpenText(0)
		if err != nil {
			t.Fatalf("OpenText(0) error: %v", err)
		}
		defer func() { _ = rc.Close() }()

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
		defer func() { _ = src.Close() }()

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
		defer func() { _ = src.Close() }()

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
		defer func() { _ = src.Close() }()

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
		defer func() { _ = src.Close() }()

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
		defer func() { _ = src.Close() }()

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

func TestZipSource_MultipleTxtEntries(t *testing.T) {
	// 複数txtエントリを含むzipを作成
	tmpDir := t.TempDir()
	zipPath := filepath.Join(tmpDir, "multi.zip")
	createMultiTxtZip(t, zipPath)

	src, err := Open(zipPath)
	if err != nil {
		t.Fatalf("Open error: %v", err)
	}
	defer func() { _ = src.Close() }()

	t.Run("has 2 text entries", func(t *testing.T) {
		entries := src.TextEntries()
		if len(entries) != 2 {
			t.Fatalf("TextEntries() len = %d, want 2", len(entries))
		}
		if entries[0].Name != "chapter1.txt" {
			t.Errorf("TextEntries()[0].Name = %q, want %q", entries[0].Name, "chapter1.txt")
		}
		if entries[1].Name != "subdir/chapter2.txt" {
			t.Errorf("TextEntries()[1].Name = %q, want %q", entries[1].Name, "subdir/chapter2.txt")
		}
	})

	t.Run("OpenText(0) returns first txt", func(t *testing.T) {
		rc, entry, err := src.OpenText(0)
		if err != nil {
			t.Fatalf("OpenText(0) error: %v", err)
		}
		defer func() { _ = rc.Close() }()

		if entry.Name != "chapter1.txt" {
			t.Errorf("entry.Name = %q, want %q", entry.Name, "chapter1.txt")
		}
		data, err := io.ReadAll(rc)
		if err != nil {
			t.Fatalf("ReadAll error: %v", err)
		}
		if string(data) != "chapter 1 content" {
			t.Errorf("content = %q, want %q", string(data), "chapter 1 content")
		}
	})

	t.Run("OpenText(1) returns second txt", func(t *testing.T) {
		rc, entry, err := src.OpenText(1)
		if err != nil {
			t.Fatalf("OpenText(1) error: %v", err)
		}
		defer func() { _ = rc.Close() }()

		if entry.Name != "subdir/chapter2.txt" {
			t.Errorf("entry.Name = %q, want %q", entry.Name, "subdir/chapter2.txt")
		}
		data, err := io.ReadAll(rc)
		if err != nil {
			t.Fatalf("ReadAll error: %v", err)
		}
		if string(data) != "chapter 2 content" {
			t.Errorf("content = %q, want %q", string(data), "chapter 2 content")
		}
	})

	t.Run("ArchiveTextParentPath for subdir txt", func(t *testing.T) {
		p := src.ArchiveTextParentPath(1)
		if p != "subdir/" {
			t.Errorf("ArchiveTextParentPath(1) = %q, want %q", p, "subdir/")
		}
	})

	t.Run("ArchiveTextParentPath for root txt", func(t *testing.T) {
		p := src.ArchiveTextParentPath(0)
		if p != "" {
			t.Errorf("ArchiveTextParentPath(0) = %q, want empty", p)
		}
	})

	t.Run("webp image detected", func(t *testing.T) {
		images := src.ImageEntries()
		found := false
		for _, img := range images {
			if strings.HasSuffix(img.Name, ".webp") {
				found = true
				break
			}
		}
		if !found {
			t.Error("expected .webp in ImageEntries()")
		}
	})
}

func TestZipSource_Txtz(t *testing.T) {
	// txtz形式のテスト
	tmpDir := t.TempDir()
	txtzPath := filepath.Join(tmpDir, "test.txtz")
	createSingleTxtZip(t, txtzPath, "内容テスト")

	src, err := Open(txtzPath)
	if err != nil {
		t.Fatalf("Open error: %v", err)
	}
	defer func() { _ = src.Close() }()

	if src.Format() != FormatTxtz {
		t.Errorf("Format() = %v, want %v", src.Format(), FormatTxtz)
	}

	entries := src.TextEntries()
	if len(entries) != 1 {
		t.Fatalf("TextEntries() len = %d, want 1", len(entries))
	}
}

func createMultiTxtZip(t *testing.T, path string) {
	t.Helper()
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("failed to create zip: %v", err)
	}
	defer func() { _ = f.Close() }()

	w := zip.NewWriter(f)
	defer func() { _ = w.Close() }()

	entries := []struct {
		name    string
		content string
	}{
		{"chapter1.txt", "chapter 1 content"},
		{"subdir/chapter2.txt", "chapter 2 content"},
		{"cover.webp", "fake webp data"},
		{"image.png", "fake png data"},
	}

	for _, e := range entries {
		fw, err := w.Create(e.name)
		if err != nil {
			t.Fatalf("failed to create entry %s: %v", e.name, err)
		}
		if _, err := fw.Write([]byte(e.content)); err != nil {
			t.Fatalf("failed to write entry %s: %v", e.name, err)
		}
	}
}

func createSingleTxtZip(t *testing.T, path, content string) {
	t.Helper()
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("failed to create zip: %v", err)
	}
	defer func() { _ = f.Close() }()

	w := zip.NewWriter(f)
	defer func() { _ = w.Close() }()

	fw, err := w.Create("content.txt")
	if err != nil {
		t.Fatalf("failed to create entry: %v", err)
	}
	if _, err := fw.Write([]byte(content)); err != nil {
		t.Fatalf("failed to write entry: %v", err)
	}
}
