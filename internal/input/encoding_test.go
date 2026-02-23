package input

import (
	"bytes"
	"testing"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func TestDetectEncoding(t *testing.T) {
	t.Run("UTF-8 text", func(t *testing.T) {
		text := "吾輩は猫である。名前はまだ無い。どこで生れたかとんと見当がつかぬ。"
		r := bytes.NewReader([]byte(text))
		enc, err := DetectEncoding(r)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if enc != "UTF-8" {
			t.Errorf("DetectEncoding = %q, want %q", enc, "UTF-8")
		}
	})

	t.Run("UTF-8 BOM text", func(t *testing.T) {
		bom := []byte{0xEF, 0xBB, 0xBF}
		text := append(bom, []byte("吾輩は猫である。名前はまだ無い。")...)
		r := bytes.NewReader(text)
		enc, err := DetectEncoding(r)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// BOM付きUTF-8は "UTF-8" として検出される
		if enc != "UTF-8" {
			t.Errorf("DetectEncoding = %q, want %q", enc, "UTF-8")
		}
	})

	t.Run("Shift_JIS text", func(t *testing.T) {
		// Shift_JISにエンコード
		encoder := japanese.ShiftJIS.NewEncoder()
		sjisText, _, err := transform.Bytes(encoder, []byte("吾輩は猫である。名前はまだ無い。どこで生れたかとんと見当がつかぬ。"))
		if err != nil {
			t.Fatalf("failed to encode to Shift_JIS: %v", err)
		}
		r := bytes.NewReader(sjisText)
		enc, err := DetectEncoding(r)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if enc != "Shift_JIS" {
			t.Errorf("DetectEncoding = %q, want %q", enc, "Shift_JIS")
		}
	})

	t.Run("empty input", func(t *testing.T) {
		r := bytes.NewReader([]byte{})
		enc, err := DetectEncoding(r)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if enc != "" {
			t.Errorf("DetectEncoding = %q, want empty string", enc)
		}
	})
}

func TestResolveEncoding(t *testing.T) {
	tests := []struct {
		name     string
		encType  string
		detected string
		want     string
	}{
		{name: "AUTO with UTF-8 detected", encType: "AUTO", detected: "UTF-8", want: "UTF-8"},
		{name: "AUTO with Shift_JIS detected", encType: "AUTO", detected: "Shift_JIS", want: "MS932"},
		{name: "AUTO with empty detected", encType: "AUTO", detected: "", want: "UTF-8"},
		{name: "fixed MS932", encType: "MS932", detected: "UTF-8", want: "MS932"},
		{name: "fixed UTF-8", encType: "UTF-8", detected: "Shift_JIS", want: "UTF-8"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ResolveEncoding(tt.encType, tt.detected)
			if got != tt.want {
				t.Errorf("ResolveEncoding(%q, %q) = %q, want %q", tt.encType, tt.detected, got, tt.want)
			}
		})
	}
}
