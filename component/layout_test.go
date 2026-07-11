package component

import (
	"testing"

	"github.com/fglo/chopstiqs/event"
	"github.com/fglo/chopstiqs/option"
	"github.com/matryer/is"
)

func TestHorizontalListLayout_Arrange(t *testing.T) {
	is := is.New(t)

	eventManager := event.NewManager()
	hl := &HorizontalListLayout{ColumnGap: 5}

	c := NewContainer(&ContainerOptions{Layout: hl})
	c.SetEventManager(eventManager)
	c.SetDimensions(200, 50)

	b1 := NewButton(&ButtonOptions{Width: option.Int(30), Height: option.Int(20)})
	b2 := NewButton(&ButtonOptions{Width: option.Int(40), Height: option.Int(20)})
	c.AddComponents(b1, b2)

	x1, _ := b1.Position()
	x2, _ := b2.Position()

	is.True(x1 >= 0)
	is.True(x2 > x1)
}

func TestHorizontalListLayout_ColumnGap(t *testing.T) {
	is := is.New(t)

	eventManager := event.NewManager()
	gap := 10
	hl := &HorizontalListLayout{ColumnGap: gap}

	c := NewContainer(&ContainerOptions{Layout: hl})
	c.SetEventManager(eventManager)
	c.SetDimensions(200, 50)

	b1 := NewButton(&ButtonOptions{Width: option.Int(30), Height: option.Int(20)})
	b2 := NewButton(&ButtonOptions{Width: option.Int(40), Height: option.Int(20)})
	c.AddComponents(b1, b2)

	x1, _ := b1.Position()
	x2, _ := b2.Position()

	expectedX2 := x1 + float64(b1.WidthWithPadding()) + float64(gap)
	is.Equal(x2, expectedX2)
}

func TestHorizontalListLayout_VerticalAlignment_Top(t *testing.T) {
	is := is.New(t)

	eventManager := event.NewManager()
	hl := &HorizontalListLayout{ColumnGap: 0}

	c := NewContainer(&ContainerOptions{Layout: hl})
	c.SetEventManager(eventManager)
	c.SetDimensions(200, 50)

	b := NewButton(&ButtonOptions{
		Width:             option.Int(30),
		Height:            option.Int(20),
		VerticalAlignment: option.AlignmentTop,
	})
	c.AddComponent(b)

	_, y := b.Position()
	is.Equal(y, float64(c.padding.Top))
}

func TestHorizontalListLayout_VerticalAlignment_Bottom(t *testing.T) {
	is := is.New(t)

	eventManager := event.NewManager()
	hl := &HorizontalListLayout{ColumnGap: 0}

	c := NewContainer(&ContainerOptions{Layout: hl})
	c.SetEventManager(eventManager)
	c.SetDimensions(200, 50)

	tallButton := NewButton(&ButtonOptions{Width: option.Int(30), Height: option.Int(40)})
	shortButton := NewButton(&ButtonOptions{
		Width:             option.Int(30),
		Height:            option.Int(20),
		VerticalAlignment: option.AlignmentBottom,
	})
	c.AddComponents(tallButton, shortButton)

	_, yShort := shortButton.Position()
	_, yTall := tallButton.Position()

	is.True(yShort > yTall)
}

func TestHorizontalListLayout_VerticalAlignment_Centered(t *testing.T) {
	is := is.New(t)

	eventManager := event.NewManager()
	hl := &HorizontalListLayout{ColumnGap: 0}

	c := NewContainer(&ContainerOptions{Layout: hl})
	c.SetEventManager(eventManager)
	c.SetDimensions(200, 50)

	tallButton := NewButton(&ButtonOptions{Width: option.Int(30), Height: option.Int(40)})
	centeredButton := NewButton(&ButtonOptions{
		Width:             option.Int(30),
		Height:            option.Int(20),
		VerticalAlignment: option.AlignmentCenteredVertically,
	})
	c.AddComponents(tallButton, centeredButton)

	_, yTall := tallButton.Position()
	_, yCentered := centeredButton.Position()

	is.True(yCentered > yTall)
}

