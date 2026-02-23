package input

import (
	"testing"
)

func TestFileTitleCreator(t *testing.T) {
	tests := []struct {
		name        string
		fileName    string
		wantTitle   string
		wantCreator string
	}{
		// パターン1: [著者名] タイトル
		{
			name:        "bracket author title",
			fileName:    "[夏目漱石] 吾輩は猫である.txt",
			wantTitle:   "吾輩は猫である",
			wantCreator: "夏目漱石",
		},
		{
			name:        "fullwidth bracket author title",
			fileName:    "［夏目漱石］ 吾輩は猫である.txt",
			wantTitle:   "吾輩は猫である",
			wantCreator: "夏目漱石",
		},
		{
			name:        "bracket with fullwidth space",
			fileName:    "[夏目漱石]　こころ.zip",
			wantTitle:   "こころ",
			wantCreator: "夏目漱石",
		},

		// パターン2: タイトル （著者名） or タイトル (著者名)
		{
			name:        "title paren author",
			fileName:    "吾輩は猫である （夏目漱石）.txt",
			wantTitle:   "吾輩は猫である",
			wantCreator: "",
		},
		{
			name:        "title halfwidth paren author",
			fileName:    "吾輩は猫である (夏目漱石).txt",
			wantTitle:   "吾輩は猫である",
			wantCreator: "",
		},

		// パターン3: タイトルのみ（括弧なし）
		{
			name:        "title only",
			fileName:    "吾輩は猫である.txt",
			wantTitle:   "吾輩は猫である",
			wantCreator: "",
		},

		// 二重拡張子
		{
			name:        "double extension",
			fileName:    "[芥川龍之介] 羅生門.txt.zip",
			wantTitle:   "羅生門",
			wantCreator: "芥川龍之介",
		},

		// 青空文庫パターン除去
		{
			name:        "aozora pattern removed",
			fileName:    "吾輩は猫である（青空文庫）.txt",
			wantTitle:   "吾輩は猫である",
			wantCreator: "",
		},
		{
			name:        "aozora new spelling removed",
			fileName:    "[夏目漱石] 吾輩は猫である(青空文庫データ).zip",
			wantTitle:   "吾輩は猫である",
			wantCreator: "夏目漱石",
		},

		// 校正パターン除去
		{
			name:        "correction pattern removed",
			fileName:    "吾輩は猫である(校正済み).txt",
			wantTitle:   "吾輩は猫である",
			wantCreator: "",
		},
		{
			name:        "revision pattern removed",
			fileName:    "[太宰治] 人間失格(Rev2).txt",
			wantTitle:   "人間失格",
			wantCreator: "太宰治",
		},
		{
			name:        "ruby pattern removed",
			fileName:    "[太宰治] 人間失格(ルビ付き).zip",
			wantTitle:   "人間失格",
			wantCreator: "太宰治",
		},
		{
			name:        "cover pattern removed",
			fileName:    "テスト(表紙付き).txt",
			wantTitle:   "テスト",
			wantCreator: "",
		},
		{
			name:        "illustration pattern removed",
			fileName:    "テスト(挿絵付き).txt",
			wantTitle:   "テスト",
			wantCreator: "",
		},
		{
			name:        "lightweight pattern removed",
			fileName:    "テスト(軽量版).txt",
			wantTitle:   "テスト",
			wantCreator: "",
		},

		// 拡張子なしファイル名
		{
			name:        "no extension",
			fileName:    "テスト作品",
			wantTitle:   "テスト作品",
			wantCreator: "",
		},

		// 空タイトルにならないケース
		{
			name:        "bracket only creator empty title",
			fileName:    "[夏目漱石].txt",
			wantTitle:   "",
			wantCreator: "夏目漱石",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTitle, gotCreator := FileTitleCreator(tt.fileName)
			if gotTitle != tt.wantTitle {
				t.Errorf("FileTitleCreator(%q) title = %q, want %q", tt.fileName, gotTitle, tt.wantTitle)
			}
			if gotCreator != tt.wantCreator {
				t.Errorf("FileTitleCreator(%q) creator = %q, want %q", tt.fileName, gotCreator, tt.wantCreator)
			}
		})
	}
}
