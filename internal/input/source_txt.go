package input

import (
	"io"
	"os"
	"path/filepath"
)

// txtSource はプレーンテキストファイルの Source 実装。
type txtSource struct {
	filePath  string
	fileName  string
	parentDir string
}

func newTxtSource(filePath string) (*txtSource, error) {
	// ファイルの存在確認
	if _, err := os.Stat(filePath); err != nil {
		return nil, err
	}

	return &txtSource{
		filePath:  filePath,
		fileName:  filepath.Base(filePath),
		parentDir: filepath.Dir(filePath) + "/",
	}, nil
}

func (s *txtSource) Format() Format {
	return FormatTxt
}

func (s *txtSource) TextEntries() []TextEntry {
	return []TextEntry{{Name: s.fileName}}
}

func (s *txtSource) ImageOnly() bool {
	return false
}

func (s *txtSource) OpenText(txtIdx int) (io.ReadCloser, TextEntry, error) {
	if txtIdx != 0 {
		return nil, TextEntry{}, ErrTextIndexOutOfRange
	}
	f, err := os.Open(s.filePath)
	if err != nil {
		return nil, TextEntry{}, err
	}
	return f, TextEntry{Name: s.fileName}, nil
}

func (s *txtSource) ImageEntries() []ImageEntry {
	return nil
}

func (s *txtSource) ArchiveTextParentPath(_ int) string {
	return s.parentDir
}

func (s *txtSource) Close() error {
	return nil
}
