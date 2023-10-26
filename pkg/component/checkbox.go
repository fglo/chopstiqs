package component

import (
	"image/color"

	"github.com/fglo/chopstiqs/internal/colorutils"
	"github.com/fglo/chopstiqs/pkg/event"
	ebiten "github.com/hajimehoshi/ebiten/v2"
)

type CheckBox struct {
	component
	Checked bool

	ToggledEvent *event.Event

	label *Label

	width  int
	height int

	firstPixelRowId       int
	secondPixelRowId      int
	lastPixelRowId        int
	penultimatePixelRowId int

	firstPixelColId       int
	secondPixelColId      int
	lastPixelColId        int
	penultimatePixelColId int

	color color.RGBA
}

type CheckBoxOptions struct {
	Color color.Color

	Label *Label

	LeftPadding   *int
	RightPadding  *int
	TopPadding    *int
	BottomPadding *int
}

type CheckBoxToggledEventArgs struct {
	CheckBox *CheckBox
}

type CheckBoxToggledHandlerFunc func(args *CheckBoxToggledEventArgs)

func NewCheckBox(options *CheckBoxOptions) *CheckBox {
	width := 10
	height := 10

	cb := &CheckBox{
		Checked:      false,
		ToggledEvent: &event.Event{},

		width:  width,
		height: height,

		color: color.RGBA{230, 230, 230, 255},
	}

	cb.component = cb.createComponent(width, height, options)

	cb.firstPixelColId = cb.leftPadding * 4
	cb.secondPixelColId = cb.firstPixelColId + 4

	cb.lastPixelColId = (width+cb.leftPadding)*4 - 4
	cb.penultimatePixelColId = cb.lastPixelColId - 4

	cb.firstPixelRowId = cb.topPadding
	cb.secondPixelRowId = cb.firstPixelRowId + 1

	cb.lastPixelRowId = height + cb.topPadding - 1
	cb.penultimatePixelRowId = cb.lastPixelRowId - 1

	if options != nil {
		if options.Label != nil {
			cb.SetLabel(options.Label)
		}

		if options.Color != nil {
			cb.color = colorutils.ColorToRGBA(options.Color)
		}
	}

	return cb
}

func (cb *CheckBox) createComponent(width, height int, options *CheckBoxOptions) component {
	var componentOptions ComponentOptions

	if options != nil {
		componentOptions = ComponentOptions{
			LeftPadding:   options.LeftPadding,
			RightPadding:  options.RightPadding,
			TopPadding:    options.TopPadding,
			BottomPadding: options.BottomPadding,
		}
	}

	component := NewComponent(width, height, &componentOptions)

	component.AddMouseButtonReleasedHandler(func(args *ComponentMouseButtonReleasedEventArgs) {
		if !cb.disabled && args.Inside {
			cb.Checked = !cb.Checked
			cb.ToggledEvent.Fire(&CheckBoxToggledEventArgs{
				CheckBox: cb,
			})
		}
	})

	return *component
}

func (cb *CheckBox) AddToggledHandler(f CheckBoxToggledHandlerFunc) *CheckBox {
	cb.ToggledEvent.AddHandler(func(args interface{}) { f(args.(*CheckBoxToggledEventArgs)) })

	return cb
}

// SetLabel sets the label of the checkbox and adjusts the checkbox's dimensions accordingly.
func (cb *CheckBox) SetLabel(label *Label) {
	width := 10
	height := 10

	cb.label = label
	cb.label.alignHorizontally = cb.label.alignToLeft
	cb.label.alignVertically = cb.label.centerVertically

	if cb.label.leftPadding == 0 {
		cb.label.leftPadding = 3
	}

	cb.label.SetPosistion(10+float64(cb.label.leftPadding), 4.5)

	cb.SetDimensions(width+cb.label.widthWithPadding, height)

	cb.firstPixelColId = cb.leftPadding * 4
	cb.secondPixelColId = cb.firstPixelColId + 4

	cb.lastPixelColId = (width+cb.leftPadding)*4 - 4
	cb.penultimatePixelColId = cb.lastPixelColId - 4

	cb.firstPixelRowId = cb.topPadding
	cb.secondPixelRowId = cb.firstPixelRowId + 1

	cb.lastPixelRowId = height + cb.topPadding - 1
	cb.penultimatePixelRowId = cb.lastPixelRowId - 1
}

func (cb *CheckBox) Set(checked bool) {
	cb.Checked = checked
}

func (cb *CheckBox) Toggle() {
	cb.Checked = !cb.Checked
	cb.ToggledEvent.Fire(&CheckBoxToggledEventArgs{
		CheckBox: cb,
	})
}

func (cb *CheckBox) Draw() *ebiten.Image {
	if cb.hidden {
		return cb.image
	}

	if cb.Checked {
		cb.image.WritePixels(cb.drawChecked())
	} else {
		cb.image.WritePixels(cb.drawUnchecked())
	}

	if cb.label != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(cb.label.Position())
		cb.image.DrawImage(cb.label.Draw(), op)
	}

	cb.component.Draw()

	return cb.image
}

func (cb *CheckBox) isBorder(rowId, colId int) bool {
	return rowId == cb.firstPixelRowId || rowId == cb.lastPixelRowId || colId == cb.firstPixelColId || colId == cb.lastPixelColId
}

func (cb *CheckBox) isColored(rowId, colId int) bool {
	return colId > cb.secondPixelColId && colId < cb.penultimatePixelColId && rowId > cb.secondPixelRowId && rowId < cb.penultimatePixelRowId
}

func (cb *CheckBox) drawUnchecked() []byte {
	arr := make([]byte, cb.component.pixelRows*cb.component.pixelCols)
	backgroundColor := cb.container.GetBackgroundColor()

	for rowId := cb.firstPixelRowId; rowId <= cb.lastPixelRowId; rowId++ {
		rowNumber := cb.component.pixelCols * rowId

		for colId := cb.firstPixelColId; colId <= cb.lastPixelColId; colId += 4 {
			if cb.isBorder(rowId, colId) {
				arr[colId+rowNumber] = cb.color.R
				arr[colId+1+rowNumber] = cb.color.G
				arr[colId+2+rowNumber] = cb.color.B
				arr[colId+3+rowNumber] = cb.color.A
			} else {
				arr[colId+rowNumber] = backgroundColor.R
				arr[colId+1+rowNumber] = backgroundColor.G
				arr[colId+2+rowNumber] = backgroundColor.B
				arr[colId+3+rowNumber] = backgroundColor.A
			}
		}
	}

	return arr
}

func (cb *CheckBox) drawChecked() []byte {
	arr := make([]byte, cb.component.pixelRows*cb.component.pixelCols)
	backgroundColor := cb.container.GetBackgroundColor()

	for rowId := cb.firstPixelRowId; rowId <= cb.lastPixelRowId; rowId++ {
		rowNumber := cb.component.pixelCols * rowId

		for colId := cb.firstPixelColId; colId <= cb.lastPixelColId; colId += 4 {
			if cb.isBorder(rowId, colId) || cb.isColored(rowId, colId) {
				arr[colId+rowNumber] = cb.color.R
				arr[colId+1+rowNumber] = cb.color.G
				arr[colId+2+rowNumber] = cb.color.B
				arr[colId+3+rowNumber] = cb.color.A
			} else {
				arr[colId+rowNumber] = backgroundColor.R
				arr[colId+1+rowNumber] = backgroundColor.G
				arr[colId+2+rowNumber] = backgroundColor.B
				arr[colId+3+rowNumber] = backgroundColor.A
			}
		}
	}

	return arr
}
