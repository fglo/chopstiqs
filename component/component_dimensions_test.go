package component

import (
	"testing"

	"github.com/fglo/chopstiqs/event"
	"github.com/fglo/chopstiqs/option"
	"github.com/matryer/is"
)

func newTestComponent(width, height int) *Button {
	eventManager := event.NewManager()
	b := NewButton(&ButtonOptions{
		Width:  option.Int(width),
		Height: option.Int(height),
	})
	b.SetEventManager(eventManager)
	return b
}

func TestComponent_SetDimensions(t *testing.T) {
	is := is.New(t)

	b := newTestComponent(50, 30)

	w, h := b.Dimensions()
	is.Equal(w, 50)
	is.Equal(h, 30)
}

func TestComponent_SetWidth(t *testing.T) {
	is := is.New(t)

	b := newTestComponent(50, 30)
	b.SetWidth(100)

	is.Equal(b.Width(), 100)
}

func TestComponent_SetHeight(t *testing.T) {
	is := is.New(t)

	b := newTestComponent(50, 30)
	b.SetHeight(80)

	is.Equal(b.Height(), 80)
}

func TestComponent_SetDimensions_IncludesPadding(t *testing.T) {
	is := is.New(t)

	p := NewPadding(5, 10, 5, 10)
	b := NewButton(&ButtonOptions{
		Width:   option.Int(50),
		Height:  option.Int(30),
		Padding: p,
	})

	w, h := b.Dimensions()
	is.Equal(w, 70)
	is.Equal(h, 40)
}

func TestComponent_SetWidth_UpdatesWidthWithPadding(t *testing.T) {
	is := is.New(t)

	p := NewPadding(0, 10, 0, 10)
	b := NewButton(&ButtonOptions{
		Width:   option.Int(50),
		Height:  option.Int(30),
		Padding: p,
	})

	b.SetWidth(100)

	is.Equal(b.WidthWithPadding(), 120)
}

func TestComponent_SetHeight_UpdatesHeightWithPadding(t *testing.T) {
	is := is.New(t)

	p := NewPadding(5, 0, 5, 0)
	b := NewButton(&ButtonOptions{
		Width:   option.Int(50),
		Height:  option.Int(30),
		Padding: p,
	})

	b.SetHeight(80)

	is.Equal(b.HeightWithPadding(), 90)
}

func TestComponent_SetPadding(t *testing.T) {
	is := is.New(t)

	b := newTestComponent(50, 30)

	newPad := NewPadding(5, 10, 5, 10)
	b.SetPadding(*newPad)

	pad := b.Padding()
	is.Equal(pad.Top, 5)
	is.Equal(pad.Right, 10)
	is.Equal(pad.Bottom, 5)
	is.Equal(pad.Left, 10)
}

func TestComponent_SetPaddingTop(t *testing.T) {
	is := is.New(t)

	b := newTestComponent(50, 30)
	b.SetPaddingTop(15)

	pad := b.Padding()
	is.Equal(pad.Top, 15)
}

func TestComponent_SetPaddingBottom(t *testing.T) {
	is := is.New(t)

	b := newTestComponent(50, 30)
	b.SetPaddingBottom(20)

	pad := b.Padding()
	is.Equal(pad.Bottom, 20)
}

func TestComponent_SetPaddingLeft(t *testing.T) {
	is := is.New(t)

	b := newTestComponent(50, 30)
	b.SetPaddingLeft(25)

	pad := b.Padding()
	is.Equal(pad.Left, 25)
}

func TestComponent_SetPaddingRight(t *testing.T) {
	is := is.New(t)

	b := newTestComponent(50, 30)
	b.SetPaddingRight(30)

	pad := b.Padding()
	is.Equal(pad.Right, 30)
}

func TestComponent_SetPaddingTop_NegativeUsesDefault(t *testing.T) {
	is := is.New(t)

	origDefault := DefaultPadding
	defer func() { DefaultPadding = origDefault }()

	DefaultPadding.Top = 7

	b := newTestComponent(50, 30)
	b.SetPaddingTop(-1)

	pad := b.Padding()
	is.Equal(pad.Top, 7)
}

