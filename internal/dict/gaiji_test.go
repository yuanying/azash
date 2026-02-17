package dict

import (
	"os"
	"strings"
	"testing"
)

func TestLoadGaijiMap_BasicUTF(t *testing.T) {
	input := "U+04E02\t\t丂\t※［＃「朽のつくり」、第4水準2-1-2］\n"
	m, err := LoadGaijiUtfMap(strings.NewReader(input), nil)
	if err != nil {
		t.Fatal(err)
	}
	v, ok := m.UtfMap["「朽のつくり」"]
	if !ok {
		t.Fatal("key not found")
	}
	if v != "丂" {
		t.Errorf("want 丂, got %s", v)
	}
}

func TestLoadGaijiMap_KeyWithCommaParam(t *testing.T) {
	// 「、」の後に「」が続く場合は次の「、」を探す
	input := "U+04E02\t\t丂\t※［＃「朽のつくり」、第4水準2-1-2］\n"
	m, err := LoadGaijiUtfMap(strings.NewReader(input), nil)
	if err != nil {
		t.Fatal(err)
	}
	// キーは「※［＃」の後から「、」の手前まで
	if _, ok := m.UtfMap["「朽のつくり」"]; !ok {
		t.Error("key should be extracted up to 、")
	}
}

func TestLoadGaijiMap_KeyWithoutComma(t *testing.T) {
	input := "U+04E04\t\t丄\t※［＃「ぼう／一」］\n"
	m, err := LoadGaijiUtfMap(strings.NewReader(input), nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := m.UtfMap["「ぼう／一」"]; !ok {
		t.Error("key should be extracted up to ］")
	}
}

func TestLoadGaijiMap_CommentAndEmpty(t *testing.T) {
	input := "# comment\n\nU+04E02\t\t丂\t※［＃「朽のつくり」、第4水準2-1-2］\n"
	m, err := LoadGaijiUtfMap(strings.NewReader(input), nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(m.UtfMap) != 1 {
		t.Errorf("want 1 entry, got %d", len(m.UtfMap))
	}
}

func TestLoadGaijiMap_NonGaijiLineSkipped(t *testing.T) {
	// 注記が ※［＃ で始まらない行はスキップ
	input := "U+04E02\t\t丂\tsome other text\n"
	m, err := LoadGaijiUtfMap(strings.NewReader(input), nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(m.UtfMap) != 0 {
		t.Errorf("want 0 entries, got %d", len(m.UtfMap))
	}
}

func TestLoadGaijiMap_IVSPriority(t *testing.T) {
	ivsInput := "U+04E19\tU+E0101\t丙󠄁\t※［＃「丙」の「人」に代えて「入」］\t→コメント\n"
	utfInput := "U+04E19\t\t丙\t※［＃「丙」の「人」に代えて「入」］\n"

	// IVS を先にロード
	m, err := LoadGaijiUtfMap(strings.NewReader(ivsInput), nil)
	if err != nil {
		t.Fatal(err)
	}
	// UTF をロード（既存キーはスキップ）
	m, err = LoadGaijiUtfMap(strings.NewReader(utfInput), m)
	if err != nil {
		t.Fatal(err)
	}

	v := m.UtfMap["「丙」の「人」に代えて「入」"]
	if v != "丙󠄁" {
		t.Errorf("IVS should take priority: want 丙󠄁, got %s", v)
	}
}

func TestLoadGaijiMap_ExistingNilMaps(t *testing.T) {
	// existing が非 nil だが内部 map が nil の場合に panic しないことを確認
	input := "U+04E02\t\t丂\t※［＃「朽のつくり」、第4水準2-1-2］\n"
	existing := &GaijiMap{}
	m, err := LoadGaijiUtfMap(strings.NewReader(input), existing)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := m.UtfMap["「朽のつくり」"]; !ok {
		t.Error("key should exist")
	}
}

func TestLoadGaijiAltMap(t *testing.T) {
	input := "\t\t［＃縦中横］!!!［＃縦中横終わり］\t※［＃感嘆符三つ］\n"
	m, err := LoadGaijiAltMap(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	v, ok := m.AltMap["感嘆符三つ"]
	if !ok {
		t.Fatal("key not found")
	}
	if v != "［＃縦中横］!!!［＃縦中横終わり］" {
		t.Errorf("want ［＃縦中横］!!!［＃縦中横終わり］, got %s", v)
	}
}

func TestLoadGaijiAltMap_DuplicateKeyFirstWins(t *testing.T) {
	// Java の loadChukiFile は重複キーを先勝ち（containsKey チェック）
	input := "\t\tFirst\t※［＃同じキー］\n\t\tSecond\t※［＃同じキー］\n"
	m, err := LoadGaijiAltMap(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	v := m.AltMap["同じキー"]
	if v != "First" {
		t.Errorf("duplicate key should be first-wins: want First, got %s", v)
	}
}

func TestLoadGaijiMap_RealUTFFile(t *testing.T) {
	f, err := os.Open("../../temp/AozoraEpub3/chuki_utf.txt")
	if err != nil {
		t.Skip("chuki_utf.txt not found:", err)
	}
	defer func() { _ = f.Close() }()

	m, err := LoadGaijiUtfMap(f, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(m.UtfMap) < 6000 {
		t.Errorf("expected at least 6000 entries, got %d", len(m.UtfMap))
	}
}

func TestLoadGaijiMap_RealIVSFile(t *testing.T) {
	ivsF, err := os.Open("../../temp/AozoraEpub3/chuki_ivs.txt")
	if err != nil {
		t.Skip("chuki_ivs.txt not found:", err)
	}
	defer func() { _ = ivsF.Close() }()

	m, err := LoadGaijiUtfMap(ivsF, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(m.UtfMap) < 100 {
		t.Errorf("expected at least 100 IVS entries, got %d", len(m.UtfMap))
	}
}

func TestLoadGaijiMap_RealAltFile(t *testing.T) {
	f, err := os.Open("../../temp/AozoraEpub3/chuki_alt.txt")
	if err != nil {
		t.Skip("chuki_alt.txt not found:", err)
	}
	defer func() { _ = f.Close() }()

	m, err := LoadGaijiAltMap(f)
	if err != nil {
		t.Fatal(err)
	}
	if len(m.AltMap) < 15 {
		t.Errorf("expected at least 15 entries, got %d", len(m.AltMap))
	}
}
