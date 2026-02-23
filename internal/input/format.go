package input

import (
	"path/filepath"
	"strings"
)

// Format は入力ファイルの形式を表す。
type Format int

const (
	FormatTxt  Format = iota // プレーンテキスト
	FormatZip                // ZIPアーカイブ
	FormatTxtz               // TXTZアーカイブ (ZIP互換)
	FormatCbz                // CBZ (画像のみZIP)
	FormatRar                // RARアーカイブ
)

func (f Format) String() string {
	switch f {
	case FormatTxt:
		return "txt"
	case FormatZip:
		return "zip"
	case FormatTxtz:
		return "txtz"
	case FormatCbz:
		return "cbz"
	case FormatRar:
		return "rar"
	default:
		return "unknown"
	}
}

// DetectFormat はファイル名の拡張子からフォーマットを判定する。
func DetectFormat(filename string) (Format, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".txt":
		return FormatTxt, nil
	case ".zip":
		return FormatZip, nil
	case ".txtz":
		return FormatTxtz, nil
	case ".cbz":
		return FormatCbz, nil
	case ".rar":
		return FormatRar, nil
	default:
		return 0, ErrUnsupportedFormat
	}
}