func TestComponent_SetPaddingBottom_NegativeUsesDefault(t *testing.T) {
	is := is.New(t)

	origDefault := DefaultPadding
	defer func() { DefaultPadding = origDefault }()

	DefaultPadding.Bottom = 8

	b := newTestComponent(50, 30)
	b.SetPaddingBottom(-1)

	pad := b.Padding()
	is.Equal(pad.Bottom, 8)
}

func TestComponent_SetPaddingLeft_NegativeUsesDefault(t *testing.T) {
	is := is.New(t)

	origDefault := DefaultPadding
	defer func() { DefaultPadding = origDefault }()

	DefaultPadding.Left = 9

	b := newTestComponent(50, 30)
	b.SetPaddingLeft(-1)

	pad := b.Padding()
	is.Equal(pad.Left, 9)
}

func TestComponent_SetPaddingRight_NegativeUsesDefault(t *testing.T) {
	is := is.New(t)

	origDefault := DefaultPadding
	defer func() { DefaultPadding = origDefault }()

	DefaultPadding.Right = 10

	b := newTestComponent(50, 30)
	b.SetPaddingRight(-1)

	pad := b.Padding()
	is.Equal(pad.Right, 10)
}

func TestComponent_Position(t *testing.T) {
	is := is.New(t)

	b := newTestComponent(50, 30)

	x, y := b.Position()
	is.Equal(x, 0.)
	is.Equal(y, 0.)

	b.SetPosition(10, 20)
	x, y = b.Position()
	is.Equal(x, 10.)
	is.Equal(y, 20.)
}

func TestComponent_SetPosX(t *testing.T) {
	is := is.New(t)

	b := newTestComponent(50, 30)
	b.SetPosX(42)

	is.Equal(b.PosX(), 42.)
}

func TestComponent_SetPosY(t *testing.T) {
	is := is.New(t)

	b := newTestComponent(50, 30)
	b.SetPosY(42)

	is.Equal(b.PosY(), 42.)
}

func TestComponent_AbsPosition_NoContainer(t *testing.T) {
	is := is.New(t)

	b := newTestComponent(50, 30)
	b.SetPosition(10, 20)

	absX, absY := b.AbsPosition()
	is.Equal(absX, 10.)
	is.Equal(absY, 20.)
}

func TestComponent_AbsPosition_WithContainer(t *testing.T) {
	is := is.New(t)

	eventManager := event.NewManager()
	c := NewContainer(&ContainerOptions{})
	c.SetEventManager(eventManager)
	c.SetPosition(100, 200)

	b := newTestComponent(50, 30)
	c.AddComponent(b)
	b.SetPosition(10, 20)

	absX, absY := b.AbsPosition()
	is.Equal(absX, 110.)
	is.Equal(absY, 220.)
}

func TestComponent_SetDisabled(t *testing.T) {
	is := is.New(t)

	b := newTestComponent(50, 30)

	is.Equal(b.Disable(), false)

	b.SetDisabled(true)
	is.Equal(b.Disable(), true)

	b.SetDisabled(false)
	is.Equal(b.Disable(), false)
}

func TestComponent_SetHidden(t *testing.T) {
	is := is.New(t)

	b := newTestComponent(50, 30)

	is.Equal(b.Hidden(), false)

	b.SetHidden(true)
	is.Equal(b.Hidden(), true)

	b.SetHidden(false)
	is.Equal(b.Hidden(), false)
}

func TestComponent_SetFocused_FiresEvent(t *testing.T) {
	is := is.New(t)

	b := newTestComponent(50, 30)

	focusedCount := 0
	b.AddFocusedHandler(func(args *ComponentFocusedEventArgs) {
		focusedCount++
	})

	b.SetFocused(true)
	b.eventManager.HandleFired()
	is.Equal(focusedCount, 1)

	b.SetFocused(true)
	b.eventManager.HandleFired()
	is.Equal(focusedCount, 1)

	b.SetFocused(false)
	b.eventManager.HandleFired()
	is.Equal(focusedCount, 2)
}

func TestComponent_MinWidth(t *testing.T) {
	is := is.New(t)

	b := newTestComponent(50, 30)

	is.Equal(b.MinWidth(), 0)
}

func TestComponent_MinHeight(t *testing.T) {
	is := is.New(t)

	b := newTestComponent(50, 30)

	is.Equal(b.MinHeight(), 0)
}
