package dict

import (
	"bufio"
	"io"
	"strings"
)

// ChukiTagMap は注記タグ変換マップを保持する。
type ChukiTagMap struct {
	Tags            map[string][]string // 注記キー → [tag, endTag]
	FlagNoBr        map[string]struct{} // フラグ '1'
	FlagNoRubyStart map[string]struct{} // フラグ '2'
	FlagNoRubyEnd   map[string]struct{} // フラグ '3'
	FlagPageBreak   map[string]struct{} // フラグ 'P','M','L'
	FlagMiddle      map[string]struct{} // フラグ 'M'
	FlagBottom      map[string]struct{} // フラグ 'L'
	FlagKunten      map[string]struct{} // フラグ 'K'
}

// LoadChukiTagMap は chuki_tag.txt 形式のタブ区切りファイルを読み込む。
func LoadChukiTagMap(r io.Reader) (*ChukiTagMap, error) {
	m := &ChukiTagMap{
		Tags:            make(map[string][]string),
		FlagNoBr:        make(map[string]struct{}),
		FlagNoRubyStart: make(map[string]struct{}),
		FlagNoRubyEnd:   make(map[string]struct{}),
		FlagPageBreak:   make(map[string]struct{}),
		FlagMiddle:      make(map[string]struct{}),
		FlagBottom:      make(map[string]struct{}),
		FlagKunten:      make(map[string]struct{}),
	}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || line[0] == '#' {
			continue
		}
		values := strings.Split(line, "\t")
		key := values[0]

		var tags []string
		switch {
		case len(values) == 1:
			tags = []string{""}
		case len(values) > 2 && values[2] != "":
			tags = []string{values[1], values[2]}
		default:
			tags = []string{values[1]}
		}
		m.Tags[key] = tags

		if len(values) > 3 && values[3] != "" {
			switch values[3][0] {
			case '1':
				m.FlagNoBr[key] = struct{}{}
			case '2':
				m.FlagNoRubyStart[key] = struct{}{}
			case '3':
				m.FlagNoRubyEnd[key] = struct{}{}
			case 'P':
				m.FlagPageBreak[key] = struct{}{}
			case 'M':
				m.FlagPageBreak[key] = struct{}{}
				m.FlagMiddle[key] = struct{}{}
			case 'L':
				m.FlagPageBreak[key] = struct{}{}
				m.FlagBottom[key] = struct{}{}
			case 'K':
				m.FlagKunten[key] = struct{}{}
			}
		}
	}
	return m, scanner.Err()
}

// ChukiSufMap は注記サフィックス変換マップを保持する。
type ChukiSufMap struct {
	Tags map[string][]string // key → [startTag, endTag]
}

// LoadChukiSufMap は chuki_tag_suf.txt 形式のタブ区切りファイルを読み込む。
func LoadChukiSufMap(r io.Reader) (*ChukiSufMap, error) {
	m := &ChukiSufMap{
		Tags: make(map[string][]string),
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
		key := values[0]

		var tags []string
		if len(values) > 2 && values[2] != "" {
			tags = []string{values[1], values[2]}
		} else {
			tags = []string{values[1]}
		}
		m.Tags[key] = tags

		// 別名キー
		if len(values) > 3 && values[3] != "" {
			m.Tags[values[3]+key] = tags
		}
	}
	return m, scanner.Err()
}
