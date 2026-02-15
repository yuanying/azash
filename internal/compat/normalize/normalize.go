package normalize

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type Summary struct {
	Files int `json:"files"`
}

var (
	uuidRe     = regexp.MustCompile(`(?i)\b[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}\b`)
	datetimeRe = regexp.MustCompile(`\b\d{4}-\d{2}-\d{2}(?:[tT ]\d{2}:\d{2}:\d{2}(?:\.\d+)?(?:Z|[+-]\d{2}:?\d{2})?)?\b`)
	multiSpace = regexp.MustCompile(`\s+`)
)

func NormalizeTree(srcRoot, dstRoot string) (Summary, error) {
	if err := os.RemoveAll(dstRoot); err != nil {
		return Summary{}, err
	}
	if err := os.MkdirAll(dstRoot, 0o755); err != nil {
		return Summary{}, err
	}

	summary := Summary{}
	err := filepath.WalkDir(srcRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(srcRoot, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		dst := filepath.Join(dstRoot, rel)
		if d.IsDir() {
			return os.MkdirAll(dst, 0o755)
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		n := normalizeContent(path, b)
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(dst, n, 0o644); err != nil {
			return err
		}
		summary.Files++
		return nil
	})
	if err != nil {
		return Summary{}, fmt.Errorf("normalize tree: %w", err)
	}
	return summary, nil
}

func normalizeContent(path string, b []byte) []byte {
	ext := strings.ToLower(filepath.Ext(path))
	b = uuidRe.ReplaceAll(b, []byte("UUID_NORMALIZED"))
	b = datetimeRe.ReplaceAll(b, []byte("DATE_NORMALIZED"))

	switch ext {
	case ".xml", ".xhtml", ".html", ".opf", ".ncx":
		if out, err := canonicalizeXML(b); err == nil {
			return out
		}
	case ".css":
		s := strings.ReplaceAll(string(b), "\r\n", "\n")
		lines := strings.Split(s, "\n")
		for i := range lines {
			lines[i] = strings.TrimSpace(lines[i])
		}
		return []byte(strings.Join(lines, "\n"))
	}
	return b
}

func canonicalizeXML(b []byte) ([]byte, error) {
	dec := xml.NewDecoder(bytes.NewReader(b))
	dec.Strict = false
	var out bytes.Buffer
	enc := xml.NewEncoder(&out)
	for {
		tok, err := dec.Token()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, err
		}
		switch t := tok.(type) {
		case xml.StartElement:
			sort.Slice(t.Attr, func(i, j int) bool {
				ai := t.Attr[i].Name.Space + ":" + t.Attr[i].Name.Local
				aj := t.Attr[j].Name.Space + ":" + t.Attr[j].Name.Local
				return ai < aj
			})
			if err := enc.EncodeToken(t); err != nil {
				return nil, err
			}
		case xml.EndElement:
			if err := enc.EncodeToken(t); err != nil {
				return nil, err
			}
		case xml.CharData:
			s := strings.TrimSpace(multiSpace.ReplaceAllString(string(t), " "))
			if s == "" {
				continue
			}
			if err := enc.EncodeToken(xml.CharData([]byte(s))); err != nil {
				return nil, err
			}
		case xml.Comment:
			continue
		default:
			if err := enc.EncodeToken(tok); err != nil {
				return nil, err
			}
		}
	}
	if err := enc.Flush(); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}
