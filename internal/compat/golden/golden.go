package golden

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"

	compat "github.com/yuanying/azash/internal/compat"
)

type Options struct {
	InputDir  string
	OutputDir string
	LogDir    string
	JavaCmd   string
	JavaCP    string
	JavaMain  string
	INIPath   string
	Ext       string
}

func Generate(ctx context.Context, samples []compat.Sample, opts Options) ([]compat.Artifact, []compat.StageFailure, error) {
	if opts.Ext == "" {
		opts.Ext = ".epub"
	}
	if err := os.MkdirAll(opts.OutputDir, 0o755); err != nil {
		return nil, nil, fmt.Errorf("mkdir output: %w", err)
	}
	if opts.LogDir != "" {
		if err := os.MkdirAll(opts.LogDir, 0o755); err != nil {
			return nil, nil, fmt.Errorf("mkdir log dir: %w", err)
		}
	}

	var artifacts []compat.Artifact
	var failures []compat.StageFailure
	for _, s := range samples {
		inputPath := filepath.Join(opts.InputDir, s.InputFile)
		if _, err := os.Stat(inputPath); err != nil {
			failures = append(failures, compat.StageFailure{Stage: "golden", SampleID: s.ID, Message: fmt.Sprintf("input missing: %v", err), Hard: true})
			continue
		}
		tmpOut, err := os.MkdirTemp("", "azash-compat-java-out-")
		if err != nil {
			return artifacts, failures, fmt.Errorf("mktemp: %w", err)
		}

		enc := s.ExpectedEncoding
		if enc == "" {
			enc = "AUTO"
		}
		args := []string{"-cp", opts.JavaCP, opts.JavaMain, "-i", opts.INIPath, "-enc", enc, "-ext", opts.Ext, "-d", tmpOut, inputPath}
		cmd := exec.CommandContext(ctx, opts.JavaCmd, args...)
		out, err := cmd.CombinedOutput()

		if opts.LogDir != "" {
			_ = os.WriteFile(filepath.Join(opts.LogDir, s.ID+".java.log"), out, 0o644)
		}

		if err != nil {
			_ = os.RemoveAll(tmpOut)
			failures = append(failures, compat.StageFailure{Stage: "golden", SampleID: s.ID, Message: fmt.Sprintf("java command failed: %v", err), Hard: true})
			continue
		}

		epubs, err := filepath.Glob(filepath.Join(tmpOut, "*.epub"))
		if err != nil || len(epubs) == 0 {
			_ = os.RemoveAll(tmpOut)
			failures = append(failures, compat.StageFailure{Stage: "golden", SampleID: s.ID, Message: "no epub produced", Hard: true})
			continue
		}
		sort.Slice(epubs, func(i, j int) bool {
			a, _ := os.Stat(epubs[i])
			b, _ := os.Stat(epubs[j])
			return a.ModTime().After(b.ModTime())
		})

		src := epubs[0]
		dst := filepath.Join(opts.OutputDir, s.ID+".java.epub")
		if err := copyFile(src, dst); err != nil {
			_ = os.RemoveAll(tmpOut)
			failures = append(failures, compat.StageFailure{Stage: "golden", SampleID: s.ID, Message: fmt.Sprintf("copy epub: %v", err), Hard: true})
			continue
		}
		_ = os.RemoveAll(tmpOut)
		artifacts = append(artifacts, compat.Artifact{SampleID: s.ID, Path: dst, Kind: "java_epub"})
	}
	return artifacts, failures, nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() { _ = in.Close() }()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() { _ = out.Close() }()
	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Close()
}
