package config

import (
	"os"
	"strings"
	"testing"
)

func TestNewConfig_Defaults(t *testing.T) {
	c := NewConfig()

	// 改ページ
	assertBool(t, "PageBreak", false, c.PageBreak)
	assertInt(t, "PageBreakSize", 0, c.PageBreakSize)
	assertInt(t, "PageBreakEmptyLine", 0, c.PageBreakEmptyLine)
	assertBool(t, "PageBreakEmpty", false, c.PageBreakEmpty)
	assertInt(t, "PageBreakEmptySize", 0, c.PageBreakEmptySize)
	assertBool(t, "PageBreakChapter", false, c.PageBreakChapter)
	assertInt(t, "PageBreakChapterSize", 0, c.PageBreakChapterSize)

	// 表示サイズ
	assertInt(t, "DispW", 600, c.DispW)
	assertInt(t, "DispH", 800, c.DispH)

	// 表紙
	assertBool(t, "CoverPage", false, c.CoverPage)
	assertInt(t, "CoverW", 600, c.CoverW)
	assertInt(t, "CoverH", 800, c.CoverH)

	// 画像
	assertInt(t, "ImageSizeType", 2, c.ImageSizeType)
	assertBool(t, "FitImage", false, c.FitImage)
	assertBool(t, "SvgImage", false, c.SvgImage)
	assertInt(t, "JpegQuality", 80, c.JpegQuality)
	assertFloat(t, "ImageScale", 1.0, c.ImageScale)
	assertInt(t, "RotateImage", 0, c.RotateImage)
	assertInt(t, "ImageFloatType", 0, c.ImageFloatType)
	assertInt(t, "ImageFloatW", 0, c.ImageFloatW)
	assertInt(t, "ImageFloatH", 0, c.ImageFloatH)

	// リサイズ
	assertBool(t, "ResizeW", false, c.ResizeW)
	assertBool(t, "ResizeH", false, c.ResizeH)
	assertInt(t, "ResizeNumW", 0, c.ResizeNumW)
	assertInt(t, "ResizeNumH", 0, c.ResizeNumH)
	assertInt(t, "SinglePageSizeW", 480, c.SinglePageSizeW)
	assertInt(t, "SinglePageSizeH", 640, c.SinglePageSizeH)
	assertInt(t, "SinglePageWidth", 600, c.SinglePageWidth)

	// 自動余白
	assertBool(t, "AutoMargin", false, c.AutoMargin)
	assertInt(t, "AutoMarginLimitH", 0, c.AutoMarginLimitH)
	assertInt(t, "AutoMarginLimitV", 0, c.AutoMarginLimitV)
	assertInt(t, "AutoMarginWhiteLevel", 80, c.AutoMarginWhiteLevel)
	assertFloat(t, "AutoMarginPadding", 0, c.AutoMarginPadding)
	assertInt(t, "AutoMarginNombre", 0, c.AutoMarginNombre)
	assertFloat(t, "AutoMarginNombreSize", 0.03, c.AutoMarginNombreSize)

	// 目次/タイトル
	assertInt(t, "TitlePage", -1, c.TitlePage)
	assertBool(t, "TitlePageWrite", false, c.TitlePageWrite)
	assertBool(t, "TocPage", false, c.TocPage)
	assertBool(t, "TocVertical", false, c.TocVertical)
	assertBool(t, "CoverPageToc", false, c.CoverPageToc)
	assertBool(t, "NavNest", false, c.NavNest)
	assertBool(t, "NcxNest", false, c.NcxNest)
	assertBool(t, "TitleToc", false, c.TitleToc)

	// 章設定
	assertBool(t, "ChapterSection", true, c.ChapterSection)
	assertBool(t, "ChapterH", false, c.ChapterH)
	assertBool(t, "ChapterH1", false, c.ChapterH1)
	assertBool(t, "ChapterH2", false, c.ChapterH2)
	assertBool(t, "ChapterH3", false, c.ChapterH3)
	assertInt(t, "ChapterNameLength", 64, c.ChapterNameLength)
	assertBool(t, "ChapterExclude", false, c.ChapterExclude)
	assertBool(t, "ChapterUseNextLine", false, c.ChapterUseNextLine)
	assertBool(t, "SameLineChapter", false, c.SameLineChapter)
	assertBool(t, "ChapterName", false, c.ChapterName)
	assertBool(t, "ChapterNumOnly", false, c.ChapterNumOnly)
	assertBool(t, "ChapterNumTitle", false, c.ChapterNumTitle)
	assertBool(t, "ChapterNumParen", false, c.ChapterNumParen)
	assertBool(t, "ChapterNumParenTitle", false, c.ChapterNumParenTitle)
	assertBool(t, "ChapterPattern", false, c.ChapterPattern)
	assertString(t, "ChapterPatternText", "", c.ChapterPatternText)

	// スタイル
	assertFloat(t, "LineHeight", 1.8, c.LineHeight)
	assertInt(t, "FontSize", 100, c.FontSize)
	assertString(t, "PageMargin", "0,0,0,0", c.PageMargin)
	assertString(t, "PageMarginUnit", "em", c.PageMarginUnit)
	assertString(t, "BodyMargin", "0,0,0,0", c.BodyMargin)
	assertString(t, "BodyMarginUnit", "em", c.BodyMarginUnit)
	assertBool(t, "BoldUseGothic", false, c.BoldUseGothic)
	assertBool(t, "GothicUseBold", false, c.GothicUseBold)

	// テキスト
	assertInt(t, "DakutenType", 0, c.DakutenType)
	assertInt(t, "SpaceHyphenation", 0, c.SpaceHyphenation)
	assertBool(t, "AutoYoko", false, c.AutoYoko)
	assertBool(t, "AutoYokoNum1", false, c.AutoYokoNum1)
	assertBool(t, "AutoYokoNum3", false, c.AutoYokoNum3)
	assertBool(t, "AutoYokoEQ1", false, c.AutoYokoEQ1)
	assertInt(t, "RemoveEmptyLine", 0, c.RemoveEmptyLine)
	assertInt(t, "MaxEmptyLine", 0, c.MaxEmptyLine)

	// その他
	assertBool(t, "MarkId", false, c.MarkId)
	assertBool(t, "CommentPrint", false, c.CommentPrint)
	assertBool(t, "CommentConvert", false, c.CommentConvert)
	assertBool(t, "NoIllust", false, c.NoIllust)
	assertBool(t, "IvsBMP", false, c.IvsBMP)
	assertBool(t, "IvsSSP", false, c.IvsSSP)

	// ガンマ
	assertFloat(t, "Gamma", 1.0, c.Gamma)
}

