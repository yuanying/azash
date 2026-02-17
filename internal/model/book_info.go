package model

import (
	"sort"
	"time"
)

// TitleType はタイトル記載種別を表す。
type TitleType int

const (
	TitleTypeTitleAuthor     TitleType = iota // 表題 → 著者名
	TitleTypeAuthorTitle                      // 著者名 → 表題
	TitleTypeSubTitleAuthor                   // 表題 → 著者名(副題優先)
	TitleTypeTitleOnly                        // 表題のみ(1行)
	TitleTypeTitleAuthorOnly                  // 表題+著者のみ(2行)
	TitleTypeNone                             // なし
)

var titleTypeNames = [...]string{
	"表題 → 著者名",
	"著者名 → 表題",
	"表題 → 著者名(副題優先)",
	"表題のみ(1行)",
	"表題+著者のみ(2行)",
	"なし",
}

// String は TitleType の表示名を返す。
func (t TitleType) String() string {
	if int(t) < 0 || int(t) >= len(titleTypeNames) {
		return ""
	}
	return titleTypeNames[t]
}

// HasTitle はタイトルを持つ種別なら true を返す。
func (t TitleType) HasTitle() bool {
	return t != TitleTypeNone
}

// HasAuthor は著者を持つ種別なら true を返す。
func (t TitleType) HasAuthor() bool {
	switch t {
	case TitleTypeTitleOnly, TitleTypeNone:
		return false
	default:
		return true
	}
}

// HasTitleAuthor はタイトルと著者の両方を持つ種別なら true を返す。
func (t TitleType) HasTitleAuthor() bool {
	switch t {
	case TitleTypeTitleOnly, TitleTypeNone:
		return false
	default:
		return true
	}
}

// TitleFirst はタイトルが著者より先に記載される種別なら true を返す。
func (t TitleType) TitleFirst() bool {
	switch t {
	case TitleTypeTitleAuthor, TitleTypeSubTitleAuthor, TitleTypeTitleOnly, TitleTypeTitleAuthorOnly:
		return true
	default:
		return false
	}
}

// TitleTypeIndexOf はインデックスから TitleType を返す。
// 範囲外の場合は TitleTypeNone を返す。
func TitleTypeIndexOf(idx int) TitleType {
	if idx < 0 || idx >= len(titleTypeNames) {
		return TitleTypeNone
	}
	return TitleType(idx)
}

// 表題ページ種別定数
const (
	TitlePageNone       = -1 // 表題は出力しない
	TitlePageNormal     = 0  // 表題はそのまま出力
	TitlePageMiddle     = 1  // 表題ページを左右中央
	TitlePageHorizontal = 2  // 表題ページ横書き
)

// BookInfo はタイトル著作者等のメタ情報を格納する。
type BookInfo struct {
	TitlePageType int

	// メタ行
	MetaLines     []string
	MetaLineStart int

	TotalLineNum int

	// タイトル
	Title     string
	TitleLine int
	TitleAs   string

	// 副題
	SubTitle     string
	SubTitleLine int

	// 原題
	OrgTitleLine    int
	SubOrgTitleLine int

	// 著作者
	Creator        string
	CreatorLine    int
	SubCreatorLine int
	CreatorAs      string

	// シリーズ・出版・言語
	SeriesLine    int
	PublisherLine int
	Publisher     string
	LanguageLine  int
	Language      string

	TitleEndLine int

	FirstCommentLineNum int

	Published time.Time
	Modified  time.Time

	Vertical bool
	RTL      bool

	SrcFile       string
	TextEntryName string

	FirstImageLineNum int
	FirstImageIdx     int

	CoverEditInfo   *CoverEditInfo
	CoverFileName   *string // nil=表紙無し, ""=先頭の挿絵
	SVGCoverImage   bool
	CoverImageIndex int
	CoverExt        string

	InsertCoverPage    bool
	InsertCoverPageToc bool
	InsertTitlePage    bool
	InsertTocPage      bool
	TocVertical        bool
	InsertTitleToc     bool
	ImageOnly          bool

	// 行マップ (lazy init)
	mapImageSectionLine map[int]string
	mapPageBreakLine    map[int]struct{}
	mapNoPageBreakLine  map[int]struct{}
	mapIgnoreLine       map[int]struct{}
	mapChapterLine      map[int]*ChapterLineInfo
}

// NewBookInfo は BookInfo を生成する。
func NewBookInfo(srcFile string) *BookInfo {
	return &BookInfo{
		SrcFile:             srcFile,
		Modified:            time.Now(),
		Vertical:            true,
		InsertTitleToc:      true,
		TotalLineNum:        -1,
		TitleLine:           -1,
		SubTitleLine:        -1,
		OrgTitleLine:        -1,
		SubOrgTitleLine:     -1,
		CreatorLine:         -1,
		SubCreatorLine:      -1,
		SeriesLine:          -1,
		PublisherLine:       -1,
		LanguageLine:        -1,
		TitleEndLine:        -1,
		FirstCommentLineNum: -1,
		FirstImageLineNum:   -1,
		FirstImageIdx:       -1,
		CoverImageIndex:     -1,
	}
}

