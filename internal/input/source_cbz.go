package input

import "io"

// cbzSource はCBZアーカイブの Source 実装。
// 内部的にzipSourceを使用するが、常にimageOnlyとして扱う。
type cbzSource struct {
	zip *zipSource
}

func newCbzSource(filePath string) (*cbzSource, error) {
	zs, err := newZipSource(filePath, FormatCbz)
	if err != nil {
		return nil, err
	}
	return &cbzSource{zip: zs}, nil
}

func (s *cbzSource) Format() Format {
	return FormatCbz
}

func (s *cbzSource) TextEntries() []TextEntry {
	return nil
}

func (s *cbzSource) ImageOnly() bool {
	return true
}

func (s *cbzSource) OpenText(_ int) (io.ReadCloser, TextEntry, error) {
	return nil, TextEntry{}, ErrNoTextEntry
}

func (s *cbzSource) ImageEntries() []ImageEntry {
	return s.zip.ImageEntries()
}

func (s *cbzSource) ArchiveTextParentPath(_ int) string {
	return ""
}

func (s *cbzSource) Close() error {
	return s.zip.Close()
}
