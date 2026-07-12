package component

import (
	"image/color"
	"testing"

	"github.com/fglo/chopstiqs/option"
	"github.com/matryer/is"
)

func TestLabel_NewLabel(t *testing.T) {
	is := is.New(t)

	l := NewLabel("hello", nil)

	is.Equal(l.text, "hello")
	is.Equal(l.Inverted, false)
}

func TestLabel_SetText(t *testing.T) {
	is := is.New(t)

	l := NewLabel("hello", nil)

	l.SetText("world")
	is.Equal(l.text, "world")
}

func TestLabel_SetText_UpdatesDimensions(t *testing.T) {
	is := is.New(t)

	l := NewLabel("hi", nil)
	w1, h1 := l.Dimensions()

	l.SetText("hello world")
	w2, h2 := l.Dimensions()

	is.True(w2 >= w1)
	is.True(h2 >= h1)
}

func TestLabel_InvertColor(t *testing.T) {
	is := is.New(t)

	l := NewLabel("hello", nil)
	is.Equal(l.Inverted, false)

	l.InvertColor()
	is.Equal(l.Inverted, true)

	l.InvertColor()
	is.Equal(l.Inverted, false)
}

func TestLabel_SetInverted(t *testing.T) {
	is := is.New(t)

	l := NewLabel("hello", nil)

	l.SetInverted(true)
	is.Equal(l.Inverted, true)

	l.SetInverted(false)
	is.Equal(l.Inverted, false)
}

func TestLabel_SetColor(t *testing.T) {
	is := is.New(t)

	l := NewLabel("hello", nil)

	red := color.RGBA{255, 0, 0, 255}
	l.SetColor(red)

	is.Equal(l.color, red)
}

func TestLabel_Dimensions_AfterSetText(t *testing.T) {
	is := is.New(t)

	l := NewLabel("hello", nil)
	w1, h1 := l.Dimensions()

	l.SetText("hello world again")
	w2, h2 := l.Dimensions()

	is.True(w2 > w1)
	is.True(h2 >= h1)
}

func TestLabel_WithOptions(t *testing.T) {
	is := is.New(t)

	red := color.RGBA{255, 0, 0, 255}
	l := NewLabel("hello", &LabelOptions{
		Color:               red,
		HorizontalAlignment: option.AlignmentCenteredHorizontally,
		VerticalAlignment:   option.AlignmentCenteredVertically,
	})

	is.Equal(l.text, "hello")
	is.Equal(l.color, red)
}

func TestLabel_EmptyText(t *testing.T) {
	is := is.New(t)

	l := NewLabel("", nil)
	is.Equal(l.text, "")
}
