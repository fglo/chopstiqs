package component

import (
	"image/color"
	"testing"
)

func TestDrawerCachesReuseBuffers(t *testing.T) {
	slider := &Slider{}
	slider.pixelRows = 4
	slider.pixelCols = 16
	slider.container = &Container{}
	slider.container.SetBackgroundColor(color.RGBA{10, 20, 30, 40})

	sliderDrawer := &DefaultSliderDrawer{}
	firstBuffer := sliderDrawer.getBuffer(slider)
	secondBuffer := sliderDrawer.getBuffer(slider)
	if &firstBuffer[0] != &secondBuffer[0] {
		t.Fatalf("expected slider drawer to reuse pixel buffer")
	}

	firstRow := sliderDrawer.getBgRow(slider)
	secondRow := sliderDrawer.getBgRow(slider)
	if &firstRow[0] != &secondRow[0] {
		t.Fatalf("expected slider drawer to reuse background row")
	}

	checkbox := &CheckBox{}
	checkbox.pixelRows = 4
	checkbox.pixelCols = 16
	checkbox.container = &Container{}
	checkbox.container.SetBackgroundColor(color.RGBA{10, 20, 30, 40})

	checkboxDrawer := &DefaultCheckBoxDrawer{}
	checkboxFirstBuffer := checkboxDrawer.getBuffer(checkbox)
	checkboxSecondBuffer := checkboxDrawer.getBuffer(checkbox)
	if &checkboxFirstBuffer[0] != &checkboxSecondBuffer[0] {
		t.Fatalf("expected checkbox drawer to reuse pixel buffer")
	}

	checkboxFirstRow := checkboxDrawer.getBgRow(checkbox)
	checkboxSecondRow := checkboxDrawer.getBgRow(checkbox)
	if &checkboxFirstRow[0] != &checkboxSecondRow[0] {
		t.Fatalf("expected checkbox drawer to reuse background row")
	}

	textInput := &TextInput{}
	textInput.pixelRows = 4
	textInput.pixelCols = 16
	textInput.container = &Container{}
	textInput.container.SetBackgroundColor(color.RGBA{10, 20, 30, 40})

	textInputDrawer := &DefaultTextInputDrawer{}
	textInputFirstBuffer := textInputDrawer.getBuffer(textInput)
	textInputSecondBuffer := textInputDrawer.getBuffer(textInput)
	if &textInputFirstBuffer[0] != &textInputSecondBuffer[0] {
		t.Fatalf("expected text input drawer to reuse pixel buffer")
	}

	textInputFirstRow := textInputDrawer.getBgRow(textInput)
	textInputSecondRow := textInputDrawer.getBgRow(textInput)
	if &textInputFirstRow[0] != &textInputSecondRow[0] {
		t.Fatalf("expected text input drawer to reuse background row")
	}

	cursor := &textInputCursor{}
	cursor.pixelRows = 4
	cursor.pixelCols = 16

	cursorDrawer := &DefaultTextInputCursorDrawer{}
	cursorFirstBuffer := cursorDrawer.getBuffer(cursor)
	cursorSecondBuffer := cursorDrawer.getBuffer(cursor)
	if &cursorFirstBuffer[0] != &cursorSecondBuffer[0] {
		t.Fatalf("expected text input cursor drawer to reuse pixel buffer")
	}

	scrollbar := &ScrollBar{}
	scrollbar.pixelRows = 4
	scrollbar.pixelCols = 16
	scrollbar.container = &Container{}
	scrollbar.container.SetBackgroundColor(color.RGBA{10, 20, 30, 40})

	scrollbarDrawer := &DefaultScrollBarDrawer{}
	scrollbarFirstBuffer := scrollbarDrawer.getBuffer(scrollbar)
	scrollbarSecondBuffer := scrollbarDrawer.getBuffer(scrollbar)
	if &scrollbarFirstBuffer[0] != &scrollbarSecondBuffer[0] {
		t.Fatalf("expected scrollbar drawer to reuse pixel buffer")
	}

	scrollbarFirstRow := scrollbarDrawer.getBgRow(scrollbar)
	scrollbarSecondRow := scrollbarDrawer.getBgRow(scrollbar)
	if &scrollbarFirstRow[0] != &scrollbarSecondRow[0] {
		t.Fatalf("expected scrollbar drawer to reuse background row")
	}
}
