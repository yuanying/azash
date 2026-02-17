package model

import "regexp"

// 見出し種別定数
const (
	TypeTitle       = 1
	TypePageBreak   = 2
	TypeChukiH      = 10
	TypeChukiH1     = 11
	TypeChukiH2     = 12
	TypeChukiH3     = 13
	TypeChapterName = 21
	TypeChapterNum  = 22
	TypePattern     = 30
)

// 見出しレベル定数
const (
	LevelTitle   = 0
	LevelSection = 1
	LevelH1      = 1
	LevelH2      = 2
	LevelH3      = 3
)

// ChapterLineInfo は見出し行の情報を格納する。
type ChapterLineInfo struct {
	LineNum          int
	Type             int
	Level            int
	PageBreakChapter bool
	EmptyNext        bool
	ChapterName      string
}

// NewChapterLineInfo は見出し行情報を生成する。
func NewChapterLineInfo(lineNum, chapterType int, pageBreak bool, level int, emptyNext bool, chapterName string) *ChapterLineInfo {
	return &ChapterLineInfo{
		LineNum:          lineNum,
		Type:             chapterType,
		PageBreakChapter: pageBreak,
		Level:            level,
		EmptyNext:        emptyNext,
		ChapterName:      chapterName,
	}
}

// TypeID は見出し種別の表示用文字列を返す。
func (c *ChapterLineInfo) TypeID() string {
	switch c.Type {
	case TypeTitle:
		return "題"
	case TypePageBreak:
		return ""
	case TypeChukiH:
		return "見"
	case TypeChukiH1:
		return "大"
	case TypeChukiH2:
		return "中"
	case TypeChukiH3:
		return "小"
	case TypeChapterName:
		return "章"
	case TypeChapterNum:
		return "数"
	case TypePattern:
		return "他"
	}
	return ""
}

// IsPattern は章名や数字やパターンでマッチした行なら true を返す。
func (c *ChapterLineInfo) IsPattern() bool {
	switch c.Type {
	case TypeTitle, TypePageBreak, TypeChukiH1, TypeChukiH2, TypeChukiH3:
		return false
	}
	return true
}

// JoinChapterName は既存の章名に新しい名前を結合する。
func (c *ChapterLineInfo) JoinChapterName(name string) {
	if c.ChapterName == "" {
		c.ChapterName = name
	} else {
		c.ChapterName = c.ChapterName + "\u3000" + name
	}
}

// ChapterTypeFromID は表示用文字から見出し種別定数を返す。
func ChapterTypeFromID(typeID rune) int {
	switch typeID {
	case '題':
		return TypeTitle
	case '改':
		return TypePageBreak
	case '見':
		return TypeChukiH
	case '大':
		return TypeChukiH1
	case '中':
		return TypeChukiH2
	case '小':
		return TypeChukiH3
	case '章':
		return TypeChapterName
	case '数':
		return TypeChapterNum
	case '他':
		return TypePattern
	}
	return 0
}

// LevelFromType は見出し種別から見出しレベルを返す。
func LevelFromType(chapterType int) int {
	switch chapterType {
	case TypeTitle:
		return LevelTitle
	case TypePageBreak:
		return LevelSection
	case TypeChukiH1:
		return LevelH1
	case TypeChukiH2:
		return LevelH2
	case TypeChukiH3:
		return LevelH3
	case TypeChapterNum:
		return LevelH2
	}
	return LevelH1
}

// ChapterInfo は目次用の章の情報を格納する。
type ChapterInfo struct {
	SectionID    string
	ChapterID    string
	ChapterName  string
	ChapterLevel int
	LevelStart   int
	LevelEnd     int
	NavClose     int
}

// NewChapterInfo は章情報を生成する。
func NewChapterInfo(sectionID, chapterID, chapterName string, chapterLevel int) *ChapterInfo {
	return &ChapterInfo{
		SectionID:    sectionID,
		ChapterID:    chapterID,
		ChapterName:  chapterName,
		ChapterLevel: chapterLevel,
		NavClose:     1,
	}
}

var tagRegexp = regexp.MustCompile(`<[^>]*>`)

// NoTagChapterName は HTML タグを除去した章名を返す。
func (c *ChapterInfo) NoTagChapterName() string {
	return tagRegexp.ReplaceAllString(c.ChapterName, "")
}

// SetTocNestLevel は目次の階層情報を設定する。
// navNest は Java 版 API 互換のために保持しているが、現在は未使用。
func SetTocNestLevel(_ /* navNest */, ncxNest bool, chapters []*ChapterInfo, insertTitleToc bool) {
	if len(chapters) == 0 {
		return
	}

	// 表題のレベルを2つめと同じにする
	if insertTitleToc && len(chapters) >= 2 {
		chapters[0].ChapterLevel = chapters[1].ChapterLevel
	}

	// 見出しレベルをコピーする
	levels := make([]int, len(chapters))
	for i, ch := range chapters {
		levels[i] = ch.ChapterLevel
	}

	// 目次の階層レベルの設定
	for i := range chapters {
		count := 0
		currLevel := levels[i]
		for j := i - 1; j >= 0; j-- {
			if levels[j] < currLevel {
				count++
				currLevel = levels[j]
			}
		}
		chapters[i].ChapterLevel = count
	}

	// Velocity用のフィールド設定
	curr := chapters[0]
	curr.LevelStart = 0

	prev := curr
	for i := 1; i < len(chapters); i++ {
		curr = chapters[i]
		diff := curr.ChapterLevel - prev.ChapterLevel
		if diff > 0 {
			curr.LevelStart = diff
			prev.LevelEnd = 0
		} else {
			curr.LevelStart = 0
			prev.LevelEnd = -diff
		}
		prev = curr
	}
	curr.LevelEnd = curr.ChapterLevel

	if ncxNest {
		prev = nil
		for _, ch := range chapters {
			curr = ch
			curr.NavClose = curr.LevelEnd + 1
			if curr.LevelStart > 0 && prev != nil {
				prev.NavClose = 0
			}
			prev = curr
		}
	}
}
