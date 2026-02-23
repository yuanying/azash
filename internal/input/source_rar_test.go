package input

import (
	"errors"
	"testing"
)

func TestRarSource(t *testing.T) {
	t.Run("Open rar returns ErrRarNotSupported", func(t *testing.T) {
		_, err := Open("test.rar")
		if !errors.Is(err, ErrRarNotSupported) {
			t.Errorf("Open(test.rar) error = %v, want %v", err, ErrRarNotSupported)
		}
	})
}
