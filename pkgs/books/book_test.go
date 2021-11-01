package books

import (
	"reflect"
	"testing"
)

func TestParseFilename(t *testing.T) {
	tests := []struct {
		desc     string
		filename string

		expectedArtists  []string
		expectedCategies []string
	}{
		{
			desc:     "no category and authors",
			filename: "comics.zip",
		},
		{
			desc:     "One category",
			filename: "(C94) comics.zip",

			expectedCategies: []string{"C94"},
		},
		{
			desc:     "Artist",
			filename: "[hoge] comics.zip",

			expectedArtists: []string{"hoge"},
		},
		{
			desc:     "Artist with group",
			filename: "[hoge (group)] comics.zip",

			expectedArtists: []string{"hoge", "group"},
		},
		{
			desc:     "Artist with group no blank",
			filename: "[hoge(group)]comics.zip",

			expectedArtists: []string{"hoge", "group"},
		},
		{
			desc:     "Category and Artist with group",
			filename: "(C200) [hoge (group)] comics.zip",

			expectedArtists:  []string{"hoge", "group"},
			expectedCategies: []string{"C200"},
		},
		{
			desc:     "Japanease",
			filename: "(日本語13) [ほげ (fuga)] コミックス (カテゴリ).zip",

			expectedArtists:  []string{"ほげ", "fuga"},
			expectedCategies: []string{"日本語13"},
		},
		{
			desc:     "Ignore last blacket",
			filename: "(C200) [hoge] comics [mustignore].zip",

			expectedArtists:  []string{"hoge"},
			expectedCategies: []string{"C200"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			book := Book{Filename: tt.filename}
			err := book.parseFilename()
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if got, want := book.Artists, tt.expectedArtists; !reflect.DeepEqual(got, want) {
				t.Errorf("book.Artists = %q, want = %q", got, want)
			}
			if got, want := book.Categories, tt.expectedCategies; !reflect.DeepEqual(got, want) {
				t.Errorf("book.Categories = %q, want = %q", got, want)
			}
		})
	}
}
