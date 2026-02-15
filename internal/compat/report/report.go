package report

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	compat "github.com/yuanying/azash/internal/compat"
)

func WriteJSON(run compat.RunReport, path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(run, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}

func WriteMarkdown(run compat.RunReport, path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	var sb strings.Builder
	sb.WriteString("# Compatibility Report\n\n")
	sb.WriteString(fmt.Sprintf("- mode: `%s`\n", run.Run.Mode))
	sb.WriteString(fmt.Sprintf("- started_at: `%s`\n", run.Run.StartedAt.Format("2006-01-02T15:04:05Z07:00")))
	sb.WriteString(fmt.Sprintf("- finished_at: `%s`\n", run.Run.FinishedAt.Format("2006-01-02T15:04:05Z07:00")))
	sb.WriteString(fmt.Sprintf("- sample_count: `%d`\n", run.SampleCount))
	sb.WriteString(fmt.Sprintf("- files_compared: `%d`\n", run.DiffSummary.FilesCompared))
	sb.WriteString(fmt.Sprintf("- files_changed: `%d`\n", run.DiffSummary.FilesChanged))
	sb.WriteString(fmt.Sprintf("- node_diffs: `%d`, text_diffs: `%d`, attr_diffs: `%d`\n\n", run.DiffSummary.NodeDiffs, run.DiffSummary.TextDiffs, run.DiffSummary.AttrDiffs))

	if len(run.Categories) > 0 {
		sb.WriteString("## Categories\n\n")
		keys := make([]string, 0, len(run.Categories))
		for k := range run.Categories {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			sb.WriteString(fmt.Sprintf("- %s: %d\n", k, run.Categories[k]))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("## Samples\n\n")
	for _, s := range run.Samples {
		sb.WriteString(fmt.Sprintf("### %s\n\n", s.SampleID))
		sb.WriteString(fmt.Sprintf("- status: `%s`\n", s.Status))
		if s.JavaEPUB != "" {
			sb.WriteString(fmt.Sprintf("- java_epub: `%s`\n", s.JavaEPUB))
		}
		if s.GoEPUB != "" {
			sb.WriteString(fmt.Sprintf("- go_epub: `%s`\n", s.GoEPUB))
		}
		sb.WriteString(fmt.Sprintf("- files_changed: `%d`\n", s.DiffSummary.FilesChanged))
		if len(s.Failures) > 0 {
			sb.WriteString("- failures:\n")
			for _, f := range s.Failures {
				sb.WriteString(fmt.Sprintf("  - [%s] %s\n", f.Stage, f.Message))
			}
		}
		sb.WriteString("\n")
	}

	if len(run.Failures) > 0 {
		sb.WriteString("## Failures\n\n")
		for _, f := range run.Failures {
			sb.WriteString(fmt.Sprintf("- sample=%s stage=%s hard=%t: %s\n", f.SampleID, f.Stage, f.Hard, f.Message))
		}
	}

	return os.WriteFile(path, []byte(sb.String()), 0o644)
}
