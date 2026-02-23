package input

import (
	"fmt"
	"io"
)

// TextEntry はテキストエントリの情報を保持する。
type TextEntry struct {
	Name string // エントリ名 (アーカイブ内パスまたはファイル名)
}

// ImageEntry はアーカイブ内の画像エントリの情報を保持する。
type ImageEntry struct {
	Name string // エントリ名 (アーカイブ内パス)
}

// Source は入力ソースの統一インタフェース。
type Source interface {
	Format() Format
	TextEntries() []TextEntry
	ImageOnly() bool
	OpenText(txtIdx int) (io.ReadCloser, TextEntry, error)
	ImageEntries() []ImageEntry
	ArchiveTextParentPath(txtIdx int) string
	Close() error
}

// Open は指定ファイルに対応する Source を生成するファクトリ関数。
func Open(filePath string) (Source, error) {
	format, err := DetectFormat(filePath)
	if err != nil {
		return nil, err
	}

	switch format {
	case FormatTxt:
		return newTxtSource(filePath)
	case FormatZip, FormatTxtz:
		return newZipSource(filePath, format)
	case FormatCbz:
		return newCbzSource(filePath)
	case FormatRar:
		return nil, ErrRarNotSupported
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedFormat, filePath)
	}
}
