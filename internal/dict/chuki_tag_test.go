package dict

import (
	"os"
	"strings"
	"testing"
)

func TestLoadChukiTagMap_Basic(t *testing.T) {
	input := "改丁\t\t\tP\n改ページ\t\t\tP\n"
	m, err := LoadChukiTagMap(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	tags, ok := m.Tags["改丁"]
	if !ok {
		t.Fatal("改丁 not found")
	}
	if len(tags) != 1 || tags[0] != "" {
		t.Errorf("改丁 tags: want [\"\"], got %v", tags)
	}
	if _, ok := m.FlagPageBreak["改丁"]; !ok {
		t.Error("改丁 should have PageBreak flag")
	}
}

func TestLoadChukiTagMap_EndTag(t *testing.T) {
	input := "ページ左下\t<div class=\"btm\">\t</div>\tL\n"
	m, err := LoadChukiTagMap(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	tags := m.Tags["ページ左下"]
	if len(tags) != 2 {
		t.Fatalf("want 2 tags, got %d", len(tags))
	}
	if tags[0] != "<div class=\"btm\">" {
		t.Errorf("start tag: want <div class=\"btm\">, got %s", tags[0])
	}
	if tags[1] != "</div>" {
		t.Errorf("end tag: want </div>, got %s", tags[1])
	}
	if _, ok := m.FlagPageBreak["ページ左下"]; !ok {
		t.Error("should have PageBreak flag")
	}
	if _, ok := m.FlagBottom["ページ左下"]; !ok {
		t.Error("should have Bottom flag")
	}
}

func TestLoadChukiTagMap_Flags(t *testing.T) {
	input := strings.Join([]string{
		"ここから大見出し\t<div class=\"font-1em50\">\t\t1",
		"ルビ開始\t<ruby>\t\t2",
		"ルビ終了\t</ruby>\t\t3",
		"ページの左右中央\t\t\tM",
		"訓点\t\t\tK",
	}, "\n") + "\n"
	m, err := LoadChukiTagMap(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := m.FlagNoBr["ここから大見出し"]; !ok {
		t.Error("should have NoBr flag")
	}
	if _, ok := m.FlagNoRubyStart["ルビ開始"]; !ok {
		t.Error("should have NoRubyStart flag")
	}
	if _, ok := m.FlagNoRubyEnd["ルビ終了"]; !ok {
		t.Error("should have NoRubyEnd flag")
	}
	if _, ok := m.FlagPageBreak["ページの左右中央"]; !ok {
		t.Error("should have PageBreak flag for M")
	}
	if _, ok := m.FlagMiddle["ページの左右中央"]; !ok {
		t.Error("should have Middle flag")
	}
	if _, ok := m.FlagKunten["訓点"]; !ok {
		t.Error("should have Kunten flag")
	}
}

func TestLoadChukiTagMap_CommentAndEmpty(t *testing.T) {
	input := "# comment\n\n改丁\t\t\tP\n## also comment\n"
	m, err := LoadChukiTagMap(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(m.Tags) != 1 {
		t.Errorf("want 1 entry, got %d", len(m.Tags))
	}
}

func TestLoadChukiTagMap_SingleColumn(t *testing.T) {
	input := "改丁\n"
	m, err := LoadChukiTagMap(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	tags := m.Tags["改丁"]
	if len(tags) != 1 || tags[0] != "" {
		t.Errorf("single column should have empty tag: got %v", tags)
	}
}

func TestLoadChukiTagMap_RealFile(t *testing.T) {
	f, err := os.Open("../../temp/AozoraEpub3/chuki_tag.txt")
	if err != nil {
		t.Skip("chuki_tag.txt not found:", err)
	}
	defer func() { _ = f.Close() }()

	m, err := LoadChukiTagMap(f)
	if err != nil {
		t.Fatal(err)
	}
	if len(m.Tags) < 100 {
		t.Errorf("expected at least 100 entries, got %d", len(m.Tags))
	}
	// 既知のエントリを検証
	if _, ok := m.Tags["改丁"]; !ok {
		t.Error("改丁 should exist")
	}
	if _, ok := m.FlagPageBreak["改丁"]; !ok {
		t.Error("改丁 should have PageBreak flag")
	}
}

func TestLoadChukiSufMap_Basic(t *testing.T) {
	input := "に傍点\t傍点\t傍点終わり\n"
	m, err := LoadChukiSufMap(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	tags, ok := m.Tags["に傍点"]
	if !ok {
		t.Fatal("に傍点 not found")
	}
	if len(tags) != 2 {
		t.Fatalf("want 2 tags, got %d", len(tags))
	}
	if tags[0] != "傍点" || tags[1] != "傍点終わり" {
		t.Errorf("tags: want [傍点, 傍点終わり], got %v", tags)
	}
}

func TestLoadChukiSufMap_NoEndTag(t *testing.T) {
	input := "は太字\t太字\t\n"
	m, err := LoadChukiSufMap(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	tags := m.Tags["は太字"]
	if len(tags) != 1 {
		t.Fatalf("want 1 tag, got %d", len(tags))
	}
	if tags[0] != "太字" {
		t.Errorf("tag: want 太字, got %s", tags[0])
	}
}

func TestLoadChukiSufMap_Alias(t *testing.T) {
	input := "に傍点\t傍点\t傍点終わり\t「」\n"
	m, err := LoadChukiSufMap(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	// 元のキー
	if _, ok := m.Tags["に傍点"]; !ok {
		t.Error("に傍点 should exist")
	}
	// 別名: 「」+ に傍点
	aliasKey := "「」に傍点"
	tags, ok := m.Tags[aliasKey]
	if !ok {
		t.Fatalf("alias key %s not found", aliasKey)
	}
	if tags[0] != "傍点" {
		t.Errorf("alias tag: want 傍点, got %s", tags[0])
	}
}

func TestLoadChukiSufMap_RealFile(t *testing.T) {
	f, err := os.Open("../../temp/AozoraEpub3/chuki_tag_suf.txt")
	if err != nil {
		t.Skip("chuki_tag_suf.txt not found:", err)
	}
	defer func() { _ = f.Close() }()

	m, err := LoadChukiSufMap(f)
	if err != nil {
		t.Fatal(err)
	}
	if len(m.Tags) < 50 {
		t.Errorf("expected at least 50 entries, got %d", len(m.Tags))
	}
	if _, ok := m.Tags["に傍点"]; !ok {
		t.Error("に傍点 should exist")
	}
}