func TestHorizontalListLayout_StretchedHorizontally(t *testing.T) {
	is := is.New(t)

	eventManager := event.NewManager()
	hl := &HorizontalListLayout{ColumnGap: 0}

	c := NewContainer(&ContainerOptions{Layout: hl})
	c.SetEventManager(eventManager)
	c.SetDimensions(100, 30)

	fixedButton := NewButton(&ButtonOptions{
		Width:  option.Int(20),
		Height: option.Int(20),
	})
	stretchedButton := NewButton(&ButtonOptions{
		Width:               option.Int(20),
		Height:              option.Int(20),
		HorizontalAlignment: option.StretchedHorizontally,
	})
	c.AddComponents(fixedButton, stretchedButton)

	c.SetDimensions(100, 30)
	hl.Rearrange(c)

	is.True(stretchedButton.Width() > 20)
}

func TestVerticalListLayout_Arrange(t *testing.T) {
	is := is.New(t)

	eventManager := event.NewManager()
	vl := &VerticalListLayout{RowGap: 5}

	c := NewContainer(&ContainerOptions{Layout: vl})
	c.SetEventManager(eventManager)
	c.SetDimensions(200, 200)

	b1 := NewButton(&ButtonOptions{Width: option.Int(30), Height: option.Int(20)})
	b2 := NewButton(&ButtonOptions{Width: option.Int(30), Height: option.Int(30)})
	c.AddComponents(b1, b2)

	_, y1 := b1.Position()
	_, y2 := b2.Position()

	is.True(y1 >= 0)
	is.True(y2 > y1)
}

func TestVerticalListLayout_RowGap(t *testing.T) {
	is := is.New(t)

	eventManager := event.NewManager()
	gap := 10
	vl := &VerticalListLayout{RowGap: gap}

	c := NewContainer(&ContainerOptions{Layout: vl})
	c.SetEventManager(eventManager)
	c.SetDimensions(200, 200)

	b1 := NewButton(&ButtonOptions{Width: option.Int(30), Height: option.Int(20)})
	b2 := NewButton(&ButtonOptions{Width: option.Int(30), Height: option.Int(30)})
	c.AddComponents(b1, b2)

	_, y1 := b1.Position()
	_, y2 := b2.Position()

	expectedY2 := y1 + float64(b1.HeightWithPadding()) + float64(gap)
	is.Equal(y2, expectedY2)
}

func TestVerticalListLayout_HorizontalAlignment_Left(t *testing.T) {
	is := is.New(t)

	eventManager := event.NewManager()
	vl := &VerticalListLayout{RowGap: 0}

	c := NewContainer(&ContainerOptions{Layout: vl})
	c.SetEventManager(eventManager)
	c.SetDimensions(200, 200)

	b := NewButton(&ButtonOptions{
		Width:               option.Int(30),
		Height:              option.Int(20),
		HorizontalAlignment: option.AlignmentLeft,
	})
	c.AddComponent(b)

	x, _ := b.Position()
	is.Equal(x, float64(c.padding.Left))
}

func TestVerticalListLayout_HorizontalAlignment_Right(t *testing.T) {
	is := is.New(t)

	eventManager := event.NewManager()
	vl := &VerticalListLayout{RowGap: 0}

	c := NewContainer(&ContainerOptions{Layout: vl})
	c.SetEventManager(eventManager)
	c.SetDimensions(200, 200)

	wideButton := NewButton(&ButtonOptions{Width: option.Int(100), Height: option.Int(20)})
	rightButton := NewButton(&ButtonOptions{
		Width:               option.Int(30),
		Height:              option.Int(20),
		HorizontalAlignment: option.AlignmentRight,
	})
	c.AddComponents(wideButton, rightButton)

	xWide, _ := wideButton.Position()
	xRight, _ := rightButton.Position()

	is.True(xRight > xWide)
}

func TestVerticalListLayout_HorizontalAlignment_Centered(t *testing.T) {
	is := is.New(t)

	eventManager := event.NewManager()
	vl := &VerticalListLayout{RowGap: 0}

	c := NewContainer(&ContainerOptions{Layout: vl})
	c.SetEventManager(eventManager)
	c.SetDimensions(200, 200)

	wideButton := NewButton(&ButtonOptions{Width: option.Int(100), Height: option.Int(20)})
	centeredButton := NewButton(&ButtonOptions{
		Width:               option.Int(30),
		Height:              option.Int(20),
		HorizontalAlignment: option.AlignmentCenteredHorizontally,
	})
	c.AddComponents(wideButton, centeredButton)

	xWide, _ := wideButton.Position()
	xCentered, _ := centeredButton.Position()

	is.True(xCentered > xWide)
}

