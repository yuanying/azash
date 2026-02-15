package diff

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	compat "github.com/yuanying/azash/internal/compat"
)

type Options struct{}

func CompareTrees(javaRoot, goRoot string, _ Options) (compat.DiffSummary, []compat.FileDiff, map[string]int, error) {
	javaFiles, err := listFiles(javaRoot)
	if err != nil {
		return compat.DiffSummary{}, nil, nil, err
	}
	goFiles, err := listFiles(goRoot)
	if err != nil {
		return compat.DiffSummary{}, nil, nil, err
	}

	allPaths := map[string]struct{}{}
	for p := range javaFiles {
		allPaths[p] = struct{}{}
	}
	for p := range goFiles {
		allPaths[p] = struct{}{}
	}
	paths := make([]string, 0, len(allPaths))
	for p := range allPaths {
		paths = append(paths, p)
	}
	sort.Strings(paths)

	summary := compat.DiffSummary{}
	categories := map[string]int{}
	fileDiffs := make([]compat.FileDiff, 0, len(paths))

	for _, rel := range paths {
		summary.FilesCompared++
		javaPath, jOK := javaFiles[rel]
		goPath, gOK := goFiles[rel]
		fd := compat.FileDiff{Path: rel}
		fd.Category = classify(rel, nil)

		if !jOK || !gOK {
			fd.Changed = true
			fd.Message = "file missing on one side"
			summary.FilesChanged++
			categories[fd.Category]++
			fileDiffs = append(fileDiffs, fd)
			continue
		}

		jb, err := os.ReadFile(javaPath)
		if err != nil {
			return summary, fileDiffs, categories, fmt.Errorf("read java file %s: %w", rel, err)
		}
		gb, err := os.ReadFile(goPath)
		if err != nil {
			return summary, fileDiffs, categories, fmt.Errorf("read go file %s: %w", rel, err)
		}

		ext := strings.ToLower(filepath.Ext(rel))
		switch ext {
		case ".xml", ".xhtml", ".html", ".opf", ".ncx":
			n, t, a, changed := compareXML(jb, gb)
			fd.NodeDiffs = n
			fd.TextDiffs = t
			fd.AttrDiffs = a
			fd.Changed = changed
			fd.Category = classify(rel, jb)
			summary.NodeDiffs += n
			summary.TextDiffs += t
			summary.AttrDiffs += a
		default:
			fd.Changed = !bytes.Equal(jb, gb)
		}

		if fd.Changed {
			summary.FilesChanged++
			categories[fd.Category]++
		}
		fileDiffs = append(fileDiffs, fd)
	}
	return summary, fileDiffs, categories, nil
}

func listFiles(root string) (map[string]string, error) {
	out := map[string]string{}
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		out[filepath.ToSlash(rel)] = path
		return nil
	})
	return out, err
}

type xmlTok struct {
	Kind  string
	Name  string
	Attrs map[string]string
	Text  string
}

func compareXML(a, b []byte) (nodeDiffs, textDiffs, attrDiffs int, changed bool) {
	ta := tokenizeXML(a)
	tb := tokenizeXML(b)
	max := len(ta)
	if len(tb) > max {
		max = len(tb)
	}
	for i := 0; i < max; i++ {
		if i >= len(ta) || i >= len(tb) {
			nodeDiffs++
			continue
		}
		x := ta[i]
		y := tb[i]
		if x.Kind != y.Kind || x.Name != y.Name {
			nodeDiffs++
		}
		if x.Kind == "text" && x.Text != y.Text {
			textDiffs++
		}
		if x.Kind == "start" {
			attrDiffs += attrMapDiff(x.Attrs, y.Attrs)
		}
	}
	changed = nodeDiffs > 0 || textDiffs > 0 || attrDiffs > 0
	return
}

func tokenizeXML(b []byte) []xmlTok {
	dec := xml.NewDecoder(bytes.NewReader(b))
	dec.Strict = false
	out := []xmlTok{}
	for {
		tok, err := dec.Token()
		if err != nil {
			break
		}
		switch t := tok.(type) {
		case xml.StartElement:
			attrs := map[string]string{}
			for _, a := range t.Attr {
				attrs[a.Name.Space+":"+a.Name.Local] = a.Value
			}
			out = append(out, xmlTok{Kind: "start", Name: t.Name.Space + ":" + t.Name.Local, Attrs: attrs})
		case xml.EndElement:
			out = append(out, xmlTok{Kind: "end", Name: t.Name.Space + ":" + t.Name.Local})
		case xml.CharData:
			s := strings.TrimSpace(string(t))
			if s != "" {
				out = append(out, xmlTok{Kind: "text", Text: s})
			}
		}
	}
	return out
}

func attrMapDiff(a, b map[string]string) int {
	d := 0
	seen := map[string]struct{}{}
	for k, va := range a {
		seen[k] = struct{}{}
		if vb, ok := b[k]; !ok || va != vb {
			d++
		}
	}
	for k := range b {
		if _, ok := seen[k]; !ok {
			d++
		}
	}
	return d
}

func classify(path string, content []byte) string {
	p := strings.ToLower(path)
	if strings.Contains(p, "mimetype") || strings.Contains(p, "meta-inf/") {
		return "packaging"
	}
	if strings.HasSuffix(p, ".opf") || strings.Contains(p, "container.xml") {
		return "metadata"
	}
	if strings.Contains(p, "nav") || strings.Contains(p, "toc") {
		return "chapter"
	}
	if strings.Contains(p, "image") || strings.HasSuffix(p, ".png") || strings.HasSuffix(p, ".jpg") || strings.HasSuffix(p, ".jpeg") {
		return "image"
	}
	if bytes.Contains(bytes.ToLower(content), []byte("<ruby")) {
		return "ruby"
	}
	if strings.Contains(p, "chuki") || bytes.Contains(content, []byte("注記")) {
		return "chuki"
	}
	if strings.Contains(p, "encoding") {
		return "encoding"
	}
	return "unknown"
}
