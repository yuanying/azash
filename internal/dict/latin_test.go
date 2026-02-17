package dict

import (
	"os"
	"strings"
	"testing"
)

func TestLoadLatinMap_Basic(t *testing.T) {
	input := "A`\tÀ\t164\t8883\n"
	m, err := LoadLatinMap(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	ch, ok := m.DecompToChar["A`"]
	if !ok {
		t.Fatal("A` not found")
	}
	if ch != 'À' {
		t.Errorf("want À, got %c", ch)
	}
	cid, ok := m.CharToCID['À']
	if !ok {
		t.Fatal("À CID not found")
	}
	if len(cid) != 2 || cid[0] != "164" || cid[1] != "8883" {
		t.Errorf("CID: want [164, 8883], got %v", cid)
	}
}

func TestLoadLatinMap_EmptyDecomp(t *testing.T) {
	// 分解表記が空の行は DecompToChar に追加しない
	input := "\tª\t140\t8859\n"
	m, err := LoadLatinMap(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(m.DecompToChar) != 0 {
		t.Errorf("want 0 decomp entries, got %d", len(m.DecompToChar))
	}
	// CID は登録される
	if _, ok := m.CharToCID['ª']; !ok {
		t.Error("ª CID should be registered")
	}
}

func TestLoadLatinMap_CommentAndEmpty(t *testing.T) {
	input := "# comment\n\nA`\tÀ\t164\t8883\n"
	m, err := LoadLatinMap(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(m.DecompToChar) != 1 {
		t.Errorf("want 1 entry, got %d", len(m.DecompToChar))
	}
}

func TestLoadLatinMap_RealFile(t *testing.T) {
	f, err := os.Open("../../temp/AozoraEpub3/chuki_latin.txt")
	if err != nil {
		t.Skip("chuki_latin.txt not found:", err)
	}
	defer func() { _ = f.Close() }()

	m, err := LoadLatinMap(f)
	if err != nil {
		t.Fatal(err)
	}
	if len(m.DecompToChar) < 50 {
		t.Errorf("expected at least 50 decomp entries, got %d", len(m.DecompToChar))
	}
	if len(m.CharToCID) < 50 {
		t.Errorf("expected at least 50 CID entries, got %d", len(m.CharToCID))
	}
}
