package widget

import (
	"image/color"

	"github.com/fglo/chopstiqs/pkg/event"
	ebiten "github.com/hajimehoshi/ebiten/v2"
)

type CheckBox struct {
	widget
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

type CheckBoxOpt func(b *CheckBox)
type CheckBoxOptions struct {
	opts []CheckBoxOpt
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

	cb.widget = cb.createWidget(width, height)

	if options != nil {
		for _, o := range options.opts {
			o(cb)
		}
	}

	return cb
}

func (o *CheckBoxOptions) ToggledHandler(f CheckBoxToggledHandlerFunc) *CheckBoxOptions {
	o.opts = append(o.opts, func(cb *CheckBox) {
		cb.ToggledEvent.AddHandler(func(args interface{}) {
			f(args.(*CheckBoxToggledEventArgs))
		})
	})

	return o
}

func (o *CheckBoxOptions) Text(labelText string, color color.RGBA) *CheckBoxOptions {
	lblOpts := &LabelOptions{}
	label := NewLabel(labelText, lblOpts.Color(color).CenteredVertically())
	label.SetPosistion(13, 5)

	o.opts = append(o.opts, func(cb *CheckBox) {
		cb.SetLabel(label)
	})

	return o
}

func (o *CheckBoxOptions) Color(color color.RGBA) *CheckBoxOptions {
	o.opts = append(o.opts, func(cb *CheckBox) {
		cb.color = color
	})

	return o
}

func (cb *CheckBox) SetLabel(label *Label) {
	cb.label = label
	cb.SetWidth(label.width + 13)
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
	arr := make([]byte, cb.widget.pixelRows*cb.widget.pixelCols)

	for i := 0; i < cb.widget.pixelRows; i++ {
		rowNumber := cb.widget.pixelCols * i

		for j := 0; j < cb.widget.pixelCols; j += 4 {
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
	arr := make([]byte, cb.widget.pixelRows*cb.widget.pixelCols)

	for i := 0; i < cb.widget.pixelRows; i++ {
		rowNumber := cb.widget.pixelCols * i

		for j := 0; j < cb.widget.pixelCols; j += 4 {
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

func (cb *CheckBox) createWidget(width, height int) widget {
	widgetOptions := &WidgetOptions{}

	widgetOptions.MouseButtonReleasedHandler(func(args *WidgetMouseButtonReleasedEventArgs) {
		if !cb.disabled && args.Inside {
			cb.Checked = !cb.Checked
			cb.ToggledEvent.Fire(&CheckBoxToggledEventArgs{
				CheckBox: cb,
			})
		}
	})

	return *NewWidget(width, height, widgetOptions)
}
