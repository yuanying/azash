package model

import "testing"

func TestNewBookInfo_Defaults(t *testing.T) {
	b := NewBookInfo("/path/to/file.txt")
	if b.SrcFile != "/path/to/file.txt" {
		t.Errorf("SrcFile = %q, want %q", b.SrcFile, "/path/to/file.txt")
	}
	if b.Modified.IsZero() {
		t.Error("Modified should be set")
	}
	if !b.Vertical {
		t.Error("Vertical should be true")
	}
	if !b.InsertTitleToc {
		t.Error("InsertTitleToc should be true")
	}

	// -1 sentinel values
	sentinelFields := map[string]int{
		"TotalLineNum":        b.TotalLineNum,
		"TitleLine":           b.TitleLine,
		"SubTitleLine":        b.SubTitleLine,
		"OrgTitleLine":        b.OrgTitleLine,
		"SubOrgTitleLine":     b.SubOrgTitleLine,
		"CreatorLine":         b.CreatorLine,
		"SubCreatorLine":      b.SubCreatorLine,
		"SeriesLine":          b.SeriesLine,
		"PublisherLine":       b.PublisherLine,
		"LanguageLine":        b.LanguageLine,
		"TitleEndLine":        b.TitleEndLine,
		"FirstCommentLineNum": b.FirstCommentLineNum,
		"FirstImageLineNum":   b.FirstImageLineNum,
		"FirstImageIdx":       b.FirstImageIdx,
		"CoverImageIndex":     b.CoverImageIndex,
	}
	for name, val := range sentinelFields {
		if val != -1 {
			t.Errorf("%s = %d, want -1", name, val)
		}
	}

	// Zero-value defaults
	if b.RTL {
		t.Error("RTL should be false")
	}
	if b.ImageOnly {
		t.Error("ImageOnly should be false")
	}
	if b.CoverFileName != nil {
		t.Error("CoverFileName should be nil")
	}
	if b.CoverEditInfo != nil {
		t.Error("CoverEditInfo should be nil")
	}
}

func TestBookInfo_ImageSectionLine(t *testing.T) {
	b := NewBookInfo("")

	// nil安全性
	if b.IsImageSectionLine(10) {
		t.Error("should be false for nil map")
	}
	if got, ok := b.GetImageSectionFileName(10); ok || got != "" {
		t.Errorf("should be (\"\", false) for nil map, got (%q, %v)", got, ok)
	}

	b.AddImageSectionLine(10, "image.png")
	if !b.IsImageSectionLine(10) {
		t.Error("should be true after add")
	}
	if got, ok := b.GetImageSectionFileName(10); !ok || got != "image.png" {
		t.Errorf("GetImageSectionFileName(10) = (%q, %v), want (\"image.png\", true)", got, ok)
	}
	if b.IsImageSectionLine(11) {
		t.Error("should be false for non-existent line")
	}
}

func TestBookInfo_PageBreakLine(t *testing.T) {
	b := NewBookInfo("")

	if b.IsPageBreakLine(5) {
		t.Error("should be false for nil map")
	}

	b.AddPageBreakLine(5)
	if !b.IsPageBreakLine(5) {
		t.Error("should be true after add")
	}
	if b.IsPageBreakLine(6) {
		t.Error("should be false for non-existent line")
	}
}

func TestBookInfo_NoPageBreakLine(t *testing.T) {
	b := NewBookInfo("")

	if b.IsNoPageBreakLine(5) {
		t.Error("should be false for nil map")
	}

	b.AddNoPageBreakLine(5)
	if !b.IsNoPageBreakLine(5) {
		t.Error("should be true after add")
	}
}

func TestBookInfo_IgnoreLine(t *testing.T) {
	b := NewBookInfo("")

	if b.IsIgnoreLine(5) {
		t.Error("should be false for nil map")
	}

	b.AddIgnoreLine(5)
	if !b.IsIgnoreLine(5) {
		t.Error("should be true after add")
	}
}

func TestBookInfo_ChapterLineInfo(t *testing.T) {
	b := NewBookInfo("")

	// nil安全性
	if got := b.GetChapterLineInfo(10); got != nil {
		t.Error("should be nil for nil map")
	}
	if got := b.GetChapterLevel(10); got != 0 {
		t.Errorf("GetChapterLevel should be 0 for nil map, got %d", got)
	}

	info := NewChapterLineInfo(10, TypeChukiH1, false, LevelH1, false, "第一章")
	b.AddChapterLineInfo(info)

	if got := b.GetChapterLineInfo(10); got != info {
		t.Error("should return the same info")
	}
	if got := b.GetChapterLevel(10); got != LevelH1 {
		t.Errorf("GetChapterLevel(10) = %d, want %d", got, LevelH1)
	}

	b.RemoveChapterLineInfo(10)
	if got := b.GetChapterLineInfo(10); got != nil {
		t.Error("should be nil after remove")
	}
}

