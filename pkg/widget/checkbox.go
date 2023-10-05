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
}

type CheckBoxOpt func(b *CheckBox)
type CheckBoxOptions struct {
	opts []CheckBoxOpt
}

type CheckBoxToggledEventArgs struct {
	CheckBox *CheckBox
}

type CheckBoxToggledHandlerFunc func(args *CheckBoxToggledEventArgs)

func NewCheckBox(posX, posY float64, options *CheckBoxOptions) *CheckBox {
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
	}

	cb.widget = cb.createWidget(posX, posY, width, height)

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
	label := NewLabel(13, 5, labelText, color, lblOpts.CenteredVertically())

	o.opts = append(o.opts, func(cb *CheckBox) {
		cb.SetLabel(label)
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
		for j := 0; j < cb.widget.pixelCols; j += 4 {
			if ((i == 0 || i == cb.lastPixelRowId) && j <= cb.lastPixelColId) || ((j == 0 || j == cb.lastPixelColId) && i <= cb.lastPixelRowId) {
				arr[j+cb.widget.pixelCols*i] = 230
				arr[j+1+cb.widget.pixelCols*i] = 230
				arr[j+2+cb.widget.pixelCols*i] = 230
			} else {
				arr[j+cb.widget.pixelCols*i] = cb.container.backgroundColor.R
				arr[j+1+cb.widget.pixelCols*i] = cb.container.backgroundColor.G
				arr[j+2+cb.widget.pixelCols*i] = cb.container.backgroundColor.B
			}
			arr[j+3+cb.widget.pixelCols*i] = 255
		}
	}

	return arr
}

func (cb *CheckBox) drawChecked() []byte {
	arr := make([]byte, cb.widget.pixelRows*cb.widget.pixelCols)

	for i := 0; i < cb.widget.pixelRows; i++ {
		for j := 0; j < cb.widget.pixelCols; j += 4 {
			if ((i == 0 || i == cb.lastPixelRowId) && j <= cb.lastPixelColId) || ((j == 0 || j == cb.lastPixelColId) && i <= cb.lastPixelRowId) ||
				(j > 4 && j < cb.penultimatePixelColId && i > 1 && i < cb.penultimatePixelRowId) {
				arr[j+cb.widget.pixelCols*i] = 230
				arr[j+1+cb.widget.pixelCols*i] = 230
				arr[j+2+cb.widget.pixelCols*i] = 230
			} else {
				arr[j+cb.widget.pixelCols*i] = cb.container.backgroundColor.R
				arr[j+1+cb.widget.pixelCols*i] = cb.container.backgroundColor.G
				arr[j+2+cb.widget.pixelCols*i] = cb.container.backgroundColor.B
			}
			arr[j+3+cb.widget.pixelCols*i] = 255
		}
	}

	return arr
}

func (cb *CheckBox) createWidget(posX, posY float64, width, height int) widget {
	widgetOptions := &WidgetOptions{}

	widgetOptions.MouseButtonReleasedHandler(func(args *WidgetMouseButtonReleasedEventArgs) {
		if !cb.disabled && args.Inside {
			cb.Checked = !cb.Checked
			cb.ToggledEvent.Fire(&CheckBoxToggledEventArgs{
				CheckBox: cb,
			})
		}
	})

	return *NewWidget(posX, posY, width, height, widgetOptions)
}
