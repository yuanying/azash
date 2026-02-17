package model

// CoverEditInfo は表紙の編集情報を格納する。
type CoverEditInfo struct {
	PanelWidth   int
	PanelHeight  int
	FitType      int
	Scale        float64
	OffsetX      float64
	OffsetY      float64
	VisibleWidth float64
}
