package dict

import (
	"bufio"
	"io"
	"strings"
	"unicode/utf8"
)

// LatinMap は拡張ラテン文字の変換マップを保持する。
type LatinMap struct {
	DecompToChar map[string]rune   // 分解表記 → 拡張ラテン文字
	CharToCID    map[rune][]string // ラテン文字 → [横CID, 縦CID]
}

// LoadLatinMap は chuki_latin.txt 形式のタブ区切りファイルを読み込む。
func LoadLatinMap(r io.Reader) (*LatinMap, error) {
	m := &LatinMap{
		DecompToChar: make(map[string]rune),
		CharToCID:    make(map[rune][]string),
	}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || line[0] == '#' {
			continue
		}
		values := strings.Split(line, "\t")
		if len(values) < 2 {
			continue
		}

		ch, _ := utf8.DecodeRuneInString(values[1])
		if ch == utf8.RuneError {
			continue
		}

		if values[0] != "" {
			m.DecompToChar[values[0]] = ch
		}

		if len(values) > 3 {
			m.CharToCID[ch] = []string{values[2], values[3]}
		}
	}
	return m, scanner.Err()
}
