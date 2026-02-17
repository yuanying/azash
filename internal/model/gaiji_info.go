package model

import "path/filepath"

// GaijiInfo は外字フォントとクラス名を格納する。
type GaijiInfo struct {
	ClassName string
	FilePath  string
}

// FileName はファイルパスからベース名を返す。
func (g *GaijiInfo) FileName() string {
	return filepath.Base(g.FilePath)
}
