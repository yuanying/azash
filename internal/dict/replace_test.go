package dict

import (
	"os"
	"strings"
	"testing"
)

func TestLoadReplaceMap_SingleChar(t *testing.T) {
	input := "－\t―\n"
	m, err := LoadReplaceMap(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	v, ok := m.Single['－']
	if !ok {
		t.Fatal("－ not found")
	}
	if v != "―" {
		t.Errorf("want ―, got %s", v)
	}
}

func TestLoadReplaceMap_DoubleChar(t *testing.T) {
	input := "。」\t」\n"
	m, err := LoadReplaceMap(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	v, ok := m.Double["。」"]
	if !ok {
		t.Fatal("。」 not found")
	}
	if v != "」" {
		t.Errorf("want 」, got %s", v)
	}
}

func TestLoadReplaceMap_ThreeCharSkipped(t *testing.T) {
	input := "あいう\t置換後\n"
	m, err := LoadReplaceMap(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(m.Single) != 0 || len(m.Double) != 0 {
		t.Error("3+ char entries should be skipped")
	}
}

func TestLoadReplaceMap_CommentAndEmpty(t *testing.T) {
	input := "### comment\n\n－\t―\n"
	m, err := LoadReplaceMap(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(m.Single) != 1 {
		t.Errorf("want 1 entry, got %d", len(m.Single))
	}
}

func TestLoadReplaceMap_NonBMPChar(t *testing.T) {
	// BMP外文字はUTF-16でサロゲートペア（2ユニット）→ Double に分類
	input := "𠀋\t代替\n"
	m, err := LoadReplaceMap(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(m.Single) != 0 {
		t.Errorf("non-BMP char should not be in Single: got %d", len(m.Single))
	}
	v, ok := m.Double["𠀋"]
	if !ok {
		t.Fatal("𠀋 not found in Double")
	}
	if v != "代替" {
		t.Errorf("want 代替, got %s", v)
	}
}

func TestLoadReplaceMap_TrailingTab(t *testing.T) {
	// Java の split("\t") は末尾空要素を落とすため、"－\t" は1要素→スキップ
	input := "－\t\n"
	m, err := LoadReplaceMap(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(m.Single) != 0 || len(m.Double) != 0 {
		t.Error("trailing tab line should be skipped (Java split compat)")
	}
}

func TestLoadReplaceMap_RealFile(t *testing.T) {
	f, err := os.Open("../../temp/AozoraEpub3/replace_sample.txt")
	if err != nil {
		t.Skip("replace_sample.txt not found:", err)
	}
	defer func() { _ = f.Close() }()

	m, err := LoadReplaceMap(f)
	if err != nil {
		t.Fatal(err)
	}
	if len(m.Single)+len(m.Double) < 1 {
		t.Error("expected at least 1 entry from replace_sample.txt")
	}
}
