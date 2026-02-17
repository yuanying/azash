package model

import "strings"

// ImageInfo は画像情報を格納する。
type ImageInfo struct {
	ID          string
	OutFileName string
	Ext         string
	Width       int
	Height      int
	OutWidth    int
	OutHeight   int
	IsCover     bool
	RotateAngle int
}

// NewImageInfo は画像情報を生成する。ext は小文字化される。
func NewImageInfo(ext string, width, height int) *ImageInfo {
	return &ImageInfo{
		Ext:       strings.ToLower(ext),
		Width:     width,
		Height:    height,
		OutWidth:  -1,
		OutHeight: -1,
	}
}

// Format は MIME 形式のフォーマット文字列を返す (例: "image/png")。
func (i *ImageInfo) Format() string {
	ext := i.Ext
	if ext == "jpg" {
		ext = "jpeg"
	}
	return "image/" + ext
}
