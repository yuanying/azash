package resource

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type sourceKind string

const (
	sourceCLI    sourceKind = "cli"
	sourceCWD    sourceKind = "cwd"
	sourceBundle sourceKind = "bundle"
	sourceEmbed  sourceKind = "embed"
)

type candidate struct {
	kind sourceKind
	root string
}

// Resolution は採用したリソース解決結果を表す。
type Resolution struct {
	Resource string
	Source   string
	Path     string
	IsDir    bool
}

// MissingResourceError はリソース欠落時の統一エラー。
type MissingResourceError struct {
	Resource string
	Searched []string
}

func (e *MissingResourceError) Error() string {
	return fmt.Sprintf("resource not found: %s (searched: %v)", e.Resource, e.Searched)
}

type Options struct {
	CLIBaseDir string
	ExecDir    string
	BundleDir  string
	EmbedFS    fs.FS
}

// Resolver は Java 互換リソース探索を拡張した探索順で解決する。
// 探索順: CLI 指定 > 実行ディレクトリ > 同梱パス > embed fallback。
type Resolver struct {
	candidates []candidate
	embedFS    fs.FS
	records    []Resolution
}

func New(opts Options) (*Resolver, error) {
	execDir := opts.ExecDir
	if execDir == "" {
		wd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("resolve working directory: %w", err)
		}
		execDir = wd
	}
	execDir = filepath.Clean(execDir)

	bundleDir := opts.BundleDir
	if bundleDir == "" {
		exePath, err := os.Executable()
		if err == nil {
			bundleDir = filepath.Dir(exePath)
		}
	}
	if bundleDir != "" {
		bundleDir = filepath.Clean(bundleDir)
	}

	var cands []candidate
	if opts.CLIBaseDir != "" {
		cands = append(cands, candidate{kind: sourceCLI, root: filepath.Clean(opts.CLIBaseDir)})
	}
	cands = append(cands, candidate{kind: sourceCWD, root: execDir})
	if bundleDir != "" && bundleDir != execDir {
		cands = append(cands, candidate{kind: sourceBundle, root: bundleDir})
	}

	return &Resolver{
		candidates: cands,
		embedFS:    opts.EmbedFS,
		records:    make([]Resolution, 0, 8),
	}, nil
}

func (r *Resolver) OpenFile(resource string) (io.ReadCloser, Resolution, error) {
	res, err := r.ResolveFile(resource)
	if err != nil {
		return nil, Resolution{}, err
	}
	rc, err := r.openResolvedFile(res)
	if err != nil {
		return nil, Resolution{}, err
	}
	return rc, res, nil
}

func (r *Resolver) OpenOptionalFile(resource string) (io.ReadCloser, *Resolution, error) {
	res, ok, err := r.ResolveOptionalFile(resource)
	if err != nil {
		return nil, nil, err
	}
	if !ok {
		return nil, nil, nil
	}
	rc, err := r.openResolvedFile(res)
	if err != nil {
		return nil, nil, err
	}
	return rc, &res, nil
}

func (r *Resolver) ResolveFile(resource string) (Resolution, error) {
	res, ok, err := r.resolve(resource, false)
	if err != nil {
		return Resolution{}, err
	}
	if !ok {
		return Resolution{}, r.missingErr(resource)
	}
	return res, nil
}

func (r *Resolver) ResolveFileGlob(pattern string) ([]Resolution, error) {
	clean, err := sanitizeResourcePath(pattern)
	if err != nil {
		return nil, err
	}

	for _, c := range r.candidates {
		fullPattern := filepath.Join(c.root, filepath.FromSlash(clean))
		matches, err := filepath.Glob(fullPattern)
		if err != nil {
			return nil, fmt.Errorf("glob %s: %w", fullPattern, err)
		}
		if len(matches) == 0 {
			continue
		}
		resolved := make([]Resolution, 0, len(matches))
		for _, m := range matches {
			info, err := os.Stat(m)
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			if err != nil {
				return nil, fmt.Errorf("stat %s: %w", m, err)
			}
			if info.IsDir() {
				continue
			}
			rel, err := filepath.Rel(c.root, m)
			if err != nil {
				return nil, fmt.Errorf("relative path %s: %w", m, err)
			}
			res := Resolution{
				Resource: filepath.ToSlash(rel),
				Source:   string(c.kind),
				Path:     m,
				IsDir:    false,
			}
			resolved = append(resolved, res)
		}
		if len(resolved) > 0 {
			r.records = append(r.records, resolved...)
			return resolved, nil
		}
	}

	if r.embedFS != nil {
		matches, err := fs.Glob(r.embedFS, clean)
		if err != nil {
			return nil, fmt.Errorf("glob embed %s: %w", clean, err)
		}
		if len(matches) == 0 {
			return nil, nil
		}
		resolved := make([]Resolution, 0, len(matches))
		for _, m := range matches {
			info, err := fs.Stat(r.embedFS, m)
			if errors.Is(err, fs.ErrNotExist) {
				continue
			}
			if err != nil {
				return nil, fmt.Errorf("stat embed %s: %w", m, err)
			}
			if info.IsDir() {
				continue
			}
			resolved = append(resolved, Resolution{
				Resource: m,
				Source:   string(sourceEmbed),
				Path:     "embed:" + m,
				IsDir:    false,
			})
		}
		if len(resolved) > 0 {
			r.records = append(r.records, resolved...)
			return resolved, nil
		}
	}

	return nil, nil
}

