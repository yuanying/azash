package model

import "testing"

func TestCoverEditInfo_Literal(t *testing.T) {
	c := CoverEditInfo{
		PanelWidth:   800,
		PanelHeight:  600,
		FitType:      1,
		Scale:        1.5,
		OffsetX:      10.0,
		OffsetY:      20.0,
		VisibleWidth: 400.0,
	}
	if c.PanelWidth != 800 {
		t.Errorf("PanelWidth = %d, want 800", c.PanelWidth)
	}
	if c.PanelHeight != 600 {
		t.Errorf("PanelHeight = %d, want 600", c.PanelHeight)
	}
	if c.FitType != 1 {
		t.Errorf("FitType = %d, want 1", c.FitType)
	}
	if c.Scale != 1.5 {
		t.Errorf("Scale = %f, want 1.5", c.Scale)
	}
	if c.OffsetX != 10.0 {
		t.Errorf("OffsetX = %f, want 10.0", c.OffsetX)
	}
	if c.OffsetY != 20.0 {
		t.Errorf("OffsetY = %f, want 20.0", c.OffsetY)
	}
	if c.VisibleWidth != 400.0 {
		t.Errorf("VisibleWidth = %f, want 400.0", c.VisibleWidth)
	}
}
