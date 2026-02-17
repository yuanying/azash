package model

import "testing"

func TestNewImageInfo_Defaults(t *testing.T) {
	img := NewImageInfo("PNG", 100, 200)
	if img.Ext != "png" {
		t.Errorf("Ext = %q, want %q", img.Ext, "png")
	}
	if img.Width != 100 {
		t.Errorf("Width = %d, want 100", img.Width)
	}
	if img.Height != 200 {
		t.Errorf("Height = %d, want 200", img.Height)
	}
	if img.OutWidth != -1 {
		t.Errorf("OutWidth = %d, want -1", img.OutWidth)
	}
	if img.OutHeight != -1 {
		t.Errorf("OutHeight = %d, want -1", img.OutHeight)
	}
	if img.IsCover {
		t.Error("IsCover should be false")
	}
	if img.RotateAngle != 0 {
		t.Errorf("RotateAngle = %d, want 0", img.RotateAngle)
	}
}

func TestImageInfo_Format(t *testing.T) {
	tests := []struct {
		ext  string
		want string
	}{
		{"png", "image/png"},
		{"jpg", "image/jpeg"},
		{"gif", "image/gif"},
		{"jpeg", "image/jpeg"},
	}
	for _, tt := range tests {
		t.Run(tt.ext, func(t *testing.T) {
			img := NewImageInfo(tt.ext, 1, 1)
			if got := img.Format(); got != tt.want {
				t.Errorf("Format() = %q, want %q", got, tt.want)
			}
		})
	}
}