func (r *Resolver) ResolveOptionalFile(resource string) (Resolution, bool, error) {
	return r.resolve(resource, false)
}

func (r *Resolver) ResolveDir(resource string) (Resolution, error) {
	res, ok, err := r.resolve(resource, true)
	if err != nil {
		return Resolution{}, err
	}
	if !ok {
		return Resolution{}, r.missingErr(resource)
	}
	return res, nil
}

func (r *Resolver) ResolveOptionalDir(resource string) (Resolution, bool, error) {
	return r.resolve(resource, true)
}

func (r *Resolver) ReadDir(resource string) ([]fs.DirEntry, Resolution, error) {
	res, err := r.ResolveDir(resource)
	if err != nil {
		return nil, Resolution{}, err
	}
	if res.Source == string(sourceEmbed) {
		entries, err := fs.ReadDir(r.embedFS, res.Resource)
		if err != nil {
			return nil, Resolution{}, fmt.Errorf("read embed dir %s: %w", res.Resource, err)
		}
		return entries, res, nil
	}
	entries, err := os.ReadDir(res.Path)
	if err != nil {
		return nil, Resolution{}, fmt.Errorf("read dir %s: %w", res.Path, err)
	}
	return entries, res, nil
}

func (r *Resolver) Diagnostics() []Resolution {
	out := make([]Resolution, len(r.records))
	copy(out, r.records)
	return out
}

func (r *Resolver) resolve(resource string, wantDir bool) (Resolution, bool, error) {
	clean, err := sanitizeResourcePath(resource)
	if err != nil {
		return Resolution{}, false, err
	}

	for _, c := range r.candidates {
		full := filepath.Join(c.root, filepath.FromSlash(clean))
		info, err := os.Stat(full)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return Resolution{}, false, fmt.Errorf("stat %s: %w", full, err)
		}
		if info.IsDir() != wantDir {
			continue
		}
		res := Resolution{
			Resource: clean,
			Source:   string(c.kind),
			Path:     full,
			IsDir:    wantDir,
		}
		r.records = append(r.records, res)
		return res, true, nil
	}

	if r.embedFS != nil {
		info, err := fs.Stat(r.embedFS, clean)
		if errors.Is(err, fs.ErrNotExist) {
			return Resolution{}, false, nil
		}
		if err != nil {
			return Resolution{}, false, fmt.Errorf("stat embed %s: %w", clean, err)
		}
		if info.IsDir() == wantDir {
			res := Resolution{
				Resource: clean,
				Source:   string(sourceEmbed),
				Path:     "embed:" + clean,
				IsDir:    wantDir,
			}
			r.records = append(r.records, res)
			return res, true, nil
		}
	}

	return Resolution{}, false, nil
}

func (r *Resolver) openResolvedFile(res Resolution) (io.ReadCloser, error) {
	if res.IsDir {
		return nil, fmt.Errorf("resource is directory: %s", res.Path)
	}
	if res.Source == string(sourceEmbed) {
		f, err := r.embedFS.Open(res.Resource)
		if err != nil {
			return nil, fmt.Errorf("open embed resource %s: %w", res.Resource, err)
		}
		rc, ok := f.(io.ReadCloser)
		if !ok {
			_ = f.Close()
			return nil, fmt.Errorf("embed resource is not readable: %s", res.Resource)
		}
		return rc, nil
	}
	f, err := os.Open(res.Path)
	if err != nil {
		return nil, fmt.Errorf("open resource %s: %w", res.Path, err)
	}
	return f, nil
}

func (r *Resolver) missingErr(resource string) error {
	searched := make([]string, 0, len(r.candidates)+1)
	clean := filepath.ToSlash(filepath.Clean(resource))
	for _, c := range r.candidates {
		searched = append(searched, fmt.Sprintf("%s:%s", c.kind, filepath.Join(c.root, filepath.FromSlash(clean))))
	}
	if r.embedFS != nil {
		searched = append(searched, "embed:"+clean)
	}
	return &MissingResourceError{
		Resource: clean,
		Searched: searched,
	}
}

func sanitizeResourcePath(resource string) (string, error) {
	clean := filepath.ToSlash(filepath.Clean(resource))
	switch {
	case clean == "", clean == ".", clean == "/", clean == "..":
		return "", fmt.Errorf("invalid resource path: %s", resource)
	case strings.HasPrefix(clean, "/"):
		return "", fmt.Errorf("absolute resource path is not allowed: %s", resource)
	case strings.HasPrefix(clean, "../"):
		return "", fmt.Errorf("resource path escapes base directory: %s", resource)
	case filepath.VolumeName(clean) != "":
		return "", fmt.Errorf("resource path must not include volume name: %s", resource)
	}
	return clean, nil
}
