package pipeline

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	compat "github.com/yuanying/azash/internal/compat"
	"github.com/yuanying/azash/internal/compat/diff"
	"github.com/yuanying/azash/internal/compat/epub"
	"github.com/yuanying/azash/internal/compat/golden"
	"github.com/yuanying/azash/internal/compat/normalize"
	"github.com/yuanying/azash/internal/compat/report"
	"github.com/yuanying/azash/internal/compat/samples"
)

type Mode string

const (
	ModeGoldenJava Mode = "golden-java"
	ModeCompare    Mode = "compare"
	ModeRun        Mode = "run"
)

type Options struct {
	SamplesFile string
	InputDir    string
	GoldenDir   string
	GoOutDir    string
	WorkDir     string
	ReportDir   string
	LogDir      string

	SampleIDs map[string]struct{}

	JavaCmd  string
	JavaCP   string
	JavaMain string
	INIPath  string

	FailOnDiff bool
}

func DefaultOptions() Options {
	return Options{
		SamplesFile: "testdata/compat/samples.csv",
		InputDir:    "testdata/compat/input",
		GoldenDir:   "testdata/compat/golden/java",
		GoOutDir:    "testdata/compat/output/go",
		WorkDir:     "testdata/compat/unpacked",
		ReportDir:   "testdata/compat/reports/latest",
		LogDir:      "testdata/compat/reports/latest/logs",
		JavaCmd:     "java",
		JavaCP:      "temp/AozoraEpub3-bin/AozoraEpub3.jar",
		JavaMain:    "AozoraEpub3",
		INIPath:     "temp/AozoraEpub3-bin/presets/reader.ini",
	}
}

