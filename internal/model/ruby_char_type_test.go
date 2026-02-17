package model

import "testing"

func TestRubyCharType_IotaValues(t *testing.T) {
	tests := []struct {
		name string
		got  RubyCharType
		want int
	}{
		{"NULL", RubyCharNull, 0},
		{"KANJI", RubyCharKanji, 1},
		{"ALPHA", RubyCharAlpha, 2},
		{"FULLALPHA", RubyCharFullAlpha, 3},
		{"HIRAGANA", RubyCharHiragana, 4},
		{"KATAKANA", RubyCharKatakana, 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if int(tt.got) != tt.want {
				t.Errorf("RubyCharType %s = %d, want %d", tt.name, int(tt.got), tt.want)
			}
		})
	}
}
