package model

// RubyCharType はルビ文字の種別を表す。
type RubyCharType int

const (
	RubyCharNull      RubyCharType = iota // NULL
	RubyCharKanji                         // KANJI
	RubyCharAlpha                         // ALPHA
	RubyCharFullAlpha                     // FULLALPHA
	RubyCharHiragana                      // HIRAGANA
	RubyCharKatakana                      // KATAKANA
)