func TestBookInfo_GetChapterLineInfoList_Sorted(t *testing.T) {
	b := NewBookInfo("")

	// nil case
	if got := b.GetChapterLineInfoList(); got != nil {
		t.Error("should be nil for nil map")
	}

	// 逆順に追加して、ソート順で返ることを確認
	b.AddChapterLineInfo(NewChapterLineInfo(30, TypeChukiH1, false, LevelH1, false, "三"))
	b.AddChapterLineInfo(NewChapterLineInfo(10, TypeChukiH1, false, LevelH1, false, "一"))
	b.AddChapterLineInfo(NewChapterLineInfo(20, TypeChukiH1, false, LevelH1, false, "二"))

	list := b.GetChapterLineInfoList()
	if len(list) != 3 {
		t.Fatalf("len = %d, want 3", len(list))
	}
	wantLines := []int{10, 20, 30}
	for i, info := range list {
		if info.LineNum != wantLines[i] {
			t.Errorf("list[%d].LineNum = %d, want %d", i, info.LineNum, wantLines[i])
		}
	}
}

func TestBookInfo_ExcludeTocChapter_NilSafe(t *testing.T) {
	b := NewBookInfo("")
	b.ExcludeTocChapter() // should not panic
}

func TestBookInfo_ExcludeTocChapter_ConsecutivePatterns(t *testing.T) {
	b := NewBookInfo("")

	// 5連続のパターン行 → 全除外
	for i := 10; i <= 14; i++ {
		b.AddChapterLineInfo(NewChapterLineInfo(i, TypePattern, false, LevelH1, false, ""))
	}

	b.ExcludeTocChapter()

	for i := 10; i <= 14; i++ {
		if b.GetChapterLineInfo(i) != nil {
			t.Errorf("line %d should be excluded", i)
		}
	}
}

func TestBookInfo_ExcludeTocChapter_ThreePatterns(t *testing.T) {
	b := NewBookInfo("")

	// 3連続のパターン行 → 全除外
	for i := 10; i <= 12; i++ {
		b.AddChapterLineInfo(NewChapterLineInfo(i, TypePattern, false, LevelH1, false, ""))
	}

	b.ExcludeTocChapter()

	for i := 10; i <= 12; i++ {
		if b.GetChapterLineInfo(i) != nil {
			t.Errorf("line %d should be excluded", i)
		}
	}
}

func TestBookInfo_ExcludeTocChapter_TwoPatterns_NotExcluded(t *testing.T) {
	b := NewBookInfo("")

	// 2つだけのパターン行 → 除外されない
	b.AddChapterLineInfo(NewChapterLineInfo(10, TypePattern, false, LevelH1, false, ""))
	b.AddChapterLineInfo(NewChapterLineInfo(11, TypePattern, false, LevelH1, false, ""))

	b.ExcludeTocChapter()

	if b.GetChapterLineInfo(10) == nil {
		t.Error("line 10 should NOT be excluded (only 2 consecutive)")
	}
	if b.GetChapterLineInfo(11) == nil {
		t.Error("line 11 should NOT be excluded (only 2 consecutive)")
	}
}

func TestBookInfo_ExcludeTocChapter_NonPattern_NotExcluded(t *testing.T) {
	b := NewBookInfo("")

	// 非パターン型は除外されない
	b.AddChapterLineInfo(NewChapterLineInfo(10, TypeChukiH1, false, LevelH1, false, ""))
	b.AddChapterLineInfo(NewChapterLineInfo(11, TypeChukiH1, false, LevelH1, false, ""))
	b.AddChapterLineInfo(NewChapterLineInfo(12, TypeChukiH1, false, LevelH1, false, ""))

	b.ExcludeTocChapter()

	for i := 10; i <= 12; i++ {
		if b.GetChapterLineInfo(i) == nil {
			t.Errorf("line %d should NOT be excluded (non-pattern type)", i)
		}
	}
}

