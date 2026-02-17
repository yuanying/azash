package model

// 画像サイズ種別定数
const (
	ImageSizeTypeAuto   = 1 // サイズ指定無し
	ImageSizeTypeHeight = 2 // 高さ指定
	ImageSizeTypeAspect = 3 // 画面縦横比
)

// SectionInfo はセクション (xhtml ファイル) の情報を格納する。
type SectionInfo struct {
	SectionID   string
	ImagePage   bool
	ImageHeight float64
	ImageFitW   bool
	ImageFitH   bool
	IsMiddle    bool
	IsBottom    bool
	StartLine   int
	EndLine     int
}

// NewSectionInfo はセクション情報を生成する。
func NewSectionInfo(sectionID string) *SectionInfo {
	return &SectionInfo{
		SectionID:   sectionID,
		ImageHeight: -1,
	}
}

// ImageHeightPercent は画像高さをパーセントで返す。
func (s *SectionInfo) ImageHeightPercent() float64 {
	return float64(int(s.ImageHeight*1000)) / 10.0
}

// ImageHeightPadding は画像高さに基づくパディングを返す。
func (s *SectionInfo) ImageHeightPadding() float64 {
	return float64(int((1-s.ImageHeight)*1000)) / 20.0
}