func TestLoadFromReader_BasicKeyValue(t *testing.T) {
	input := "DispW=800\nDispH=1024\n"
	c, warnings, err := LoadFromReader(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(warnings) != 0 {
		t.Errorf("unexpected warnings: %v", warnings)
	}
	assertInt(t, "DispW", 800, c.DispW)
	assertInt(t, "DispH", 1024, c.DispH)
}

func TestLoadFromReader_CommentAndEmptyLine(t *testing.T) {
	input := "# comment\n\nDispW=800\n"
	c, _, err := LoadFromReader(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	assertInt(t, "DispW", 800, c.DispW)
	// 他はデフォルト
	assertInt(t, "DispH", 800, c.DispH)
}

func TestLoadFromReader_BoolConversion(t *testing.T) {
	input := "CoverPage=1\nFitImage=0\nSvgImage=\n"
	c, _, err := LoadFromReader(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	assertBool(t, "CoverPage", true, c.CoverPage)
	assertBool(t, "FitImage", false, c.FitImage)
	assertBool(t, "SvgImage", false, c.SvgImage)
}

func TestLoadFromReader_IntParseWithDefault(t *testing.T) {
	input := "DispW=abc\n"
	c, _, err := LoadFromReader(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	// パース失敗時はデフォルト値
	assertInt(t, "DispW", 600, c.DispW)
}

func TestLoadFromReader_FloatParse(t *testing.T) {
	input := "LineHeight=2.5\nImageScale=0.5\n"
	c, _, err := LoadFromReader(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	assertFloat(t, "LineHeight", 2.5, c.LineHeight)
	assertFloat(t, "ImageScale", 0.5, c.ImageScale)
}

func TestLoadFromReader_Margin(t *testing.T) {
	input := "PageMargin=1,2,3,4\nPageMarginUnit=0\nBodyMargin=5,6,7,8\nBodyMarginUnit=1\n"
	c, _, err := LoadFromReader(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	assertString(t, "PageMargin", "1,2,3,4", c.PageMargin)
	assertString(t, "PageMarginUnit", "em", c.PageMarginUnit)
	assertString(t, "BodyMargin", "5,6,7,8", c.BodyMargin)
	assertString(t, "BodyMarginUnit", "%", c.BodyMarginUnit)
}

func TestLoadFromReader_MarginUnitAbsent(t *testing.T) {
	// Java 版: PageMarginUnit が未指定の場合、Objects.equals(null, "0") → false → "%"
	input := "PageMargin=1,2,3,4\nBodyMargin=5,6,7,8\n"
	c, _, err := LoadFromReader(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	assertString(t, "PageMarginUnit", "%", c.PageMarginUnit)
	assertString(t, "BodyMarginUnit", "%", c.BodyMarginUnit)
}

func TestLoadFromReader_MarginNon4CSV(t *testing.T) {
	// Java 版: PageMargin が4要素に分割できない場合、unit を適用せずデフォルト "0,0,0,0" にリセット
	input := "PageMargin=1,2,3\nPageMarginUnit=0\nBodyMargin=5\nBodyMarginUnit=0\n"
	c, _, err := LoadFromReader(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	// 4要素でないのでデフォルトにリセットされ、unit はデフォルトの "em"
	assertString(t, "PageMargin", "0,0,0,0", c.PageMargin)
	assertString(t, "PageMarginUnit", "em", c.PageMarginUnit)
	assertString(t, "BodyMargin", "0,0,0,0", c.BodyMargin)
	assertString(t, "BodyMarginUnit", "em", c.BodyMarginUnit)
}

func TestLoadFromReader_ExclamationComment(t *testing.T) {
	input := "! comment\nDispW=800\n"
	c, _, err := LoadFromReader(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	assertInt(t, "DispW", 800, c.DispW)
}

func TestLoadFromReader_LeadingWhitespace(t *testing.T) {
	input := "  DispW=800\n\tDispH=600\n"
	c, _, err := LoadFromReader(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	assertInt(t, "DispW", 800, c.DispW)
	assertInt(t, "DispH", 600, c.DispH)
}

func TestLoadFromReader_WhitespaceAroundEquals(t *testing.T) {
	// Java Properties: key と value の前後空白をトリム
	input := "DispW = 800\n DispH =600\n"
	c, _, err := LoadFromReader(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	assertInt(t, "DispW", 800, c.DispW)
	assertInt(t, "DispH", 600, c.DispH)
}

func TestLoadFromReader_UnknownKey(t *testing.T) {
	input := "UnknownKey=value\nDispW=800\n"
	c, warnings, err := LoadFromReader(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(warnings) != 1 {
		t.Fatalf("want 1 warning, got %d", len(warnings))
	}
	if !strings.Contains(warnings[0], "UnknownKey") {
		t.Errorf("warning should mention UnknownKey: %s", warnings[0])
	}
	assertInt(t, "DispW", 800, c.DispW)
}

func TestLoadFromReader_ChapterSectionAbsent(t *testing.T) {
	input := "DispW=800\n"
	c, _, err := LoadFromReader(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	// ChapterSection はキー不在時 true（Java 版 L236 互換）
	assertBool(t, "ChapterSection", true, c.ChapterSection)
}

func TestLoadFromReader_ChapterSectionPresent(t *testing.T) {
	input := "ChapterSection=0\n"
	c, _, err := LoadFromReader(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	assertBool(t, "ChapterSection", false, c.ChapterSection)
}

func TestLoadFromReader_PageBreakSizeKBtoBytes(t *testing.T) {
	input := "PageBreak=1\nPageBreakSize=300\nPageBreakEmpty=1\nPageBreakEmptySize=250\nPageBreakChapter=1\nPageBreakChapterSize=200\n"
	c, _, err := LoadFromReader(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	assertBool(t, "PageBreak", true, c.PageBreak)
	assertInt(t, "PageBreakSize", 300*1024, c.PageBreakSize)
	assertInt(t, "PageBreakEmptySize", 250*1024, c.PageBreakEmptySize)
	assertInt(t, "PageBreakChapterSize", 200*1024, c.PageBreakChapterSize)
}

func TestLoadFromReader_PageBreakDisabled(t *testing.T) {
	// PageBreak=0 のとき、サブプロパティはパースされない
	input := "PageBreakSize=300\nPageBreakEmptySize=250\n"
	c, _, err := LoadFromReader(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	assertInt(t, "PageBreakSize", 0, c.PageBreakSize)
	assertInt(t, "PageBreakEmptySize", 0, c.PageBreakEmptySize)
}

func TestLoadFromReader_RotateImage(t *testing.T) {
	input := "RotateImage=1\n"
	c, _, err := LoadFromReader(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	assertInt(t, "RotateImage", 90, c.RotateImage)

	input = "RotateImage=2\n"
	c, _, err = LoadFromReader(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	assertInt(t, "RotateImage", -90, c.RotateImage)
}

func TestLoadFromReader_TitlePage(t *testing.T) {
	// TitlePageWrite=1 のとき TitlePage がパースされる
	input := "TitlePageWrite=1\nTitlePage=2\n"
	c, _, err := LoadFromReader(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	assertInt(t, "TitlePage", 2, c.TitlePage)

	// TitlePageWrite=0 のとき TitlePage はデフォルト
	input = "TitlePage=2\n"
	c, _, err = LoadFromReader(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	assertInt(t, "TitlePage", -1, c.TitlePage)
}

func TestLoadFromReader_TypoKey(t *testing.T) {
	// Java 版のタイポ hapterNumParenTitle を互換キーとして認識
	input := "hapterNumParenTitle=1\n"
	c, _, err := LoadFromReader(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	assertBool(t, "ChapterNumParenTitle", true, c.ChapterNumParenTitle)
}

func TestLoadFromReader_AutoMarginConditional(t *testing.T) {
	// AutoMargin=0 のとき、サブプロパティはパースされない
	input := "AutoMarginLimitH=100\nAutoMarginWhiteLevel=50\n"
	c, _, err := LoadFromReader(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	assertInt(t, "AutoMarginLimitH", 0, c.AutoMarginLimitH)
	assertInt(t, "AutoMarginWhiteLevel", 80, c.AutoMarginWhiteLevel)

	// AutoMargin=1 のとき
	input = "AutoMargin=1\nAutoMarginLimitH=100\nAutoMarginWhiteLevel=50\n"
	c, _, err = LoadFromReader(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	assertInt(t, "AutoMarginLimitH", 100, c.AutoMarginLimitH)
	assertInt(t, "AutoMarginWhiteLevel", 50, c.AutoMarginWhiteLevel)
}

func TestLoadFromReader_GammaConditional(t *testing.T) {
	// Gamma=1 のとき GammaValue がパースされる
	input := "Gamma=1\nGammaValue=2.2\n"
	c, _, err := LoadFromReader(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	assertFloat(t, "Gamma", 2.2, c.Gamma)

	// Gamma=0 のとき GammaValue はパースされない
	input = "GammaValue=2.2\n"
	c, _, err = LoadFromReader(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	assertFloat(t, "Gamma", 1.0, c.Gamma)
}

func TestLoadFromReader_ChapterPatternConditional(t *testing.T) {
	// ChapterPattern=1 のとき ChapterPatternText がパースされる
	input := "ChapterPattern=1\nChapterPatternText=第.+章\n"
	c, _, err := LoadFromReader(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	assertBool(t, "ChapterPattern", true, c.ChapterPattern)
	assertString(t, "ChapterPatternText", "第.+章", c.ChapterPatternText)

	// ChapterPattern=0 のとき
	input = "ChapterPatternText=第.+章\n"
	c, _, err = LoadFromReader(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	assertString(t, "ChapterPatternText", "", c.ChapterPatternText)
}

func TestLoadFromReader_ResizeConditional(t *testing.T) {
	// ResizeW=1 のとき ResizeNumW がパースされる
	input := "ResizeW=1\nResizeNumW=2048\n"
	c, _, err := LoadFromReader(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	assertInt(t, "ResizeNumW", 2048, c.ResizeNumW)

	// ResizeW=0 のとき
	input = "ResizeNumW=2048\n"
	c, _, err = LoadFromReader(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	assertInt(t, "ResizeNumW", 0, c.ResizeNumW)
}

func TestLoadFromReader_RealINIFile(t *testing.T) {
	f, err := os.Open("../../temp/AozoraEpub3/presets/reader.ini")
	if err != nil {
		t.Skip("reader.ini not found:", err)
	}
	defer func() { _ = f.Close() }()

	c, warnings, err := LoadFromReader(f)
	if err != nil {
		t.Fatal(err)
	}

	// reader.ini の既知値を検証
	assertInt(t, "DispW", 584, c.DispW)
	assertInt(t, "DispH", 754, c.DispH)
	assertInt(t, "CoverW", 600, c.CoverW)
	assertInt(t, "CoverH", 800, c.CoverH)
	assertBool(t, "PageBreak", true, c.PageBreak)
	assertInt(t, "PageBreakSize", 300*1024, c.PageBreakSize)
	assertInt(t, "ImageSizeType", 3, c.ImageSizeType)
	assertBool(t, "FitImage", true, c.FitImage)
	assertBool(t, "SvgImage", true, c.SvgImage)
	assertBool(t, "ResizeW", true, c.ResizeW)
	assertBool(t, "ResizeH", true, c.ResizeH)
	assertInt(t, "ResizeNumW", 2048, c.ResizeNumW)
	assertInt(t, "ResizeNumH", 2048, c.ResizeNumH)
	assertInt(t, "SpaceHyphenation", 2, c.SpaceHyphenation)
	assertFloat(t, "LineHeight", 1.8, c.LineHeight)
	assertInt(t, "FontSize", 100, c.FontSize)
	assertBool(t, "GothicUseBold", true, c.GothicUseBold)

	// reader.ini のキーはすべて既知キーとして処理される
	for _, w := range warnings {
		t.Logf("warning: %s", w)
	}
}

func assertBool(t *testing.T, name string, want, got bool) {
	t.Helper()
	if want != got {
		t.Errorf("%s: want %v, got %v", name, want, got)
	}
}

func assertInt(t *testing.T, name string, want, got int) {
	t.Helper()
	if want != got {
		t.Errorf("%s: want %d, got %d", name, want, got)
	}
}

func assertFloat(t *testing.T, name string, want, got float64) {
	t.Helper()
	if want != got {
		t.Errorf("%s: want %f, got %f", name, want, got)
	}
}

func assertString(t *testing.T, name string, want, got string) {
	t.Helper()
	if want != got {
		t.Errorf("%s: want %q, got %q", name, want, got)
	}
}
