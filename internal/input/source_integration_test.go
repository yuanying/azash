package input

import (
	"io"
	"testing"
)

func TestIntegration_AllSamples(t *testing.T) {
	tests := []struct {
		name           string
		inputFile      string
		wantFormat     Format
		wantTxtCount   int
		wantImageOnly  bool
		wantImageCount int // -1 means don't check
	}{
		{
			name:          "sample-001 吾輩は猫である",
			inputFile:     "aozora_789_ruby.zip",
			wantFormat:    FormatZip,
			wantTxtCount:  1,
			wantImageOnly: false,
		},
		{
			name:          "sample-002 こころ",
			inputFile:     "aozora_773_ruby.zip",
			wantFormat:    FormatZip,
			wantTxtCount:  1,
			wantImageOnly: false,
		},
		{
			name:          "sample-003 羅生門",
			inputFile:     "aozora_127_ruby.zip",
			wantFormat:    FormatZip,
			wantTxtCount:  1,
			wantImageOnly: false,
		},
		{
			name:          "sample-004 銀河鉄道の夜初期形",
			inputFile:     "aozora_60681_ruby.zip",
			wantFormat:    FormatZip,
			wantTxtCount:  1,
			wantImageOnly: false,
		},
		{
			name:          "sample-005 実さんの精神分析",
			inputFile:     "aozora_46705_ruby.zip",
			wantFormat:    FormatZip,
			wantTxtCount:  1,
			wantImageOnly: false,
		},
		{
			name:           "sample-006 test_png (image-only)",
			inputFile:      "test_png.zip",
			wantFormat:     FormatZip,
			wantTxtCount:   0,
			wantImageOnly:  true,
			wantImageCount: 11,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src, err := Open(testdataPath(tt.inputFile))
			if err != nil {
				t.Fatalf("Open(%q) error: %v", tt.inputFile, err)
			}
			defer func() { _ = src.Close() }()

			if src.Format() != tt.wantFormat {
				t.Errorf("Format() = %v, want %v", src.Format(), tt.wantFormat)
			}

			txtEntries := src.TextEntries()
			if len(txtEntries) != tt.wantTxtCount {
				t.Errorf("TextEntries() count = %d, want %d", len(txtEntries), tt.wantTxtCount)
			}

			if src.ImageOnly() != tt.wantImageOnly {
				t.Errorf("ImageOnly() = %v, want %v", src.ImageOnly(), tt.wantImageOnly)
			}

			if tt.wantImageCount > 0 {
				images := src.ImageEntries()
				if len(images) != tt.wantImageCount {
					t.Errorf("ImageEntries() count = %d, want %d", len(images), tt.wantImageCount)
				}
			}
		})
	}
}

func TestIntegration_TextReadAndEncoding(t *testing.T) {
	tests := []struct {
		name      string
		inputFile string
		wantEnc   string
	}{
		{
			name:      "sample-003 羅生門 encoding detection",
			inputFile: "aozora_127_ruby.zip",
			wantEnc:   "Shift_JIS",
		},
		{
			name:      "sample-001 吾輩は猫である encoding detection",
			inputFile: "aozora_789_ruby.zip",
			wantEnc:   "Shift_JIS",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src, err := Open(testdataPath(tt.inputFile))
			if err != nil {
				t.Fatalf("Open error: %v", err)
			}
			defer func() { _ = src.Close() }()

			rc, _, err := src.OpenText(0)
			if err != nil {
				t.Fatalf("OpenText(0) error: %v", err)
			}
			defer func() { _ = rc.Close() }()

			enc, err := DetectEncoding(rc)
			if err != nil {
				t.Fatalf("DetectEncoding error: %v", err)
			}
			if enc != tt.wantEnc {
				t.Errorf("DetectEncoding = %q, want %q", enc, tt.wantEnc)
			}

			resolved := ResolveEncoding("AUTO", enc)
			if resolved != "MS932" {
				t.Errorf("ResolveEncoding(AUTO, %q) = %q, want MS932", enc, resolved)
			}
		})
	}
}

func TestIntegration_OpenTextAndReadFull(t *testing.T) {
	src, err := Open(testdataPath("aozora_127_ruby.zip"))
	if err != nil {
		t.Fatalf("Open error: %v", err)
	}
	defer func() { _ = src.Close() }()

	rc, entry, err := src.OpenText(0)
	if err != nil {
		t.Fatalf("OpenText(0) error: %v", err)
	}
	defer func() { _ = rc.Close() }()

	if entry.Name != "rashomon.txt" {
		t.Errorf("entry.Name = %q, want %q", entry.Name, "rashomon.txt")
	}

	data, err := io.ReadAll(rc)
	if err != nil {
		t.Fatalf("ReadAll error: %v", err)
	}

	// rashomon.txt は約13KB
	if len(data) < 10000 {
		t.Errorf("data length = %d, expected > 10000", len(data))
	}
}
