package model

// 画像単ページ種別定数
const (
	ImagePageNone  = 0  // 画像単ページではない
	ImagePageW     = 1  // 単ページ 幅100%
	ImagePageH     = 2  // 単ページ 高さ100%
	ImagePageNoFit = 5  // 単ページ サイズ変更無し
	ImagePageAuto  = 10 // 単ページ サイズに応じて自動調整
)

// インライン画像種別定数
const (
	ImageInlineW       = 11 // 幅合わせ
	ImageInlineH       = 12 // 高さ合わせ
	ImageInlineTop     = 20 // 回り込み上
	ImageInlineBottom  = 21 // 回り込み下
	ImageInlineTopW    = 25 // 回り込み上 幅合わせ
	ImageInlineBottomW = 26 // 回り込み下 幅合わせ
)

// ページ種別定数
const (
	PageNormal = 0 // 通常
	PageMiddle = 1 // 中央寄せ
	PageBottom = 2 // 下付き
)

// PageBreakType は改ページ情報を格納する。
type PageBreakType struct {
	IgnoreEmptyPage bool
	PageType        int
	ImagePageType   int
	NoChapter       bool
	SrcFileName     string
	DstFileName     string
}

// NewPageBreakType は改ページ情報を生成する。
func NewPageBreakType(ignoreEmptyPage bool, pageType, imagePageType int) *PageBreakType {
	return &PageBreakType{
		IgnoreEmptyPage: ignoreEmptyPage,
		PageType:        pageType,
		ImagePageType:   imagePageType,
	}
}
