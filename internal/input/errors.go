package input

import "errors"

var (
	// ErrUnsupportedFormat はサポート外のファイル形式で返される。
	ErrUnsupportedFormat = errors.New("unsupported format: only txt, zip, txtz, cbz, rar are supported")

	// ErrNoTextEntry はアーカイブ内にテキストエントリがない場合に返される。
	ErrNoTextEntry = errors.New("no text entry in archive")

	// ErrTextIndexOutOfRange はtxtIdxが範囲外の場合に返される。
	ErrTextIndexOutOfRange = errors.New("text index out of range")

	// ErrRarNotSupported はRAR形式が未対応であることを示す。
	ErrRarNotSupported = errors.New("rar format is not yet supported")
)