// Clear はすべての行マップをクリアする。
func (b *BookInfo) Clear() {
	clear(b.mapImageSectionLine)
	clear(b.mapPageBreakLine)
	clear(b.mapNoPageBreakLine)
	clear(b.mapIgnoreLine)
}

// ──────────────────────────────────────────
// 画像単体ページ行マップ
// ──────────────────────────────────────────

// AddImageSectionLine は画像単体ページの行番号とファイル名を保存する。
func (b *BookInfo) AddImageSectionLine(lineNum int, imageFileName string) {
	if b.mapImageSectionLine == nil {
		b.mapImageSectionLine = make(map[int]string)
	}
	b.mapImageSectionLine[lineNum] = imageFileName
}

// IsImageSectionLine は画像単体ページの行なら true を返す。
func (b *BookInfo) IsImageSectionLine(lineNum int) bool {
	if b.mapImageSectionLine == nil {
		return false
	}
	_, ok := b.mapImageSectionLine[lineNum]
	return ok
}

// GetImageSectionFileName は画像単体ページの画像ファイル名を返す。
// 未登録の場合は ("", false) を返す。
func (b *BookInfo) GetImageSectionFileName(lineNum int) (string, bool) {
	if b.mapImageSectionLine == nil {
		return "", false
	}
	v, ok := b.mapImageSectionLine[lineNum]
	return v, ok
}

// ──────────────────────────────────────────
// 強制改ページ行マップ
// ──────────────────────────────────────────

// AddPageBreakLine は強制改ページ行番号を保存する。
func (b *BookInfo) AddPageBreakLine(lineNum int) {
	if b.mapPageBreakLine == nil {
		b.mapPageBreakLine = make(map[int]struct{})
	}
	b.mapPageBreakLine[lineNum] = struct{}{}
}

// IsPageBreakLine は強制改ページ行なら true を返す。
func (b *BookInfo) IsPageBreakLine(lineNum int) bool {
	if b.mapPageBreakLine == nil {
		return false
	}
	_, ok := b.mapPageBreakLine[lineNum]
	return ok
}

// ──────────────────────────────────────────
// 改ページしない行マップ
// ──────────────────────────────────────────

// AddNoPageBreakLine は改ページしない行番号を保存する。
func (b *BookInfo) AddNoPageBreakLine(lineNum int) {
	if b.mapNoPageBreakLine == nil {
		b.mapNoPageBreakLine = make(map[int]struct{})
	}
	b.mapNoPageBreakLine[lineNum] = struct{}{}
}

// IsNoPageBreakLine は改ページしない行なら true を返す。
func (b *BookInfo) IsNoPageBreakLine(lineNum int) bool {
	if b.mapNoPageBreakLine == nil {
		return false
	}
	_, ok := b.mapNoPageBreakLine[lineNum]
	return ok
}

// ──────────────────────────────────────────
// 出力しない行マップ
// ──────────────────────────────────────────

// AddIgnoreLine は出力しない行番号を保存する。
func (b *BookInfo) AddIgnoreLine(lineNum int) {
	if b.mapIgnoreLine == nil {
		b.mapIgnoreLine = make(map[int]struct{})
	}
	b.mapIgnoreLine[lineNum] = struct{}{}
}

// IsIgnoreLine は出力しない行なら true を返す。
func (b *BookInfo) IsIgnoreLine(lineNum int) bool {
	if b.mapIgnoreLine == nil {
		return false
	}
	_, ok := b.mapIgnoreLine[lineNum]
	return ok
}

// ──────────────────────────────────────────
// 見出し行マップ
// ──────────────────────────────────────────

// AddChapterLineInfo は見出し行情報を保存する。
func (b *BookInfo) AddChapterLineInfo(info *ChapterLineInfo) {
	if b.mapChapterLine == nil {
		b.mapChapterLine = make(map[int]*ChapterLineInfo)
	}
	b.mapChapterLine[info.LineNum] = info
}

// RemoveChapterLineInfo は見出し行情報を削除する。
func (b *BookInfo) RemoveChapterLineInfo(lineNum int) {
	if b.mapChapterLine != nil {
		delete(b.mapChapterLine, lineNum)
	}
}

// GetChapterLineInfo は見出し行の情報を返す。
func (b *BookInfo) GetChapterLineInfo(lineNum int) *ChapterLineInfo {
	if b.mapChapterLine == nil {
		return nil
	}
	return b.mapChapterLine[lineNum]
}

