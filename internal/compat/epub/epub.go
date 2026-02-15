package epub

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type EntryMeta struct {
	Name  string `json:"name"`
	Index int    `json:"index"`
	Size  int64  `json:"size"`
}

func ListEntries(epubPath string) ([]EntryMeta, error) {
	r, err := zip.OpenReader(epubPath)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	out := make([]EntryMeta, 0, len(r.File))
	for i, f := range r.File {
		out = append(out, EntryMeta{Name: f.Name, Index: i, Size: int64(f.UncompressedSize64)})
	}
	return out, nil
}

func Unpack(epubPath, outDir string) error {
	r, err := zip.OpenReader(epubPath)
	if err != nil {
		return err
	}
	defer r.Close()

	if err := os.RemoveAll(outDir); err != nil {
		return err
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}

	for _, f := range r.File {
		target := filepath.Join(outDir, f.Name)
		if !stringsHasPathPrefix(target, outDir) {
			return fmt.Errorf("illegal path traversal: %s", f.Name)
		}
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(target, 0o755); err != nil {
				return err
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}
		wf, err := os.Create(target)
		if err != nil {
			rc.Close()
			return err
		}
		if _, err := io.Copy(wf, rc); err != nil {
			wf.Close()
			rc.Close()
			return err
		}
		if err := wf.Close(); err != nil {
			rc.Close()
			return err
		}
		if err := rc.Close(); err != nil {
			return err
		}
	}
	return nil
}

func stringsHasPathPrefix(path, root string) bool {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return false
	}
	return rel != ".." && rel != "." && !startsWithDotDot(rel)
}

func startsWithDotDot(s string) bool {
	return len(s) >= 2 && s[0:2] == ".."
}