func TestBookInfo_MetaLineText(t *testing.T) {
	b := NewBookInfo("")
	b.MetaLines = []string{"タイトル", "著者名", "シリーズ"}
	b.MetaLineStart = 5
	b.TitleLine = 5
	b.CreatorLine = 6
	b.SeriesLine = 7

	if got := b.TitleText(); got != "タイトル" {
		t.Errorf("TitleText() = %q, want %q", got, "タイトル")
	}
	if got := b.CreatorText(); got != "著者名" {
		t.Errorf("CreatorText() = %q, want %q", got, "著者名")
	}
	if got := b.SeriesText(); got != "シリーズ" {
		t.Errorf("SeriesText() = %q, want %q", got, "シリーズ")
	}
}

func TestBookInfo_MetaLineText_Unset(t *testing.T) {
	b := NewBookInfo("")
	// -1 = 未設定 → 空文字列
	if got := b.TitleText(); got != "" {
		t.Errorf("TitleText() = %q, want empty", got)
	}
	if got := b.SubTitleText(); got != "" {
		t.Errorf("SubTitleText() = %q, want empty", got)
	}
	if got := b.OrgTitleText(); got != "" {
		t.Errorf("OrgTitleText() = %q, want empty", got)
	}
	if got := b.SubOrgTitleText(); got != "" {
		t.Errorf("SubOrgTitleText() = %q, want empty", got)
	}
	if got := b.CreatorText(); got != "" {
		t.Errorf("CreatorText() = %q, want empty", got)
	}
	if got := b.SubCreatorText(); got != "" {
		t.Errorf("SubCreatorText() = %q, want empty", got)
	}
	if got := b.PublisherText(); got != "" {
		t.Errorf("PublisherText() = %q, want empty", got)
	}
	if got := b.LanguageText(); got != "" {
		t.Errorf("LanguageText() = %q, want empty", got)
	}
}

func TestBookInfo_MetaLineText_OutOfBounds(t *testing.T) {
	b := NewBookInfo("")
	b.MetaLines = []string{"line"}
	b.MetaLineStart = 0
	b.TitleLine = 10 // out of bounds

	if got := b.TitleText(); got != "" {
		t.Errorf("TitleText() for out-of-bounds = %q, want empty", got)
	}
}

func TestBookInfo_Clear(t *testing.T) {
	b := NewBookInfo("")
	b.AddImageSectionLine(1, "img.png")
	b.AddPageBreakLine(2)
	b.AddNoPageBreakLine(3)
	b.AddIgnoreLine(4)

	b.Clear()

	if b.IsImageSectionLine(1) {
		t.Error("ImageSectionLine should be cleared")
	}
	if b.IsPageBreakLine(2) {
		t.Error("PageBreakLine should be cleared")
	}
	if b.IsNoPageBreakLine(3) {
		t.Error("NoPageBreakLine should be cleared")
	}
	if b.IsIgnoreLine(4) {
		t.Error("IgnoreLine should be cleared")
	}
}

func TestBookInfo_Clear_NilSafe(t *testing.T) {
	b := NewBookInfo("")
	b.Clear() // should not panic
}

func TestTitleType_HasTitle(t *testing.T) {
	tests := []struct {
		tt   TitleType
		want bool
	}{
		{TitleTypeTitleAuthor, true},
		{TitleTypeAuthorTitle, true},
		{TitleTypeSubTitleAuthor, true},
		{TitleTypeTitleOnly, true},
		{TitleTypeTitleAuthorOnly, true},
		{TitleTypeNone, false},
	}
	for _, tt := range tests {
		if got := tt.tt.HasTitle(); got != tt.want {
			t.Errorf("TitleType(%d).HasTitle() = %v, want %v", tt.tt, got, tt.want)
		}
	}
}

func TestTitleType_HasAuthor(t *testing.T) {
	tests := []struct {
		tt   TitleType
		want bool
	}{
		{TitleTypeTitleAuthor, true},
		{TitleTypeAuthorTitle, true},
		{TitleTypeSubTitleAuthor, true},
		{TitleTypeTitleOnly, false},
		{TitleTypeTitleAuthorOnly, true},
		{TitleTypeNone, false},
	}
	for _, tt := range tests {
		if got := tt.tt.HasAuthor(); got != tt.want {
			t.Errorf("TitleType(%d).HasAuthor() = %v, want %v", tt.tt, got, tt.want)
		}
	}
}

