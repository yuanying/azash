package config

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Config は AozoraEpub3 の変換設定を保持する。
// Java 版 AozoraEpub3.java main() の全 getProperty() 呼び出しと互換のフィールドを持つ。
type Config struct {
	// 改ページ
	PageBreak            bool
	PageBreakSize        int // バイト単位（INI は KB）
	PageBreakEmpty       bool
	PageBreakEmptyLine   int
	PageBreakEmptySize   int // バイト単位（INI は KB）
	PageBreakChapter     bool
	PageBreakChapterSize int // バイト単位（INI は KB）

	// 表示サイズ
	DispW int
	DispH int

	// 表紙
	CoverPage bool
	CoverW    int
	CoverH    int

	// 画像
	ImageSizeType  int
	FitImage       bool
	SvgImage       bool
	JpegQuality    int
	ImageScale     float64
	RotateImage    int
	ImageFloatType int
	ImageFloatW    int
	ImageFloatH    int

	// リサイズ
	ResizeW         bool
	ResizeH         bool
	ResizeNumW      int
	ResizeNumH      int
	SinglePageSizeW int
	SinglePageSizeH int
	SinglePageWidth int

	// 自動余白
	AutoMargin           bool
	AutoMarginLimitH     int
	AutoMarginLimitV     int
	AutoMarginWhiteLevel int
	AutoMarginPadding    float64
	AutoMarginNombre     int
	AutoMarginNombreSize float64

	// 目次/タイトル
	TitlePage      int
	TitlePageWrite bool
	TocPage        bool
	TocVertical    bool
	CoverPageToc   bool
	NavNest        bool
	NcxNest        bool
	TitleToc       bool

	// 章設定
	ChapterSection       bool
	ChapterH             bool
	ChapterH1            bool
	ChapterH2            bool
	ChapterH3            bool
	ChapterNameLength    int
	ChapterExclude       bool
	ChapterUseNextLine   bool
	SameLineChapter      bool
	ChapterName          bool
	ChapterNumOnly       bool
	ChapterNumTitle      bool
	ChapterNumParen      bool
	ChapterNumParenTitle bool
	ChapterPattern       bool
	ChapterPatternText   string

	// スタイル
	LineHeight     float64
	FontSize       int
	PageMargin     string
	PageMarginUnit string
	BodyMargin     string
	BodyMarginUnit string
	BoldUseGothic  bool
	GothicUseBold  bool

	// テキスト
	DakutenType      int
	SpaceHyphenation int
	AutoYoko         bool
	AutoYokoNum1     bool
	AutoYokoNum3     bool
	AutoYokoEQ1      bool
	RemoveEmptyLine  int
	MaxEmptyLine     int

	// その他
	MarkId         bool
	CommentPrint   bool
	CommentConvert bool
	NoIllust       bool
	IvsBMP         bool
	IvsSSP         bool

	// ガンマ
	Gamma float64
}

// NewConfig は Java 版と互換のデフォルト値を持つ Config を生成する。
func NewConfig() *Config {
	return &Config{
		// 表示サイズ
		DispW: 600,
		DispH: 800,

		// 表紙
		CoverW: 600,
		CoverH: 800,

		// 画像
		ImageSizeType: 2,
		JpegQuality:   80,
		ImageScale:    1.0,

		// リサイズ
		SinglePageSizeW: 480,
		SinglePageSizeH: 640,
		SinglePageWidth: 600,

		// 自動余白
		AutoMarginWhiteLevel: 80,
		AutoMarginNombreSize: 0.03,

		// 目次/タイトル
		TitlePage: -1,

		// 章設定
		ChapterSection:    true,
		ChapterNameLength: 64,

		// スタイル
		LineHeight:     1.8,
		FontSize:       100,
		PageMargin:     "0,0,0,0",
		PageMarginUnit: "em",
		BodyMargin:     "0,0,0,0",
		BodyMarginUnit: "em",

		// ガンマ
		Gamma: 1.0,
	}
}

