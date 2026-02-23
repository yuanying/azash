package input

import (
	"archive/zip"
	"io"
	"path"
	"strings"

	"golang.org/x/text/encoding/japanese"
)

// 画像拡張子セット (Java版 ImageInfoReader 互換: png/jpg/jpeg/gif/webp)
var imageExts = map[string]bool{
	".png": true, ".jpg": true, ".jpeg": true, ".gif": true, ".webp": true,
}

// zipSource はzip/txtzアーカイブの Source 実装。
type zipSource struct {
	filePath     string
	format       Format
	reader       *zip.ReadCloser
	textEntries  []TextEntry
	imageEntries []ImageEntry
	textFiles    []*zip.File // テキストエントリに対応するzip.File
}

func newZipSource(filePath string, format Format) (*zipSource, error) {
	r, err := zip.OpenReader(filePath)
	if err != nil {
		return nil, err
	}

	s := &zipSource{
		filePath: filePath,
		format:   format,
		reader:   r,
	}

	s.scanEntries()
	return s, nil
}

// scanEntries はアーカイブ内のエントリをスキャンしてtxtと画像に分類する。
func (s *zipSource) scanEntries() {
	for _, f := range s.reader.File {
		name := s.decodeName(f)

		// ディレクトリはスキップ
		if f.FileInfo().IsDir() {
			continue
		}

		ext := strings.ToLower(path.Ext(name))
		if ext == ".txt" {
			s.textEntries = append(s.textEntries, TextEntry{Name: name})
			s.textFiles = append(s.textFiles, f)
		} else if imageExts[ext] {
			s.imageEntries = append(s.imageEntries, ImageEntry{Name: name})
		}
	}
}

// decodeName はzip.Fileのファイル名をデコードする。
// NonUTF8フラグがセットされている場合はMS932 (Shift_JIS) としてデコードする。
func (s *zipSource) decodeName(f *zip.File) string {
	if f.NonUTF8 {
		decoded, err := japanese.ShiftJIS.NewDecoder().String(f.Name)
		if err == nil {
			return decoded
		}
	}
	return f.Name
}

func (s *zipSource) Format() Format {
	return s.format
}

func (s *zipSource) TextEntries() []TextEntry {
	return s.textEntries
}

func (s *zipSource) ImageOnly() bool {
	return len(s.textEntries) == 0
}

func (s *zipSource) OpenText(txtIdx int) (io.ReadCloser, TextEntry, error) {
	if len(s.textEntries) == 0 {
		return nil, TextEntry{}, ErrNoTextEntry
	}
	if txtIdx < 0 || txtIdx >= len(s.textEntries) {
		return nil, TextEntry{}, ErrTextIndexOutOfRange
	}

	rc, err := s.textFiles[txtIdx].Open()
	if err != nil {
		return nil, TextEntry{}, err
	}
	return rc, s.textEntries[txtIdx], nil
}

func (s *zipSource) ImageEntries() []ImageEntry {
	return s.imageEntries
}

func (s *zipSource) ArchiveTextParentPath(txtIdx int) string {
	if txtIdx < 0 || txtIdx >= len(s.textEntries) {
		return ""
	}
	dir := path.Dir(s.textEntries[txtIdx].Name)
	if dir == "." {
		return ""
	}
	return dir + "/"
}

func (s *zipSource) Close() error {
	if s.reader != nil {
		return s.reader.Close()
	}
	return nil
}
