package config

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
