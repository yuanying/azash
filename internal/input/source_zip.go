package input

import "io"

// zipSource はzip/txtzアーカイブの Source 実装。(Commit 5で実装)
type zipSource struct{}

func newZipSource(_ string, _ Format) (*zipSource, error) {
	panic("not implemented yet")
}

func (s *zipSource) Format() Format                                      { panic("not implemented yet") }
func (s *zipSource) TextEntries() []TextEntry                            { panic("not implemented yet") }
func (s *zipSource) ImageOnly() bool                                     { panic("not implemented yet") }
func (s *zipSource) OpenText(_ int) (io.ReadCloser, TextEntry, error)    { panic("not implemented yet") }
func (s *zipSource) ImageEntries() []ImageEntry                          { panic("not implemented yet") }
func (s *zipSource) ArchiveTextParentPath(_ int) string                  { panic("not implemented yet") }
func (s *zipSource) Close() error                                        { panic("not implemented yet") }
