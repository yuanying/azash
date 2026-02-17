package dict

import (
	"bufio"
	"io"
	"strings"
)

// GaijiMap は外字注記の変換マップを保持する。
type GaijiMap struct {
	UtfMap map[string]string // 注記キー → UTF-8 文字
	AltMap map[string]string // 注記キー → 代替文字列
}

// LoadGaijiUtfMap は chuki_utf.txt / chuki_ivs.txt 形式のファイルを読み込む。
// existing が nil でない場合、既存のマップにマージし、重複キーはスキップする。
// これにより IVS → UTF の優先度制御が可能。
func LoadGaijiUtfMap(r io.Reader, existing *GaijiMap) (*GaijiMap, error) {
	m := existing
	if m == nil {
		m = &GaijiMap{}
	}
	if m.UtfMap == nil {
		m.UtfMap = make(map[string]string)
	}
	if m.AltMap == nil {
		m.AltMap = make(map[string]string)
	}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || line[0] == '#' {
			continue
		}
		utfChar, chukiKey, ok := parseGaijiLine(line)
		if !ok {
			continue
		}
		if _, exists := m.UtfMap[chukiKey]; !exists {
			m.UtfMap[chukiKey] = utfChar
		}
	}
	return m, scanner.Err()
}

// LoadGaijiAltMap は chuki_alt.txt 形式のファイルを読み込む。
func LoadGaijiAltMap(r io.Reader) (*GaijiMap, error) {
	m := &GaijiMap{
		UtfMap: make(map[string]string),
		AltMap: make(map[string]string),
	}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || line[0] == '#' {
			continue
		}
		utfChar, chukiKey, ok := parseGaijiLine(line)
		if !ok {
			continue
		}
		// Java の loadChukiFile は重複キーを先勝ち（containsKey チェック）
		if _, exists := m.AltMap[chukiKey]; !exists {
			m.AltMap[chukiKey] = utfChar
		}
	}
	return m, scanner.Err()
}

// parseGaijiLine は外字行をパースし、UTF-8文字と注記キーを返す。
// Java 版 AozoraGaijiConverter.loadChukiFile 互換。
func parseGaijiLine(line string) (utfChar, chukiKey string, ok bool) {
	// タブ位置を順に探す: col1 \t col2 \t UTF文字 \t 注記 [\t ...]
	tab1 := strings.IndexByte(line, '\t')
	if tab1 == -1 {
		return "", "", false
	}
	tab2 := strings.IndexByte(line[tab1+1:], '\t')
	if tab2 == -1 {
		return "", "", false
	}
	tab2 += tab1 + 1
	tab3 := strings.IndexByte(line[tab2+1:], '\t')
	if tab3 == -1 {
		return "", "", false
	}
	tab3 += tab2 + 1

	utfChar = line[tab2+1 : tab3]

	// 注記部分を抽出（4列目以降のタブがあればそこまで）
	annotation := line[tab3+1:]
	if nextTab := strings.IndexByte(annotation, '\t'); nextTab != -1 {
		annotation = annotation[:nextTab]
	}

	// 注記は「※［＃」で始まる必要がある
	if !strings.HasPrefix(annotation, "※［＃") {
		return "", "", false
	}

	// ※［＃ の後のコンテンツ
	content := annotation[len("※［＃"):]

	// キー終了位置を探す: 「、」または「］」
	commaIdx := strings.Index(content, "、")
	// 「、」の直後が「」の場合は次の「、」を探す（Java版互換）
	if commaIdx != -1 {
		after := content[commaIdx+len("、"):]
		if strings.HasPrefix(after, "「") {
			nextComma := strings.Index(after, "、")
			if nextComma != -1 {
				commaIdx = commaIdx + len("、") + nextComma
			} else {
				commaIdx = -1
			}
		}
	}

	if commaIdx != -1 {
		chukiKey = content[:commaIdx]
	} else {
		bracketIdx := strings.Index(content, "］")
		if bracketIdx != -1 {
			chukiKey = content[:bracketIdx]
		} else {
			chukiKey = content
		}
	}

	ok = true
	return
}