func TestVerticalListLayout_StretchedVertically(t *testing.T) {
	is := is.New(t)

	eventManager := event.NewManager()
	vl := &VerticalListLayout{RowGap: 0}

	c := NewContainer(&ContainerOptions{Layout: vl})
	c.SetEventManager(eventManager)
	c.SetDimensions(100, 100)

	fixedButton := NewButton(&ButtonOptions{
		Width:  option.Int(30),
		Height: option.Int(20),
	})
	stretchedButton := NewButton(&ButtonOptions{
		Width:             option.Int(30),
		Height:            option.Int(20),
		VerticalAlignment: option.StretchedVertically,
	})
	c.AddComponents(fixedButton, stretchedButton)

	c.SetDimensions(100, 100)
	vl.Rearrange(c)

	is.True(stretchedButton.Height() > 20)
}

func TestGridLayout_Arrange(t *testing.T) {
	is := is.New(t)

	eventManager := event.NewManager()
	gl := &GridLayout{Columns: 2, Rows: 2, ColumnGap: 5, RowGap: 5}
	gl.Setup()

	c := NewContainer(&ContainerOptions{Layout: gl})
	c.SetEventManager(eventManager)
	c.SetDimensions(200, 200)

	b1 := NewButton(&ButtonOptions{Width: option.Int(30), Height: option.Int(20)})
	b2 := NewButton(&ButtonOptions{Width: option.Int(40), Height: option.Int(25)})
	b3 := NewButton(&ButtonOptions{Width: option.Int(35), Height: option.Int(30)})
	b4 := NewButton(&ButtonOptions{Width: option.Int(25), Height: option.Int(15)})
	c.AddComponents(b1, b2, b3, b4)

	x1, y1 := b1.Position()
	x2, y2 := b2.Position()
	x3, y3 := b3.Position()
	x4, y4 := b4.Position()

	is.True(x1 >= 0)
	is.True(y1 >= 0)
	is.True(x2 > x1)
	is.True(y2 == y1)
	is.True(x3 == x1)
	is.True(y3 > y1)
	is.True(x4 > x3)
	is.True(y4 == y3)
}

func TestGridLayout_Overflow_HidesComponents(t *testing.T) {
	is := is.New(t)

	eventManager := event.NewManager()
	gl := &GridLayout{Columns: 2, Rows: 1, ColumnGap: 0, RowGap: 0}
	gl.Setup()

	c := NewContainer(&ContainerOptions{Layout: gl})
	c.SetEventManager(eventManager)
	c.SetDimensions(200, 200)

	b1 := NewButton(&ButtonOptions{Width: option.Int(30), Height: option.Int(20)})
	b2 := NewButton(&ButtonOptions{Width: option.Int(30), Height: option.Int(20)})
	b3 := NewButton(&ButtonOptions{Width: option.Int(30), Height: option.Int(20)})
	c.AddComponents(b1, b2, b3)

	is.Equal(b3.Hidden(), true)
}

func TestGridLayout_FixedColumnWidths(t *testing.T) {
	is := is.New(t)

	eventManager := event.NewManager()
	gl := &GridLayout{
		Columns:       2,
		Rows:          2,
		ColumnsWidths: []int{50, 80},
		ColumnGap:     0,
		RowGap:        0,
	}
	gl.Setup()

	c := NewContainer(&ContainerOptions{Layout: gl})
	c.SetEventManager(eventManager)
	c.SetDimensions(200, 200)

	narrowButton := NewButton(&ButtonOptions{Width: option.Int(10), Height: option.Int(20)})
	wideButton := NewButton(&ButtonOptions{Width: option.Int(100), Height: option.Int(20)})
	c.AddComponents(narrowButton, wideButton)

	_, y1 := narrowButton.Position()
	_, y2 := wideButton.Position()

	is.Equal(y1, y2)
}

func TestGridLayout_FixedRowHeights(t *testing.T) {
	is := is.New(t)

	eventManager := event.NewManager()
	gl := &GridLayout{
		Columns:       1,
		Rows:          2,
		ColumnsWidths: []int{50},
		RowsHeights:   []int{40, 60},
		ColumnGap:     0,
		RowGap:        0,
	}
	gl.Setup()

	c := NewContainer(&ContainerOptions{Layout: gl})
	c.SetEventManager(eventManager)
	c.SetDimensions(200, 200)

	b1 := NewButton(&ButtonOptions{Width: option.Int(30), Height: option.Int(10)})
	b2 := NewButton(&ButtonOptions{Width: option.Int(30), Height: option.Int(10)})
	c.AddComponents(b1, b2)

	_, y1 := b1.Position()
	_, y2 := b2.Position()

	is.True(y2 > y1)
}
