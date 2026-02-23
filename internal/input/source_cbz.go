package input

import "io"

// cbzSource はCBZアーカイブの Source 実装。(Commit 6で実装)
type cbzSource struct{}

func newCbzSource(_ string) (*cbzSource, error) {
	panic("not implemented yet")
}

func (s *cbzSource) Format() Format                                      { panic("not implemented yet") }
func (s *cbzSource) TextEntries() []TextEntry                            { panic("not implemented yet") }
func (s *cbzSource) ImageOnly() bool                                     { panic("not implemented yet") }
func (s *cbzSource) OpenText(_ int) (io.ReadCloser, TextEntry, error)    { panic("not implemented yet") }
func (s *cbzSource) ImageEntries() []ImageEntry                          { panic("not implemented yet") }
func (s *cbzSource) ArchiveTextParentPath(_ int) string                  { panic("not implemented yet") }
func (s *cbzSource) Close() error                                        { panic("not implemented yet") }
