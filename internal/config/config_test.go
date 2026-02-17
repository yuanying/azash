package config

import "testing"

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
