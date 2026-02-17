package model

import "testing"

func TestGaijiInfo_FileName(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"/path/to/font.ttf", "font.ttf"},
		{"gaiji/kanji.png", "kanji.png"},
		{"single.svg", "single.svg"},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			g := &GaijiInfo{ClassName: "cls", FilePath: tt.path}
			if got := g.FileName(); got != tt.want {
				t.Errorf("FileName() = %q, want %q", got, tt.want)
			}
		})
	}
}
