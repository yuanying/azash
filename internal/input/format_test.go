package input

import (
	"testing"
)

func TestDetectFormat(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     Format
		wantErr  bool
	}{
		// txt
		{name: "txt lowercase", filename: "test.txt", want: FormatTxt},
		{name: "txt uppercase", filename: "test.TXT", want: FormatTxt},
		{name: "txt mixed case", filename: "test.Txt", want: FormatTxt},

		// zip
		{name: "zip lowercase", filename: "test.zip", want: FormatZip},
		{name: "zip uppercase", filename: "test.ZIP", want: FormatZip},

		// txtz
		{name: "txtz lowercase", filename: "test.txtz", want: FormatTxtz},
		{name: "txtz uppercase", filename: "test.TXTZ", want: FormatTxtz},

		// cbz
		{name: "cbz lowercase", filename: "test.cbz", want: FormatCbz},
		{name: "cbz uppercase", filename: "test.CBZ", want: FormatCbz},

		// rar
		{name: "rar lowercase", filename: "test.rar", want: FormatRar},
		{name: "rar uppercase", filename: "test.RAR", want: FormatRar},

		// unsupported
		{name: "epub unsupported", filename: "test.epub", wantErr: true},
		{name: "pdf unsupported", filename: "test.pdf", wantErr: true},
		{name: "no extension", filename: "testfile", wantErr: true},
		{name: "empty string", filename: "", wantErr: true},
		{name: "dot only", filename: ".", wantErr: true},

		// path with directory
		{name: "with path", filename: "/path/to/test.zip", want: FormatZip},
		{name: "with jp path", filename: "/書庫/テスト.txt", want: FormatTxt},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DetectFormat(tt.filename)
			if tt.wantErr {
				if err == nil {
					t.Errorf("DetectFormat(%q) expected error, got %v", tt.filename, got)
				}
				return
			}
			if err != nil {
				t.Errorf("DetectFormat(%q) unexpected error: %v", tt.filename, err)
				return
			}
			if got != tt.want {
				t.Errorf("DetectFormat(%q) = %v, want %v", tt.filename, got, tt.want)
			}
		})
	}
}

func TestFormat_String(t *testing.T) {
	tests := []struct {
		f    Format
		want string
	}{
		{FormatTxt, "txt"},
		{FormatZip, "zip"},
		{FormatTxtz, "txtz"},
		{FormatCbz, "cbz"},
		{FormatRar, "rar"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.f.String(); got != tt.want {
				t.Errorf("Format.String() = %q, want %q", got, tt.want)
			}
		})
	}
}