func Execute(ctx context.Context, mode Mode, opts Options) (compat.RunReport, error) {
	start := time.Now()
	run := compat.RunReport{
		Run: compat.RunMeta{
			StartedAt: start,
			Tool:      "azash-compat",
			Version:   "dev",
			Mode:      string(mode),
		},
		Categories: map[string]int{},
		Artifacts:  map[string]any{},
	}

	smpls, err := samples.LoadCSV(opts.SamplesFile)
	if err != nil {
		return run, err
	}
	smpls = filterSamples(smpls, opts.SampleIDs)
	run.SampleCount = len(smpls)

	if err := os.MkdirAll(opts.ReportDir, 0o755); err != nil {
		return run, err
	}
	if err := os.MkdirAll(opts.WorkDir, 0o755); err != nil {
		return run, err
	}

	if mode == ModeGoldenJava || mode == ModeRun {
		artifacts, fails, err := golden.Generate(ctx, smpls, golden.Options{
			InputDir:  opts.InputDir,
			OutputDir: opts.GoldenDir,
			LogDir:    opts.LogDir,
			JavaCmd:   opts.JavaCmd,
			JavaCP:    opts.JavaCP,
			JavaMain:  opts.JavaMain,
			INIPath:   opts.INIPath,
			Ext:       ".epub",
		})
		if err != nil {
			return run, err
		}
		run.Artifacts["golden"] = artifacts
		run.Failures = append(run.Failures, fails...)
	}

	if mode == ModeCompare || mode == ModeRun {
		for _, s := range smpls {
			sr := compat.SampleReport{SampleID: s.ID, Status: "matched", Categories: map[string]int{}}
			javaEPUB := filepath.Join(opts.GoldenDir, s.ID+".java.epub")
			goEPUB := filepath.Join(opts.GoOutDir, s.ID+".go.epub")
			sr.JavaEPUB = javaEPUB
			sr.GoEPUB = goEPUB

			if _, err := os.Stat(javaEPUB); err != nil {
				f := compat.StageFailure{Stage: "compare", SampleID: s.ID, Message: "java epub missing", Hard: true}
				sr.Status = "error"
				sr.Failures = append(sr.Failures, f)
				run.Failures = append(run.Failures, f)
				run.Samples = append(run.Samples, sr)
				continue
			}
			if _, err := os.Stat(goEPUB); err != nil {
				f := compat.StageFailure{Stage: "go-generate", SampleID: s.ID, Message: "go epub missing; compare skipped", Hard: false}
				sr.Status = "skipped"
				sr.Failures = append(sr.Failures, f)
				run.Failures = append(run.Failures, f)
				run.Samples = append(run.Samples, sr)
				continue
			}

			javaRaw := filepath.Join(opts.WorkDir, s.ID, "java", "raw")
			goRaw := filepath.Join(opts.WorkDir, s.ID, "go", "raw")
			javaNorm := filepath.Join(opts.WorkDir, s.ID, "java", "normalized")
			goNorm := filepath.Join(opts.WorkDir, s.ID, "go", "normalized")
			if err := epub.Unpack(javaEPUB, javaRaw); err != nil {
				f := compat.StageFailure{Stage: "unpack", SampleID: s.ID, Message: err.Error(), Hard: true}
				sr.Status = "error"
				sr.Failures = append(sr.Failures, f)
				run.Failures = append(run.Failures, f)
				run.Samples = append(run.Samples, sr)
				continue
			}
			if err := epub.Unpack(goEPUB, goRaw); err != nil {
				f := compat.StageFailure{Stage: "unpack", SampleID: s.ID, Message: err.Error(), Hard: true}
				sr.Status = "error"
				sr.Failures = append(sr.Failures, f)
				run.Failures = append(run.Failures, f)
				run.Samples = append(run.Samples, sr)
				continue
			}
			if _, err := normalize.NormalizeTree(javaRaw, javaNorm); err != nil {
				f := compat.StageFailure{Stage: "normalize", SampleID: s.ID, Message: err.Error(), Hard: true}
				sr.Status = "error"
				sr.Failures = append(sr.Failures, f)
				run.Failures = append(run.Failures, f)
				run.Samples = append(run.Samples, sr)
				continue
			}
			if _, err := normalize.NormalizeTree(goRaw, goNorm); err != nil {
				f := compat.StageFailure{Stage: "normalize", SampleID: s.ID, Message: err.Error(), Hard: true}
				sr.Status = "error"
				sr.Failures = append(sr.Failures, f)
				run.Failures = append(run.Failures, f)
				run.Samples = append(run.Samples, sr)
				continue
			}

			ds, fds, cats, err := diff.CompareTrees(javaNorm, goNorm, diff.Options{})
			if err != nil {
				f := compat.StageFailure{Stage: "diff", SampleID: s.ID, Message: err.Error(), Hard: true}
				sr.Status = "error"
				sr.Failures = append(sr.Failures, f)
				run.Failures = append(run.Failures, f)
				run.Samples = append(run.Samples, sr)
				continue
			}
			sr.DiffSummary = ds
			sr.FileDiffs = fds
			sr.Categories = cats
			if ds.FilesChanged > 0 {
				sr.Status = "different"
			}
			mergeSummary(&run.DiffSummary, ds)
			mergeCategories(run.Categories, cats)
			run.Samples = append(run.Samples, sr)
		}
	}

	run.Run.FinishedAt = time.Now()
	sort.Slice(run.Samples, func(i, j int) bool { return run.Samples[i].SampleID < run.Samples[j].SampleID })

	jsonPath := filepath.Join(opts.ReportDir, "report.json")
	mdPath := filepath.Join(opts.ReportDir, "summary.md")
	if err := report.WriteJSON(run, jsonPath); err != nil {
		return run, err
	}
	if err := report.WriteMarkdown(run, mdPath); err != nil {
		return run, err
	}

	if hardFailureExists(run.Failures) {
		return run, errors.New("compat pipeline has hard failures")
	}
	if opts.FailOnDiff && run.DiffSummary.FilesChanged > 0 {
		return run, fmt.Errorf("diff detected: %d files changed", run.DiffSummary.FilesChanged)
	}
	return run, nil
}

func filterSamples(in []compat.Sample, ids map[string]struct{}) []compat.Sample {
	if len(ids) == 0 {
		return in
	}
	out := make([]compat.Sample, 0, len(ids))
	for _, s := range in {
		if _, ok := ids[s.ID]; ok {
			out = append(out, s)
		}
	}
	return out
}

func mergeSummary(dst *compat.DiffSummary, src compat.DiffSummary) {
	dst.FilesCompared += src.FilesCompared
	dst.FilesChanged += src.FilesChanged
	dst.NodeDiffs += src.NodeDiffs
	dst.TextDiffs += src.TextDiffs
	dst.AttrDiffs += src.AttrDiffs
}

func mergeCategories(dst map[string]int, src map[string]int) {
	for k, v := range src {
		dst[k] += v
	}
}

func hardFailureExists(f []compat.StageFailure) bool {
	for _, x := range f {
		if x.Hard {
			return true
		}
	}
	return false
}
