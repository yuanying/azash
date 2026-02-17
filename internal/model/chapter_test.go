package model

import "testing"

func TestChapterLineInfo_TypeID(t *testing.T) {
	tests := []struct {
		chType int
		want   string
	}{
		{TypeTitle, "題"},
		{TypePageBreak, ""},
		{TypeChukiH, "見"},
		{TypeChukiH1, "大"},
		{TypeChukiH2, "中"},
		{TypeChukiH3, "小"},
		{TypeChapterName, "章"},
		{TypeChapterNum, "数"},
		{TypePattern, "他"},
		{999, ""},
	}
	for _, tt := range tests {
		c := &ChapterLineInfo{Type: tt.chType}
		if got := c.TypeID(); got != tt.want {
			t.Errorf("TypeID() for type %d = %q, want %q", tt.chType, got, tt.want)
		}
	}
}

func TestChapterLineInfo_IsPattern(t *testing.T) {
	patternTypes := []int{TypeChukiH, TypeChapterName, TypeChapterNum, TypePattern}
	for _, tp := range patternTypes {
		c := &ChapterLineInfo{Type: tp}
		if !c.IsPattern() {
			t.Errorf("IsPattern() for type %d should be true", tp)
		}
	}
	nonPatternTypes := []int{TypeTitle, TypePageBreak, TypeChukiH1, TypeChukiH2, TypeChukiH3}
	for _, tp := range nonPatternTypes {
		c := &ChapterLineInfo{Type: tp}
		if c.IsPattern() {
			t.Errorf("IsPattern() for type %d should be false", tp)
		}
	}
}

func TestChapterLineInfo_JoinChapterName(t *testing.T) {
	t.Run("empty initial", func(t *testing.T) {
		c := &ChapterLineInfo{}
		c.JoinChapterName("第一章")
		if c.ChapterName != "第一章" {
			t.Errorf("ChapterName = %q, want %q", c.ChapterName, "第一章")
		}
	})
	t.Run("join with existing", func(t *testing.T) {
		c := &ChapterLineInfo{ChapterName: "第一章"}
		c.JoinChapterName("序文")
		want := "第一章\u3000序文"
		if c.ChapterName != want {
			t.Errorf("ChapterName = %q, want %q", c.ChapterName, want)
		}
	})
}

func TestChapterTypeFromID(t *testing.T) {
	tests := []struct {
		id   rune
		want int
	}{
		{'題', TypeTitle},
		{'改', TypePageBreak},
		{'見', TypeChukiH},
		{'大', TypeChukiH1},
		{'中', TypeChukiH2},
		{'小', TypeChukiH3},
		{'章', TypeChapterName},
		{'数', TypeChapterNum},
		{'他', TypePattern},
		{'?', 0},
	}
	for _, tt := range tests {
		if got := ChapterTypeFromID(tt.id); got != tt.want {
			t.Errorf("ChapterTypeFromID(%c) = %d, want %d", tt.id, got, tt.want)
		}
	}
}

func TestLevelFromType(t *testing.T) {
	tests := []struct {
		chType int
		want   int
	}{
		{TypeTitle, LevelTitle},
		{TypePageBreak, LevelSection},
		{TypeChukiH1, LevelH1},
		{TypeChukiH2, LevelH2},
		{TypeChukiH3, LevelH3},
		{TypeChapterNum, LevelH2},
		{TypeChukiH, LevelH1},
		{TypeChapterName, LevelH1},
		{TypePattern, LevelH1},
	}
	for _, tt := range tests {
		if got := LevelFromType(tt.chType); got != tt.want {
			t.Errorf("LevelFromType(%d) = %d, want %d", tt.chType, got, tt.want)
		}
	}
}

