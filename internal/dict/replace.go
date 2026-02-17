package dict

import (
	"bufio"
	"io"
	"strings"
	"unicode/utf8"
)

// ReplaceMap は文字置換マップを保持する。
type ReplaceMap struct {
	Single map[rune]string   // 1文字置換
	Double map[string]string // 2文字置換
}

// LoadReplaceMap は replace.txt 形式のタブ区切りファイルを読み込む。
// Java 版は String.length()（UTF-16 コードユニット数）で 1/2 文字を判定するため、
// BMP 外文字はサロゲートペアで 2 ユニットとなり Double に分類される。
func LoadReplaceMap(r io.Reader) (*ReplaceMap, error) {
	m := &ReplaceMap{
		Single: make(map[rune]string),
		Double: make(map[string]string),
	}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || line[0] == '#' {
			continue
		}
		// Java 版は split("\t") で末尾空要素を落とすため、Go でも同様にトリム
		values := splitTabJava(line)
		if len(values) < 2 {
			continue
		}
		from := values[0]
		to := values[1]

		// Java 版は UTF-16 コードユニット長で判定
		switch utf16Units(from) {
		case 1:
			ch, _ := utf8.DecodeRuneInString(from)
			m.Single[ch] = to
		case 2:
			m.Double[from] = to
		}
	}
	return m, scanner.Err()
}

// splitTabJava は Java の String.split("\t") と互換のタブ分割を行う。
// Java の split は末尾の空要素を落とす。
func splitTabJava(s string) []string {
	v := strings.Split(s, "\t")
	for len(v) > 0 && v[len(v)-1] == "" {
		v = v[:len(v)-1]
	}
	return v
}

// utf16Units は文字列の UTF-16 コードユニット数をアロケーションなしで計算する。
func utf16Units(s string) int {
	n := 0
	for _, r := range s {
		if r <= 0xFFFF {
			n++
		} else {
			n += 2
		}
	}
	return n
}
