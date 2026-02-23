package input

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTxtSource(t *testing.T) {
	// テスト用テキストファイルを一時ディレクトリに作成
	tmpDir := t.TempDir()
	txtPath := filepath.Join(tmpDir, "test.txt")
	content := "吾輩は猫である。名前はまだ無い。"
	if err := os.WriteFile(txtPath, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	t.Run("Open returns txtSource", func(t *testing.T) {
		src, err := Open(txtPath)
		if err != nil {
			t.Fatalf("Open(%q) error: %v", txtPath, err)
		}
		defer src.Close()

		if src.Format() != FormatTxt {
			t.Errorf("Format() = %v, want %v", src.Format(), FormatTxt)
		}
	})

	t.Run("TextEntries returns single entry", func(t *testing.T) {
		src, err := Open(txtPath)
		if err != nil {
			t.Fatalf("Open(%q) error: %v", txtPath, err)
		}
		defer src.Close()

		entries := src.TextEntries()
		if len(entries) != 1 {
			t.Fatalf("TextEntries() len = %d, want 1", len(entries))
		}
		if entries[0].Name != "test.txt" {
			t.Errorf("TextEntries()[0].Name = %q, want %q", entries[0].Name, "test.txt")
		}
	})

	t.Run("ImageOnly is false", func(t *testing.T) {
		src, err := Open(txtPath)
		if err != nil {
			t.Fatalf("Open(%q) error: %v", txtPath, err)
		}
		defer src.Close()

		if src.ImageOnly() {
			t.Error("ImageOnly() = true, want false")
		}
	})

	t.Run("OpenText reads content", func(t *testing.T) {
		src, err := Open(txtPath)
		if err != nil {
			t.Fatalf("Open(%q) error: %v", txtPath, err)
		}
		defer src.Close()

		rc, entry, err := src.OpenText(0)
		if err != nil {
			t.Fatalf("OpenText(0) error: %v", err)
		}
		defer rc.Close()

		if entry.Name != "test.txt" {
			t.Errorf("entry.Name = %q, want %q", entry.Name, "test.txt")
		}

		data, err := io.ReadAll(rc)
		if err != nil {
			t.Fatalf("ReadAll error: %v", err)
		}
		if string(data) != content {
			t.Errorf("content = %q, want %q", string(data), content)
		}
	})

	t.Run("OpenText out of range", func(t *testing.T) {
		src, err := Open(txtPath)
		if err != nil {
			t.Fatalf("Open(%q) error: %v", txtPath, err)
		}
		defer src.Close()

		_, _, err = src.OpenText(1)
		if !errors.Is(err, ErrTextIndexOutOfRange) {
			t.Errorf("OpenText(1) error = %v, want %v", err, ErrTextIndexOutOfRange)
		}
	})

	t.Run("ArchiveTextParentPath returns parent dir", func(t *testing.T) {
		src, err := Open(txtPath)
		if err != nil {
			t.Fatalf("Open(%q) error: %v", txtPath, err)
		}
		defer src.Close()

		parentPath := src.ArchiveTextParentPath(0)
		if !strings.HasSuffix(parentPath, "/") {
			t.Errorf("ArchiveTextParentPath(0) = %q, should end with /", parentPath)
		}
		if !strings.Contains(parentPath, tmpDir) {
			t.Errorf("ArchiveTextParentPath(0) = %q, should contain %q", parentPath, tmpDir)
		}
	})

	t.Run("ImageEntries is empty", func(t *testing.T) {
		src, err := Open(txtPath)
		if err != nil {
			t.Fatalf("Open(%q) error: %v", txtPath, err)
		}
		defer src.Close()

		images := src.ImageEntries()
		if len(images) != 0 {
			t.Errorf("ImageEntries() len = %d, want 0", len(images))
		}
	})

	t.Run("nonexistent file returns error", func(t *testing.T) {
		_, err := Open(filepath.Join(tmpDir, "nonexistent.txt"))
		if err == nil {
			t.Error("Open(nonexistent) expected error, got nil")
		}
	})
}
