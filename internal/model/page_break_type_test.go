package model

import "testing"

func TestNewPageBreakType_Defaults(t *testing.T) {
	p := NewPageBreakType(true, PageMiddle, ImagePageW)
	if !p.IgnoreEmptyPage {
		t.Error("IgnoreEmptyPage should be true")
	}
	if p.PageType != PageMiddle {
		t.Errorf("PageType = %d, want %d", p.PageType, PageMiddle)
	}
	if p.ImagePageType != ImagePageW {
		t.Errorf("ImagePageType = %d, want %d", p.ImagePageType, ImagePageW)
	}
	if p.NoChapter {
		t.Error("NoChapter should be false by default")
	}
	if p.SrcFileName != "" {
		t.Errorf("SrcFileName = %q, want empty", p.SrcFileName)
	}
	if p.DstFileName != "" {
		t.Errorf("DstFileName = %q, want empty", p.DstFileName)
	}
}

func TestPageBreakType_NoChapter(t *testing.T) {
	p := &PageBreakType{
		IgnoreEmptyPage: true,
		PageType:        PageNormal,
		ImagePageType:   ImagePageNone,
		NoChapter:       true,
	}
	if !p.NoChapter {
		t.Error("NoChapter should be true")
	}
}

func TestPageBreakType_Constants(t *testing.T) {
	// 定数値が Java 版と一致していることを確認
	if ImagePageNone != 0 {
		t.Errorf("ImagePageNone = %d, want 0", ImagePageNone)
	}
	if ImagePageW != 1 {
		t.Errorf("ImagePageW = %d, want 1", ImagePageW)
	}
	if ImagePageH != 2 {
		t.Errorf("ImagePageH = %d, want 2", ImagePageH)
	}
	if ImagePageNoFit != 5 {
		t.Errorf("ImagePageNoFit = %d, want 5", ImagePageNoFit)
	}
	if ImagePageAuto != 10 {
		t.Errorf("ImagePageAuto = %d, want 10", ImagePageAuto)
	}
	if PageNormal != 0 {
		t.Errorf("PageNormal = %d, want 0", PageNormal)
	}
	if PageMiddle != 1 {
		t.Errorf("PageMiddle = %d, want 1", PageMiddle)
	}
	if PageBottom != 2 {
		t.Errorf("PageBottom = %d, want 2", PageBottom)
	}
}