// LoadFromReader は Java Properties 形式の INI ファイルを読み込み Config を構築する。
// 未知キーは警告リストに追加される。
func LoadFromReader(r io.Reader) (*Config, []string, error) {
	props, err := parseProperties(r)
	if err != nil {
		return nil, nil, err
	}

	c := NewConfig()
	var warnings []string

	// 既知キーのセット（処理済みフラグ）
	known := make(map[string]bool)
	markKnown := func(keys ...string) {
		for _, k := range keys {
			known[k] = true
		}
	}

	getBool := func(key string) bool {
		markKnown(key)
		return props[key] == "1"
	}
	getInt := func(key string, def int) int {
		markKnown(key)
		v, ok := props[key]
		if !ok || v == "" {
			return def
		}
		n, err := strconv.Atoi(v)
		if err != nil {
			return def
		}
		return n
	}
	getFloat := func(key string, def float64) float64 {
		markKnown(key)
		v, ok := props[key]
		if !ok || v == "" {
			return def
		}
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return def
		}
		return f
	}
	getString := func(key string, def string) string {
		markKnown(key)
		v, ok := props[key]
		if !ok {
			return def
		}
		return v
	}

	// 表紙
	c.CoverPage = getBool("CoverPage")

	// タイトル
	c.TitlePageWrite = getBool("TitlePageWrite")
	if c.TitlePageWrite {
		c.TitlePage = getInt("TitlePage", c.TitlePage)
	} else {
		markKnown("TitlePage")
	}

	c.MarkId = getBool("MarkId")
	c.CommentPrint = getBool("CommentPrint")
	c.CommentConvert = getBool("CommentConvert")
	c.AutoYoko = getBool("AutoYoko")
	c.AutoYokoNum1 = getBool("AutoYokoNum1")
	c.AutoYokoNum3 = getBool("AutoYokoNum3")
	c.AutoYokoEQ1 = getBool("AutoYokoEQ1")
	c.SpaceHyphenation = getInt("SpaceHyphenation", c.SpaceHyphenation)
	c.TocPage = getBool("TocPage")
	c.TocVertical = getBool("TocVertical")
	c.CoverPageToc = getBool("CoverPageToc")
	c.RemoveEmptyLine = getInt("RemoveEmptyLine", c.RemoveEmptyLine)
	c.MaxEmptyLine = getInt("MaxEmptyLine", c.MaxEmptyLine)

	// 画面サイズと画像リサイズ
	c.DispW = getInt("DispW", c.DispW)
	c.DispH = getInt("DispH", c.DispH)
	c.CoverW = getInt("CoverW", c.CoverW)
	c.CoverH = getInt("CoverH", c.CoverH)

	c.ResizeW = getBool("ResizeW")
	if c.ResizeW {
		c.ResizeNumW = getInt("ResizeNumW", c.ResizeNumW)
	} else {
		markKnown("ResizeNumW")
	}
	c.ResizeH = getBool("ResizeH")
	if c.ResizeH {
		c.ResizeNumH = getInt("ResizeNumH", c.ResizeNumH)
	} else {
		markKnown("ResizeNumH")
	}

	c.SinglePageSizeW = getInt("SinglePageSizeW", c.SinglePageSizeW)
	c.SinglePageSizeH = getInt("SinglePageSizeH", c.SinglePageSizeH)
	c.SinglePageWidth = getInt("SinglePageWidth", c.SinglePageWidth)
	c.ImageScale = getFloat("ImageScale", c.ImageScale)
	c.ImageFloatType = getInt("ImageFloatType", c.ImageFloatType)
	c.ImageFloatW = getInt("ImageFloatW", c.ImageFloatW)
	c.ImageFloatH = getInt("ImageFloatH", c.ImageFloatH)
	c.ImageSizeType = getInt("ImageSizeType", c.ImageSizeType)
	c.FitImage = getBool("FitImage")
	c.SvgImage = getBool("SvgImage")

	// RotateImage: "1"→90, "2"→-90, else→0
	markKnown("RotateImage")
	switch props["RotateImage"] {
	case "1":
		c.RotateImage = 90
	case "2":
		c.RotateImage = -90
	}

	c.JpegQuality = getInt("JpegQuality", c.JpegQuality)

	// Gamma: Gamma="1" のとき GammaValue をパース
	markKnown("Gamma", "GammaValue")
	if props["Gamma"] == "1" {
		c.Gamma = getFloat("GammaValue", c.Gamma)
	}

	// 自動余白
	c.AutoMargin = getBool("AutoMargin")
	if c.AutoMargin {
		c.AutoMarginLimitH = getInt("AutoMarginLimitH", c.AutoMarginLimitH)
		c.AutoMarginLimitV = getInt("AutoMarginLimitV", c.AutoMarginLimitV)
		c.AutoMarginWhiteLevel = getInt("AutoMarginWhiteLevel", c.AutoMarginWhiteLevel)
		c.AutoMarginPadding = getFloat("AutoMarginPadding", c.AutoMarginPadding)
		c.AutoMarginNombre = getInt("AutoMarginNombre", c.AutoMarginNombre)
		// Java版 L183 はタイポで AutoMarginNombreSize を autoMarginPadding に代入し、
		// nobreSize は常に 0.03f のまま。Go 版は意図的にこのバグを再現せず、
		// 各フィールドに正しく値を格納する。
		c.AutoMarginNombreSize = getFloat("AutoMarginNombreSize", c.AutoMarginNombreSize)
	} else {
		markKnown("AutoMarginLimitH", "AutoMarginLimitV", "AutoMarginWhiteLevel",
			"AutoMarginPadding", "AutoMarginNombre", "AutoMarginNombreSize")
	}

	// 目次階層化
	c.NavNest = getBool("NavNest")
	c.NcxNest = getBool("NcxNest")

	// スタイル
	// Java版: PageMargin が4要素に分割できた場合にのみ unit を適用。
	// 4要素でない場合はデフォルト "0,0,0,0" にリセットし、unit もデフォルトのまま。
	// unit は Objects.equals(props.getProperty("PageMarginUnit"), "0") で判定。
	// getProperty は key 不在時 null を返すため、"0" 以外はすべて "%" になる。
	c.PageMargin = getString("PageMargin", c.PageMargin)
	markKnown("PageMarginUnit")
	if len(strings.Split(c.PageMargin, ",")) == 4 {
		if props["PageMarginUnit"] == "0" {
			c.PageMarginUnit = "em"
		} else if _, ok := props["PageMargin"]; ok {
			c.PageMarginUnit = "%"
		}
	} else {
		c.PageMargin = "0,0,0,0"
	}
	c.BodyMargin = getString("BodyMargin", c.BodyMargin)
	markKnown("BodyMarginUnit")
	if len(strings.Split(c.BodyMargin, ",")) == 4 {
		if props["BodyMarginUnit"] == "0" {
			c.BodyMarginUnit = "em"
		} else if _, ok := props["BodyMargin"]; ok {
			c.BodyMarginUnit = "%"
		}
	} else {
		c.BodyMargin = "0,0,0,0"
	}
	c.LineHeight = getFloat("LineHeight", c.LineHeight)
	c.FontSize = getInt("FontSize", c.FontSize)
	c.BoldUseGothic = getBool("BoldUseGothic")
	c.GothicUseBold = getBool("gothicUseBold")
	// GothicUseBold の大文字始まりも互換キーとして認識
	if getBool("GothicUseBold") {
		c.GothicUseBold = true
	}

	// 自動改ページ
	c.PageBreak = getBool("PageBreak")
	if c.PageBreak {
		c.PageBreakSize = getInt("PageBreakSize", 0) * 1024
		c.PageBreakEmpty = getBool("PageBreakEmpty")
		if c.PageBreakEmpty {
			c.PageBreakEmptyLine = getInt("PageBreakEmptyLine", 0)
			c.PageBreakEmptySize = getInt("PageBreakEmptySize", 0) * 1024
		} else {
			markKnown("PageBreakEmptyLine", "PageBreakEmptySize")
		}
		c.PageBreakChapter = getBool("PageBreakChapter")
		if c.PageBreakChapter {
			c.PageBreakChapterSize = getInt("PageBreakChapterSize", 0) * 1024
		} else {
			markKnown("PageBreakChapterSize")
		}
	} else {
		markKnown("PageBreakSize", "PageBreakEmpty", "PageBreakEmptyLine",
			"PageBreakEmptySize", "PageBreakChapter", "PageBreakChapterSize")
	}

	// 章設定
	c.ChapterNameLength = getInt("ChapterNameLength", c.ChapterNameLength)
	c.TitleToc = getBool("TitleToc")
	c.ChapterExclude = getBool("ChapterExclude")
	c.ChapterUseNextLine = getBool("ChapterUseNextLine")

	// ChapterSection: キー不在時 true（Java版 L236 互換）
	markKnown("ChapterSection")
	if v, ok := props["ChapterSection"]; ok {
		c.ChapterSection = v == "1"
	}

	c.ChapterH = getBool("ChapterH")
	c.ChapterH1 = getBool("ChapterH1")
	c.ChapterH2 = getBool("ChapterH2")
	c.ChapterH3 = getBool("ChapterH3")
	c.SameLineChapter = getBool("SameLineChapter")
	c.ChapterName = getBool("ChapterName")
	c.ChapterNumOnly = getBool("ChapterNumOnly")
	c.ChapterNumTitle = getBool("ChapterNumTitle")
	c.ChapterNumParen = getBool("ChapterNumParen")
	// Java版タイポ: hapterNumParenTitle
	c.ChapterNumParenTitle = getBool("ChapterNumParenTitle")
	if getBool("hapterNumParenTitle") {
		c.ChapterNumParenTitle = true
	}

	c.ChapterPattern = getBool("ChapterPattern")
	if c.ChapterPattern {
		c.ChapterPatternText = getString("ChapterPatternText", "")
	} else {
		markKnown("ChapterPatternText")
	}

	// その他
	c.NoIllust = getBool("NoIllust")
	c.DakutenType = getInt("DakutenType", c.DakutenType)
	c.IvsBMP = getBool("IvsBMP")
	c.IvsSSP = getBool("IvsSSP")

	// プリセット関連（無視だが既知キーとして登録）
	markKnown("Ext", "PresetName", "PresetIcon",
		"Gaiji32", "ImageFloatPage", "ImageFloatBlock")

	// 未知キーの警告
	for key := range props {
		if !known[key] {
			warnings = append(warnings, fmt.Sprintf("unknown key: %s", key))
		}
	}

	return c, warnings, nil
}

// parseProperties は Java Properties 形式を最小互換でパースする。
// #/! コメント、先頭空白トリム、= 区切りに対応。
func parseProperties(r io.Reader) (map[string]string, error) {
	props := make(map[string]string)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimLeft(scanner.Text(), " \t")
		if line == "" || line[0] == '#' || line[0] == '!' {
			continue
		}
		idx := strings.IndexByte(line, '=')
		if idx < 0 {
			continue
		}
		key := strings.TrimSpace(line[:idx])
		value := strings.TrimSpace(line[idx+1:])
		props[key] = value
	}
	return props, scanner.Err()
}
