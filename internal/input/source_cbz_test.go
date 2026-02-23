package input

import (
	"archive/zip"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestCbzSource(t *testing.T) {
	// テスト用CBZファイルを作成 (画像のみのzip)
	tmpDir := t.TempDir()
	cbzPath := filepath.Join(tmpDir, "test.cbz")
	createTestCbz(t, cbzPath)

	t.Run("Format is cbz", func(t *testing.T) {
		src, err := Open(cbzPath)
		if err != nil {
			t.Fatalf("Open error: %v", err)
		}
		defer func() { _ = src.Close() }()

		if src.Format() != FormatCbz {
			t.Errorf("Format() = %v, want %v", src.Format(), FormatCbz)
		}
	})

	t.Run("ImageOnly is always true", func(t *testing.T) {
		src, err := Open(cbzPath)
		if err != nil {
			t.Fatalf("Open error: %v", err)
		}
		defer func() { _ = src.Close() }()

		if !src.ImageOnly() {
			t.Error("ImageOnly() = false, want true")
		}
	})

	t.Run("TextEntries is always empty", func(t *testing.T) {
		src, err := Open(cbzPath)
		if err != nil {
			t.Fatalf("Open error: %v", err)
		}
		defer func() { _ = src.Close() }()

		entries := src.TextEntries()
		if len(entries) != 0 {
			t.Errorf("TextEntries() len = %d, want 0", len(entries))
		}
	})

	t.Run("OpenText returns ErrNoTextEntry", func(t *testing.T) {
		src, err := Open(cbzPath)
		if err != nil {
			t.Fatalf("Open error: %v", err)
		}
		defer func() { _ = src.Close() }()

		_, _, err = src.OpenText(0)
		if !errors.Is(err, ErrNoTextEntry) {
			t.Errorf("OpenText(0) error = %v, want %v", err, ErrNoTextEntry)
		}
	})

	t.Run("ImageEntries returns images", func(t *testing.T) {
		src, err := Open(cbzPath)
		if err != nil {
			t.Fatalf("Open error: %v", err)
		}
		defer func() { _ = src.Close() }()

		images := src.ImageEntries()
		if len(images) != 2 {
			t.Errorf("ImageEntries() len = %d, want 2", len(images))
		}
	})
}

func createTestCbz(t *testing.T, path string) {
	t.Helper()
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("failed to create cbz: %v", err)
	}
	defer func() { _ = f.Close() }()

	w := zip.NewWriter(f)
	defer func() { _ = w.Close() }()

	for _, name := range []string{"page01.png", "page02.jpg"} {
		fw, err := w.Create(name)
		if err != nil {
			t.Fatalf("failed to create entry: %v", err)
		}
		if _, err := fw.Write([]byte("fake image data")); err != nil {
			t.Fatalf("failed to write entry: %v", err)
		}
	}
}
