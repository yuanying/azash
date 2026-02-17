package model

import "testing"

func TestNewSectionInfo_Defaults(t *testing.T) {
	s := NewSectionInfo("0001")
	if s.SectionID != "0001" {
		t.Errorf("SectionID = %q, want %q", s.SectionID, "0001")
	}
	if s.ImageHeight != -1 {
		t.Errorf("ImageHeight = %f, want -1", s.ImageHeight)
	}
	if s.ImagePage {
		t.Error("ImagePage should be false")
	}
	if s.ImageFitW {
		t.Error("ImageFitW should be false")
	}
	if s.ImageFitH {
		t.Error("ImageFitH should be false")
	}
	if s.IsMiddle {
		t.Error("IsMiddle should be false")
	}
	if s.IsBottom {
		t.Error("IsBottom should be false")
	}
	if s.StartLine != 0 {
		t.Errorf("StartLine = %d, want 0", s.StartLine)
	}
	if s.EndLine != 0 {
		t.Errorf("EndLine = %d, want 0", s.EndLine)
	}
}

func TestSectionInfo_ImageHeightPercent(t *testing.T) {
	tests := []struct {
		name   string
		height float64
		want   float64
	}{
		{"75%", 0.75, 75.0},
		{"100%", 1.0, 100.0},
		{"50%", 0.5, 50.0},
		{"75.6%", 0.756, 75.6},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SectionInfo{ImageHeight: tt.height}
			if got := s.ImageHeightPercent(); got != tt.want {
				t.Errorf("ImageHeightPercent() = %f, want %f", got, tt.want)
			}
		})
	}
}

func TestSectionInfo_ImageHeightPadding(t *testing.T) {
	tests := []struct {
		name   string
		height float64
		want   float64
	}{
		{"75% → 12.5", 0.75, 12.5},
		{"100% → 0", 1.0, 0.0},
		{"50% → 25", 0.5, 25.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SectionInfo{ImageHeight: tt.height}
			if got := s.ImageHeightPadding(); got != tt.want {
				t.Errorf("ImageHeightPadding() = %f, want %f", got, tt.want)
			}
		})
	}
}

func TestSectionInfo_ImageSizeTypeConstants(t *testing.T) {
	if ImageSizeTypeAuto != 1 {
		t.Errorf("ImageSizeTypeAuto = %d, want 1", ImageSizeTypeAuto)
	}
	if ImageSizeTypeHeight != 2 {
		t.Errorf("ImageSizeTypeHeight = %d, want 2", ImageSizeTypeHeight)
	}
	if ImageSizeTypeAspect != 3 {
		t.Errorf("ImageSizeTypeAspect = %d, want 3", ImageSizeTypeAspect)
	}
}
