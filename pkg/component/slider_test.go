package component

import (
	"testing"

	"github.com/fglo/chopstiqs/pkg/event"
)

func TestSlider_calcHandlePosition(t *testing.T) {
	type fields struct {
		component             component
		min                   float64
		max                   float64
		step                  float64
		stepPixels            float64
		value                 float64
		pressed               bool
		hovering              bool
		handle                *Button
		sliding               bool
		SlidedEvent           *event.Event
		PressedEvent          *event.Event
		ReleasedEvent         *event.Event
		ClickedEvent          *event.Event
		firstPixelRowId       int
		secondPixelRowId      int
		lastPixelRowId        int
		penultimatePixelRowId int
		firstPixelColId       int
		secondPixelColId      int
		lastPixelColId        int
		penultimatePixelColId int
		drawer                SliderDrawer
		handleDrawer          ButtonDrawer
	}
	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Slider{
				component:             tt.fields.component,
				min:                   tt.fields.min,
				max:                   tt.fields.max,
				step:                  tt.fields.step,
				stepPixels:            tt.fields.stepPixels,
				value:                 tt.fields.value,
				pressed:               tt.fields.pressed,
				hovering:              tt.fields.hovering,
				handle:                tt.fields.handle,
				sliding:               tt.fields.sliding,
				SlidedEvent:           tt.fields.SlidedEvent,
				PressedEvent:          tt.fields.PressedEvent,
				ReleasedEvent:         tt.fields.ReleasedEvent,
				ClickedEvent:          tt.fields.ClickedEvent,
				firstPixelRowId:       tt.fields.firstPixelRowId,
				secondPixelRowId:      tt.fields.secondPixelRowId,
				lastPixelRowId:        tt.fields.lastPixelRowId,
				penultimatePixelRowId: tt.fields.penultimatePixelRowId,
				firstPixelColId:       tt.fields.firstPixelColId,
				secondPixelColId:      tt.fields.secondPixelColId,
				lastPixelColId:        tt.fields.lastPixelColId,
				penultimatePixelColId: tt.fields.penultimatePixelColId,
				drawer:                tt.fields.drawer,
				handleDrawer:          tt.fields.handleDrawer,
			}
			if got := s.calcHandlePosition(); got != tt.want {
				t.Errorf("Slider.calcHandlePosition() = %v, want %v", got, tt.want)
			}
		})
	}
}