func TestTitleType_HasTitleAuthor(t *testing.T) {
	tests := []struct {
		tt   TitleType
		want bool
	}{
		{TitleTypeTitleAuthor, true},
		{TitleTypeAuthorTitle, true},
		{TitleTypeSubTitleAuthor, true},
		{TitleTypeTitleOnly, false},
		{TitleTypeTitleAuthorOnly, true},
		{TitleTypeNone, false},
	}
	for _, tt := range tests {
		if got := tt.tt.HasTitleAuthor(); got != tt.want {
			t.Errorf("TitleType(%d).HasTitleAuthor() = %v, want %v", tt.tt, got, tt.want)
		}
	}
}

func TestTitleType_TitleFirst(t *testing.T) {
	tests := []struct {
		tt   TitleType
		want bool
	}{
		{TitleTypeTitleAuthor, true},
		{TitleTypeAuthorTitle, false},
		{TitleTypeSubTitleAuthor, true},
		{TitleTypeTitleOnly, true},
		{TitleTypeTitleAuthorOnly, true},
		{TitleTypeNone, false},
	}
	for _, tt := range tests {
		if got := tt.tt.TitleFirst(); got != tt.want {
			t.Errorf("TitleType(%d).TitleFirst() = %v, want %v", tt.tt, got, tt.want)
		}
	}
}

func TestTitleType_String(t *testing.T) {
	tests := []struct {
		tt   TitleType
		want string
	}{
		{TitleTypeTitleAuthor, "表題 → 著者名"},
		{TitleTypeAuthorTitle, "著者名 → 表題"},
		{TitleTypeSubTitleAuthor, "表題 → 著者名(副題優先)"},
		{TitleTypeTitleOnly, "表題のみ(1行)"},
		{TitleTypeTitleAuthorOnly, "表題+著者のみ(2行)"},
		{TitleTypeNone, "なし"},
		{TitleType(-1), ""},
		{TitleType(100), ""},
	}
	for _, tt := range tests {
		if got := tt.tt.String(); got != tt.want {
			t.Errorf("TitleType(%d).String() = %q, want %q", tt.tt, got, tt.want)
		}
	}
}

func TestTitleType_IotaValues(t *testing.T) {
	if TitleTypeTitleAuthor != 0 {
		t.Errorf("TitleTypeTitleAuthor = %d, want 0", TitleTypeTitleAuthor)
	}
	if TitleTypeAuthorTitle != 1 {
		t.Errorf("TitleTypeAuthorTitle = %d, want 1", TitleTypeAuthorTitle)
	}
	if TitleTypeSubTitleAuthor != 2 {
		t.Errorf("TitleTypeSubTitleAuthor = %d, want 2", TitleTypeSubTitleAuthor)
	}
	if TitleTypeTitleOnly != 3 {
		t.Errorf("TitleTypeTitleOnly = %d, want 3", TitleTypeTitleOnly)
	}
	if TitleTypeTitleAuthorOnly != 4 {
		t.Errorf("TitleTypeTitleAuthorOnly = %d, want 4", TitleTypeTitleAuthorOnly)
	}
	if TitleTypeNone != 5 {
		t.Errorf("TitleTypeNone = %d, want 5", TitleTypeNone)
	}
}

func TestTitleTypeIndexOf_Bounds(t *testing.T) {
	// 範囲外は TitleTypeNone にフォールバック
	if got := TitleTypeIndexOf(-1); got != TitleTypeNone {
		t.Errorf("TitleTypeIndexOf(-1) = %d, want %d", got, TitleTypeNone)
	}
	if got := TitleTypeIndexOf(100); got != TitleTypeNone {
		t.Errorf("TitleTypeIndexOf(100) = %d, want %d", got, TitleTypeNone)
	}
	// 正常範囲
	if got := TitleTypeIndexOf(0); got != TitleTypeTitleAuthor {
		t.Errorf("TitleTypeIndexOf(0) = %d, want %d", got, TitleTypeTitleAuthor)
	}
	if got := TitleTypeIndexOf(5); got != TitleTypeNone {
		t.Errorf("TitleTypeIndexOf(5) = %d, want %d", got, TitleTypeNone)
	}
}

func TestTitlePageConstants(t *testing.T) {
	if TitlePageNone != -1 {
		t.Errorf("TitlePageNone = %d, want -1", TitlePageNone)
	}
	if TitlePageNormal != 0 {
		t.Errorf("TitlePageNormal = %d, want 0", TitlePageNormal)
	}
	if TitlePageMiddle != 1 {
		t.Errorf("TitlePageMiddle = %d, want 1", TitlePageMiddle)
	}
	if TitlePageHorizontal != 2 {
		t.Errorf("TitlePageHorizontal = %d, want 2", TitlePageHorizontal)
	}
}
