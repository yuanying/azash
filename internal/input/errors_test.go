package input

import (
	"errors"
	"testing"
)

func TestErrors(t *testing.T) {
	t.Run("ErrUnsupportedFormat is sentinel", func(t *testing.T) {
		err := ErrUnsupportedFormat
		if !errors.Is(err, ErrUnsupportedFormat) {
			t.Error("expected errors.Is to match ErrUnsupportedFormat")
		}
	})

	t.Run("ErrNoTextEntry is sentinel", func(t *testing.T) {
		err := ErrNoTextEntry
		if !errors.Is(err, ErrNoTextEntry) {
			t.Error("expected errors.Is to match ErrNoTextEntry")
		}
	})

	t.Run("ErrTextIndexOutOfRange is sentinel", func(t *testing.T) {
		err := ErrTextIndexOutOfRange
		if !errors.Is(err, ErrTextIndexOutOfRange) {
			t.Error("expected errors.Is to match ErrTextIndexOutOfRange")
		}
	})

	t.Run("ErrRarNotSupported is sentinel", func(t *testing.T) {
		err := ErrRarNotSupported
		if !errors.Is(err, ErrRarNotSupported) {
			t.Error("expected errors.Is to match ErrRarNotSupported")
		}
	})

	t.Run("DetectFormat returns ErrUnsupportedFormat for unknown", func(t *testing.T) {
		_, err := DetectFormat("test.pdf")
		if !errors.Is(err, ErrUnsupportedFormat) {
			t.Errorf("expected ErrUnsupportedFormat, got %v", err)
		}
	})
}
