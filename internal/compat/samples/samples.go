package samples

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

	compat "github.com/yuanying/azash/internal/compat"
)

var requiredColumns = []string{
	"id",
	"title",
	"source_url",
	"downloaded_at",
	"input_file",
	"expected_encoding",
	"tags",
	"notes",
}

func LoadCSV(path string) ([]compat.Sample, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open samples csv: %w", err)
	}
	defer func() { _ = f.Close() }()

	r := csv.NewReader(f)
	r.FieldsPerRecord = -1

	header, err := r.Read()
	if err != nil {
		return nil, fmt.Errorf("read header: %w", err)
	}
	index := map[string]int{}
	for i, col := range header {
		index[strings.TrimSpace(col)] = i
	}
	for _, col := range requiredColumns {
		if _, ok := index[col]; !ok {
			return nil, fmt.Errorf("missing required column %q", col)
		}
	}

	var out []compat.Sample
	ids := map[string]struct{}{}
	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read record: %w", err)
		}
		if len(rec) == 0 {
			continue
		}
		s := compat.Sample{
			ID:               field(rec, index["id"]),
			Title:            field(rec, index["title"]),
			SourceURL:        field(rec, index["source_url"]),
			DownloadedAt:     field(rec, index["downloaded_at"]),
			InputFile:        field(rec, index["input_file"]),
			ExpectedEncoding: field(rec, index["expected_encoding"]),
			Notes:            field(rec, index["notes"]),
		}
		for _, t := range strings.Split(field(rec, index["tags"]), ",") {
			t = strings.TrimSpace(strings.Trim(t, `"`))
			if t != "" {
				s.Tags = append(s.Tags, t)
			}
		}
		if s.ID == "" {
			return nil, fmt.Errorf("record with empty id")
		}
		if _, ok := ids[s.ID]; ok {
			return nil, fmt.Errorf("duplicate sample id %q", s.ID)
		}
		ids[s.ID] = struct{}{}
		if s.ExpectedEncoding == "" {
			s.ExpectedEncoding = "AUTO"
		}
		out = append(out, s)
	}
	return out, nil
}

func field(rec []string, idx int) string {
	if idx < 0 || idx >= len(rec) {
		return ""
	}
	return strings.TrimSpace(rec[idx])
}
