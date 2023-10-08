package component

import (
	"image/color"

	"github.com/fglo/chopstiqs/pkg/colorutils"
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

	pixelCols int
	pixelRows int

	lastPixelRowId        int
	penultimatePixelRowId int
	lastPixelColId        int
	penultimatePixelColId int

	color color.RGBA
}

type CheckBoxOptions struct {
	Color color.Color

	Label *Label
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

		pixelCols:             width * 4,
		lastPixelColId:        width*4 - 4,
		penultimatePixelColId: width*4 - 8,

		pixelRows:             height,
		lastPixelRowId:        height - 1,
		penultimatePixelRowId: height - 2,

		color: color.RGBA{230, 230, 230, 255},
	}

	cb.component = cb.createComponent(width, height)

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

func (cb *CheckBox) AddToggledHandler(f CheckBoxToggledHandlerFunc) *CheckBox {
	cb.ToggledEvent.AddHandler(func(args interface{}) { f(args.(*CheckBoxToggledEventArgs)) })

	return cb
}

// SetLabel sets the label of the checkbox and adjusts the checkbox's dimensions accordingly.
func (cb *CheckBox) SetLabel(label *Label) {
	cb.label = label
	cb.label.alignHorizontally = cb.label.alignToLeft
	cb.label.alignVertically = cb.label.centerVertically

	cb.label.SetPosistion(13, 4.5)

	cb.SetDimensions(cb.label.width+13, 10)
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

	return cb.image
}

func (cb *CheckBox) drawUnchecked() []byte {
	arr := make([]byte, cb.component.pixelRows*cb.component.pixelCols)

	for i := 0; i < cb.component.pixelRows; i++ {
		rowNumber := cb.component.pixelCols * i

		for j := 0; j < cb.component.pixelCols; j += 4 {
			if ((i == 0 || i == cb.lastPixelRowId) && j <= cb.lastPixelColId) || ((j == 0 || j == cb.lastPixelColId) && i <= cb.lastPixelRowId) {
				arr[j+rowNumber] = cb.color.R
				arr[j+1+rowNumber] = cb.color.G
				arr[j+2+rowNumber] = cb.color.B
				arr[j+3+rowNumber] = cb.color.A
			} else {
				arr[j+rowNumber] = cb.container.backgroundColor.R
				arr[j+1+rowNumber] = cb.container.backgroundColor.G
				arr[j+2+rowNumber] = cb.container.backgroundColor.B
				arr[j+3+rowNumber] = cb.container.backgroundColor.A
			}
		}
	}

	return arr
}

func (cb *CheckBox) drawChecked() []byte {
	arr := make([]byte, cb.component.pixelRows*cb.component.pixelCols)

	for i := 0; i < cb.component.pixelRows; i++ {
		rowNumber := cb.component.pixelCols * i

		for j := 0; j < cb.component.pixelCols; j += 4 {
			if ((i == 0 || i == cb.lastPixelRowId) && j <= cb.lastPixelColId) || ((j == 0 || j == cb.lastPixelColId) && i <= cb.lastPixelRowId) ||
				(j > 4 && j < cb.penultimatePixelColId && i > 1 && i < cb.penultimatePixelRowId) {
				arr[j+rowNumber] = cb.color.R
				arr[j+1+rowNumber] = cb.color.G
				arr[j+2+rowNumber] = cb.color.B
				arr[j+3+rowNumber] = cb.color.A
			} else {
				arr[j+rowNumber] = cb.container.backgroundColor.R
				arr[j+1+rowNumber] = cb.container.backgroundColor.G
				arr[j+2+rowNumber] = cb.container.backgroundColor.B
				arr[j+3+rowNumber] = cb.container.backgroundColor.A
			}
		}
	}

	return arr
}

func (cb *CheckBox) createComponent(width, height int) component {
	componentOptions := &ComponentOptions{}

	componentOptions.AddMouseButtonReleasedHandler(func(args *ComponentMouseButtonReleasedEventArgs) {
		if !cb.disabled && args.Inside {
			cb.Checked = !cb.Checked
			cb.ToggledEvent.Fire(&CheckBoxToggledEventArgs{
				CheckBox: cb,
			})
		}
	})

	return *NewComponent(width, height, componentOptions)
}