// GetChapterLevel は見出し行なら目次階層レベルを返す。
func (b *BookInfo) GetChapterLevel(lineNum int) int {
	if b.mapChapterLine == nil {
		return 0
	}
	info := b.mapChapterLine[lineNum]
	if info == nil {
		return 0
	}
	return info.Level
}

// GetChapterLineInfoList は行番号順に並び替えた目次一覧リストを返す。
func (b *BookInfo) GetChapterLineInfoList() []*ChapterLineInfo {
	if b.mapChapterLine == nil {
		return nil
	}
	lines := make([]int, 0, len(b.mapChapterLine))
	for lineNum := range b.mapChapterLine {
		lines = append(lines, lineNum)
	}
	sort.Ints(lines)
	list := make([]*ChapterLineInfo, 0, len(lines))
	for _, lineNum := range lines {
		list = append(list, b.mapChapterLine[lineNum])
	}
	return list
}

// ExcludeTocChapter は目次ページの自動抽出見出しを除外する。
func (b *BookInfo) ExcludeTocChapter() {
	if b.mapChapterLine == nil {
		return
	}

	// 前2行と後ろ2行が自動抽出見出しの行を抽出
	excludeLines := make(map[int]struct{})
	for lineNum := range b.mapChapterLine {
		if b.isPattern(lineNum) {
			prevIsPattern := false
			if b.isPattern(lineNum - 1) {
				prevIsPattern = true
			} else if b.mapChapterLine[lineNum].EmptyNext && b.isPattern(lineNum-2) {
				prevIsPattern = true
			}
			nextIsPattern := false
			if b.isPattern(lineNum + 1) {
				nextIsPattern = true
			} else if b.isPattern(lineNum + 2) {
				nextIsPattern = true
			}
			if prevIsPattern && nextIsPattern {
				excludeLines[lineNum] = struct{}{}
			}
		}
	}

	// 先頭と最後
	excludeLines2 := make(map[int]struct{})
	for lineNum := range b.mapChapterLine {
		if _, excluded := excludeLines[lineNum]; !excluded && b.isPattern(lineNum) {
			_, prevExcluded := excludeLines[lineNum-1]
			_, prev2Excluded := excludeLines[lineNum-2]
			_, nextExcluded := excludeLines[lineNum+1]
			_, next2Excluded := excludeLines[lineNum+2]

			switch {
			case prevExcluded:
				excludeLines2[lineNum] = struct{}{}
			case b.mapChapterLine[lineNum].EmptyNext && prev2Excluded:
				excludeLines2[lineNum] = struct{}{}
			case nextExcluded:
				excludeLines2[lineNum] = struct{}{}
			case next2Excluded:
				excludeLines2[lineNum] = struct{}{}
			}
		}
	}

	for lineNum := range excludeLines {
		delete(b.mapChapterLine, lineNum)
	}
	for lineNum := range excludeLines2 {
		delete(b.mapChapterLine, lineNum)
	}
}

func (b *BookInfo) isPattern(lineNum int) bool {
	info := b.GetChapterLineInfo(lineNum)
	if info == nil {
		return false
	}
	return info.IsPattern()
}

// ──────────────────────────────────────────
// メタテキスト取得ヘルパー
// ──────────────────────────────────────────

func (b *BookInfo) metaLineText(line int) string {
	if line == -1 {
		return ""
	}
	idx := line - b.MetaLineStart
	if idx < 0 || idx >= len(b.MetaLines) {
		return ""
	}
	return b.MetaLines[idx]
}

// TitleText はタイトル行のテキストを返す。
func (b *BookInfo) TitleText() string { return b.metaLineText(b.TitleLine) }

// SubTitleText は副題行のテキストを返す。
func (b *BookInfo) SubTitleText() string { return b.metaLineText(b.SubTitleLine) }

// OrgTitleText は原題行のテキストを返す。
func (b *BookInfo) OrgTitleText() string { return b.metaLineText(b.OrgTitleLine) }

// SubOrgTitleText は原副題行のテキストを返す。
func (b *BookInfo) SubOrgTitleText() string { return b.metaLineText(b.SubOrgTitleLine) }

// CreatorText は著作者行のテキストを返す。
func (b *BookInfo) CreatorText() string { return b.metaLineText(b.CreatorLine) }

// SubCreatorText は複著作者行のテキストを返す。
func (b *BookInfo) SubCreatorText() string { return b.metaLineText(b.SubCreatorLine) }

// SeriesText はシリーズ行のテキストを返す。
func (b *BookInfo) SeriesText() string { return b.metaLineText(b.SeriesLine) }

// PublisherText は刊行者行のテキストを返す。
func (b *BookInfo) PublisherText() string { return b.metaLineText(b.PublisherLine) }

// LanguageText は言語行のテキストを返す。
func (b *BookInfo) LanguageText() string { return b.metaLineText(b.LanguageLine) }
