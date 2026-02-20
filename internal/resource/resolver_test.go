package resource

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"testing/fstest"
)

func TestResolveFile_OrderCLIFirst(t *testing.T) {
	base := t.TempDir()
	cliDir := filepath.Join(base, "cli")
	cwdDir := filepath.Join(base, "cwd")
	bundleDir := filepath.Join(base, "bundle")
	mustWriteFile(t, filepath.Join(cliDir, "chuki_utf.txt"), "cli")
	mustWriteFile(t, filepath.Join(cwdDir, "chuki_utf.txt"), "cwd")
	mustWriteFile(t, filepath.Join(bundleDir, "chuki_utf.txt"), "bundle")

	r, err := New(Options{
		CLIBaseDir: cliDir,
		ExecDir:    cwdDir,
		BundleDir:  bundleDir,
	})
	if err != nil {
		t.Fatal(err)
	}

	res, err := r.ResolveFile("chuki_utf.txt")
	if err != nil {
		t.Fatal(err)
	}
	if res.Source != "cli" {
		t.Fatalf("want source cli, got %s", res.Source)
	}
}

func TestResolveFile_FallbackToBundle(t *testing.T) {
	base := t.TempDir()
	cwdDir := filepath.Join(base, "cwd")
	bundleDir := filepath.Join(base, "bundle")
	mustWriteFile(t, filepath.Join(bundleDir, "chuki_utf.txt"), "bundle")

	r, err := New(Options{
		ExecDir:   cwdDir,
		BundleDir: bundleDir,
	})
	if err != nil {
		t.Fatal(err)
	}

	res, err := r.ResolveFile("chuki_utf.txt")
	if err != nil {
		t.Fatal(err)
	}
	if res.Source != "bundle" {
		t.Fatalf("want source bundle, got %s", res.Source)
	}
}

func TestResolveFile_EmbedFallback(t *testing.T) {
	cwdDir := t.TempDir()
	r, err := New(Options{
		ExecDir: cwdDir,
		EmbedFS: fstest.MapFS{
			"chuki_utf.txt": &fstest.MapFile{Data: []byte("embed")},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	res, err := r.ResolveFile("chuki_utf.txt")
	if err != nil {
		t.Fatal(err)
	}
	if res.Source != "embed" {
		t.Fatalf("want source embed, got %s", res.Source)
	}
}

func TestResolveDir(t *testing.T) {
	base := t.TempDir()
	cliDir := filepath.Join(base, "cli")
	mustMkdirAll(t, filepath.Join(cliDir, "template"))

	r, err := New(Options{
		CLIBaseDir: cliDir,
		ExecDir:    filepath.Join(base, "cwd"),
	})
	if err != nil {
		t.Fatal(err)
	}

	res, err := r.ResolveDir("template")
	if err != nil {
		t.Fatal(err)
	}
	if !res.IsDir {
		t.Fatal("want directory resolution")
	}
}

func TestResolveFile_MissingError(t *testing.T) {
	cwdDir := t.TempDir()
	r, err := New(Options{
		ExecDir: cwdDir,
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = r.ResolveFile("missing.txt")
	if err == nil {
		t.Fatal("expected error")
	}
	var miss *MissingResourceError
	if !errors.As(err, &miss) {
		t.Fatalf("want MissingResourceError, got %T", err)
	}
	if miss.Resource != "missing.txt" {
		t.Fatalf("unexpected resource: %s", miss.Resource)
	}
	if len(miss.Searched) == 0 {
		t.Fatal("searched path list should not be empty")
	}
}

func TestResolveOptionalFile_NotFound(t *testing.T) {
	cwdDir := t.TempDir()
	r, err := New(Options{
		ExecDir: cwdDir,
	})
	if err != nil {
		t.Fatal(err)
	}

	_, ok, err := r.ResolveOptionalFile("replace.txt")
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("want not found")
	}
}

func TestOpenFile(t *testing.T) {
	cwdDir := t.TempDir()
	mustWriteFile(t, filepath.Join(cwdDir, "chuki_utf.txt"), "ok")

	r, err := New(Options{
		ExecDir: cwdDir,
	})
	if err != nil {
		t.Fatal(err)
	}

	rc, _, err := r.OpenFile("chuki_utf.txt")
	if err != nil {
		t.Fatal(err)
	}
	_ = rc.Close()
}

func TestDiagnostics(t *testing.T) {
	cwdDir := t.TempDir()
	mustWriteFile(t, filepath.Join(cwdDir, "chuki_utf.txt"), "ok")
	mustMkdirAll(t, filepath.Join(cwdDir, "template"))

	r, err := New(Options{
		ExecDir: cwdDir,
	})
	if err != nil {
		t.Fatal(err)
	}

	_, _ = r.ResolveFile("chuki_utf.txt")
	_, _ = r.ResolveDir("template")
	records := r.Diagnostics()
	if len(records) != 2 {
		t.Fatalf("want 2 records, got %d", len(records))
	}
}

func TestResolveFile_RejectPathEscape(t *testing.T) {
	cwdDir := t.TempDir()
	r, err := New(Options{ExecDir: cwdDir})
	if err != nil {
		t.Fatal(err)
	}

	tests := []string{
		"../secret.txt",
		"/etc/passwd",
		"..",
	}
	for _, tc := range tests {
		t.Run(tc, func(t *testing.T) {
			_, err := r.ResolveFile(tc)
			if err == nil {
				t.Fatalf("expected error for %q", tc)
			}
		})
	}
}

func TestResolveFileGlob(t *testing.T) {
	base := t.TempDir()
	cliDir := filepath.Join(base, "cli")
	cwdDir := filepath.Join(base, "cwd")
	mustWriteFile(t, filepath.Join(cliDir, "replace.txt"), "cli-main")
	mustWriteFile(t, filepath.Join(cliDir, "replace_extra.txt"), "cli-extra")
	mustWriteFile(t, filepath.Join(cwdDir, "replace.txt"), "cwd-main")

	r, err := New(Options{
		CLIBaseDir: cliDir,
		ExecDir:    cwdDir,
	})
	if err != nil {
		t.Fatal(err)
	}

	got, err := r.ResolveFileGlob("replace*.txt")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Fatalf("want 2 files from cli root, got %d", len(got))
	}
	names := []string{got[0].Resource, got[1].Resource}
	want := []string{"replace.txt", "replace_extra.txt"}
	if !reflect.DeepEqual(names, want) {
		t.Fatalf("unexpected resources: got %v, want %v", names, want)
	}
	for _, res := range got {
		if res.Source != "cli" {
			t.Fatalf("expected cli source, got %s", res.Source)
		}
	}
}

func TestReadDir(t *testing.T) {
	cwdDir := t.TempDir()
	mustWriteFile(t, filepath.Join(cwdDir, "web", "example.com", "extract.txt"), "ok")

	r, err := New(Options{ExecDir: cwdDir})
	if err != nil {
		t.Fatal(err)
	}

	entries, _, err := r.ReadDir("web")
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || entries[0].Name() != "example.com" || !entries[0].IsDir() {
		t.Fatalf("unexpected entries: %#v", entries)
	}
}

func mustWriteFile(t *testing.T, path string, content string) {
	t.Helper()
	mustMkdirAll(t, filepath.Dir(path))
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func mustMkdirAll(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatal(err)
	}
}
