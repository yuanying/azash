package compat

import "time"

// Sample describes one compatibility test input.
type Sample struct {
	ID               string
	Title            string
	SourceURL        string
	DownloadedAt     string
	InputFile        string
	ExpectedEncoding string
	Tags             []string
	Notes            string
}

// Artifact describes an output artifact file.
type Artifact struct {
	SampleID string `json:"sample_id"`
	Path     string `json:"path"`
	Kind     string `json:"kind"`
}

// StageFailure records one stage failure/warning during pipeline execution.
type StageFailure struct {
	Stage    string `json:"stage"`
	SampleID string `json:"sample_id,omitempty"`
	Message  string `json:"message"`
	Hard     bool   `json:"hard"`
}

// DiffSummary accumulates difference counters.
type DiffSummary struct {
	FilesCompared int `json:"files_compared"`
	FilesChanged  int `json:"files_changed"`
	NodeDiffs     int `json:"node_diffs"`
	TextDiffs     int `json:"text_diffs"`
	AttrDiffs     int `json:"attr_diffs"`
}

// FileDiff describes per-file diff information.
type FileDiff struct {
	Path      string `json:"path"`
	Category  string `json:"category"`
	Changed   bool   `json:"changed"`
	NodeDiffs int    `json:"node_diffs,omitempty"`
	TextDiffs int    `json:"text_diffs,omitempty"`
	AttrDiffs int    `json:"attr_diffs,omitempty"`
	Message   string `json:"message,omitempty"`
}

// SampleReport aggregates one sample run.
type SampleReport struct {
	SampleID    string                 `json:"sample_id"`
	Status      string                 `json:"status"`
	JavaEPUB    string                 `json:"java_epub,omitempty"`
	GoEPUB      string                 `json:"go_epub,omitempty"`
	DiffSummary DiffSummary            `json:"diff_summary"`
	Categories  map[string]int         `json:"categories,omitempty"`
	FileDiffs   []FileDiff             `json:"file_diffs,omitempty"`
	Failures    []StageFailure         `json:"failures,omitempty"`
	Artifacts   map[string]interface{} `json:"artifacts,omitempty"`
}

// RunMeta describes run metadata.
type RunMeta struct {
	StartedAt  time.Time `json:"started_at"`
	FinishedAt time.Time `json:"finished_at"`
	Tool       string    `json:"tool"`
	Version    string    `json:"version"`
	Mode       string    `json:"mode"`
}

// RunReport is the top-level JSON report.
type RunReport struct {
	Run         RunMeta        `json:"run"`
	SampleCount int            `json:"sample_count"`
	Samples     []SampleReport `json:"samples"`
	DiffSummary DiffSummary    `json:"diff_summary"`
	Categories  map[string]int `json:"categories"`
	Failures    []StageFailure `json:"failures,omitempty"`
	Artifacts   map[string]any `json:"artifacts,omitempty"`
}