func TestChapterInfo_NoTagChapterName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"no tags", "第一章", "第一章"},
		{"with tags", "<b>第一章</b>", "第一章"},
		{"nested tags", "<span class=\"x\"><b>序文</b></span>", "序文"},
		{"empty", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ChapterInfo{ChapterName: tt.input}
			if got := c.NoTagChapterName(); got != tt.want {
				t.Errorf("NoTagChapterName() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestNewChapterLineInfo(t *testing.T) {
	c := NewChapterLineInfo(10, TypeChukiH1, true, LevelH1, false, "第一章")
	if c.LineNum != 10 {
		t.Errorf("LineNum = %d, want 10", c.LineNum)
	}
	if c.Type != TypeChukiH1 {
		t.Errorf("Type = %d, want %d", c.Type, TypeChukiH1)
	}
	if !c.PageBreakChapter {
		t.Error("PageBreakChapter should be true")
	}
	if c.Level != LevelH1 {
		t.Errorf("Level = %d, want %d", c.Level, LevelH1)
	}
	if c.EmptyNext {
		t.Error("EmptyNext should be false")
	}
	if c.ChapterName != "第一章" {
		t.Errorf("ChapterName = %q, want %q", c.ChapterName, "第一章")
	}
}

func TestNewChapterInfo(t *testing.T) {
	c := NewChapterInfo("sec001", "ch001", "第一章", 1)
	if c.SectionID != "sec001" {
		t.Errorf("SectionID = %q, want %q", c.SectionID, "sec001")
	}
	if c.ChapterID != "ch001" {
		t.Errorf("ChapterID = %q, want %q", c.ChapterID, "ch001")
	}
	if c.ChapterName != "第一章" {
		t.Errorf("ChapterName = %q, want %q", c.ChapterName, "第一章")
	}
	if c.ChapterLevel != 1 {
		t.Errorf("ChapterLevel = %d, want 1", c.ChapterLevel)
	}
	if c.NavClose != 1 {
		t.Errorf("NavClose = %d, want 1", c.NavClose)
	}
	if c.LevelStart != 0 {
		t.Errorf("LevelStart = %d, want 0", c.LevelStart)
	}
	if c.LevelEnd != 0 {
		t.Errorf("LevelEnd = %d, want 0", c.LevelEnd)
	}
}

func TestSetTocNestLevel_Empty(t *testing.T) {
	SetTocNestLevel(false, false, nil, false)
	SetTocNestLevel(false, false, []*ChapterInfo{}, false)
}

func TestSetTocNestLevel_Flat(t *testing.T) {
	chapters := []*ChapterInfo{
		{ChapterLevel: 1},
		{ChapterLevel: 1},
		{ChapterLevel: 1},
	}
	SetTocNestLevel(false, false, chapters, false)
	for i, ch := range chapters {
		if ch.ChapterLevel != 0 {
			t.Errorf("chapters[%d].ChapterLevel = %d, want 0", i, ch.ChapterLevel)
		}
	}
	// 全て同一レベルなので LevelStart/LevelEnd は 0
	for i, ch := range chapters {
		if ch.LevelStart != 0 {
			t.Errorf("chapters[%d].LevelStart = %d, want 0", i, ch.LevelStart)
		}
		if ch.LevelEnd != 0 {
			t.Errorf("chapters[%d].LevelEnd = %d, want 0", i, ch.LevelEnd)
		}
	}
}

func TestSetTocNestLevel_Nested(t *testing.T) {
	// Input levels: 1, 1, 2, 2, 1
	// After ancestor counting: 0, 0, 1, 1, 0
	chapters := []*ChapterInfo{
		{ChapterLevel: 1},
		{ChapterLevel: 1},
		{ChapterLevel: 2},
		{ChapterLevel: 2},
		{ChapterLevel: 1},
	}
	SetTocNestLevel(false, false, chapters, false)

	wantLevels := []int{0, 0, 1, 1, 0}
	for i, ch := range chapters {
		if ch.ChapterLevel != wantLevels[i] {
			t.Errorf("chapters[%d].ChapterLevel = %d, want %d", i, ch.ChapterLevel, wantLevels[i])
		}
	}
	// levelStart
	wantLevelStart := []int{0, 0, 1, 0, 0}
	for i, ch := range chapters {
		if ch.LevelStart != wantLevelStart[i] {
			t.Errorf("chapters[%d].LevelStart = %d, want %d", i, ch.LevelStart, wantLevelStart[i])
		}
	}
	// levelEnd
	wantLevelEnd := []int{0, 0, 0, 1, 0}
	for i, ch := range chapters {
		if ch.LevelEnd != wantLevelEnd[i] {
			t.Errorf("chapters[%d].LevelEnd = %d, want %d", i, ch.LevelEnd, wantLevelEnd[i])
		}
	}
}

func TestSetTocNestLevel_InsertTitleToc(t *testing.T) {
	chapters := []*ChapterInfo{
		{ChapterLevel: 0},
		{ChapterLevel: 1},
		{ChapterLevel: 1},
	}
	SetTocNestLevel(false, false, chapters, true)
	// insertTitleToc: chapters[0].ChapterLevel should be set to chapters[1].ChapterLevel before processing
	// Original levels become [1, 1, 1] then all become 0
	for i, ch := range chapters {
		if ch.ChapterLevel != 0 {
			t.Errorf("chapters[%d].ChapterLevel = %d, want 0", i, ch.ChapterLevel)
		}
	}
}

func TestSetTocNestLevel_NcxNest(t *testing.T) {
	// Input levels: 1, 1, 2, 2, 1
	// After processing: levels = 0, 0, 1, 1, 0
	chapters := []*ChapterInfo{
		{ChapterLevel: 1},
		{ChapterLevel: 1},
		{ChapterLevel: 2},
		{ChapterLevel: 2},
		{ChapterLevel: 1},
	}
	SetTocNestLevel(false, true, chapters, false)

	// navClose: levelEnd + 1, with parent nodes set to 0
	wantNavClose := []int{1, 0, 1, 2, 1}
	for i, ch := range chapters {
		if ch.NavClose != wantNavClose[i] {
			t.Errorf("chapters[%d].NavClose = %d, want %d", i, ch.NavClose, wantNavClose[i])
		}
	}
}

func TestChapterLineInfo_Constants(t *testing.T) {
	if TypeTitle != 1 {
		t.Errorf("TypeTitle = %d, want 1", TypeTitle)
	}
	if TypePageBreak != 2 {
		t.Errorf("TypePageBreak = %d, want 2", TypePageBreak)
	}
	if TypeChukiH != 10 {
		t.Errorf("TypeChukiH = %d, want 10", TypeChukiH)
	}
	if TypeChukiH1 != 11 {
		t.Errorf("TypeChukiH1 = %d, want 11", TypeChukiH1)
	}
	if TypeChukiH2 != 12 {
		t.Errorf("TypeChukiH2 = %d, want 12", TypeChukiH2)
	}
	if TypeChukiH3 != 13 {
		t.Errorf("TypeChukiH3 = %d, want 13", TypeChukiH3)
	}
	if TypeChapterName != 21 {
		t.Errorf("TypeChapterName = %d, want 21", TypeChapterName)
	}
	if TypeChapterNum != 22 {
		t.Errorf("TypeChapterNum = %d, want 22", TypeChapterNum)
	}
	if TypePattern != 30 {
		t.Errorf("TypePattern = %d, want 30", TypePattern)
	}
	if LevelTitle != 0 {
		t.Errorf("LevelTitle = %d, want 0", LevelTitle)
	}
	if LevelSection != 1 {
		t.Errorf("LevelSection = %d, want 1", LevelSection)
	}
	if LevelH1 != 1 {
		t.Errorf("LevelH1 = %d, want 1", LevelH1)
	}
	if LevelH2 != 2 {
		t.Errorf("LevelH2 = %d, want 2", LevelH2)
	}
	if LevelH3 != 3 {
		t.Errorf("LevelH3 = %d, want 3", LevelH3)
	}
}
