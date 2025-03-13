package component

import (
	"image/color"

	"github.com/fglo/chopstiqs/event"
	"github.com/fglo/chopstiqs/option"
	ebiten "github.com/hajimehoshi/ebiten/v2"
)

type CheckBox struct {
	component
	checked bool

	ToggledEvent *event.Event

	label *Label

	cbWidth  int
	cbHeight int

	firstPixelRowId       int
	secondPixelRowId      int
	lastPixelRowId        int
	penultimatePixelRowId int

	firstPixelColId       int
	secondPixelColId      int
	lastPixelColId        int
	penultimatePixelColId int

	drawer CheckBoxDrawer
}

type CheckBoxOptions struct {
	Width  option.OptInt
	Height option.OptInt

	Label *Label

	Padding *Padding

	Drawer CheckBoxDrawer
}

type CheckBoxToggledEventArgs struct {
	CheckBox *CheckBox
}

type CheckBoxToggledHandlerFunc func(args *CheckBoxToggledEventArgs)

func NewCheckBox(opt *CheckBoxOptions) *CheckBox {

	cb := &CheckBox{
		checked:      false,
		ToggledEvent: &event.Event{},

		cbWidth:  10,
		cbHeight: 10,

		drawer: DefaultCheckBoxDrawer{
			Color: color.RGBA{230, 230, 230, 255},
		},
	}

	width := cb.cbWidth
	height := cb.cbHeight

	cb.SetDimensions(width, height)

	if opt != nil {
		if opt.Width.IsSet() && opt.Height.IsSet() {
			cb.cbWidth = opt.Width.Val()
			width = cb.cbWidth
		}

		if opt.Height.IsSet() {
			cb.cbHeight = opt.Height.Val()
			height = cb.cbHeight
		}

		cb.SetDimensions(width, height)

		if opt.Label != nil {
			cb.SetLabel(opt.Label)
		}

		if opt.Drawer != nil {
			cb.drawer = opt.Drawer
		}
	}

	cb.setUpComponent(opt)

	cb.setDrawingDimensions()

	return cb
}

func (cb *CheckBox) setUpComponent(opt *CheckBoxOptions) {
	var componentOptions ComponentOptions

	if opt != nil {
		componentOptions = ComponentOptions{
			Padding: opt.Padding,
		}
	}

	cb.component.setUpComponent(&componentOptions)

	cb.component.AddMouseButtonReleasedHandler(func(args *ComponentMouseButtonReleasedEventArgs) {
		if !cb.disabled && args.Inside {
			cb.checked = !cb.checked
			cb.eventManager.Fire(cb.ToggledEvent, &CheckBoxToggledEventArgs{
				CheckBox: cb,
			})
		}
	})
}

func (cb *CheckBox) setDrawingDimensions() {
	cb.firstPixelColId = cb.padding.Left * 4
	cb.secondPixelColId = cb.firstPixelColId + 4

	cb.lastPixelColId = (cb.cbWidth+cb.padding.Left)*4 - 4
	cb.penultimatePixelColId = cb.lastPixelColId - 4

	cb.firstPixelRowId = cb.padding.Top
	cb.secondPixelRowId = cb.firstPixelRowId + 1

	cb.lastPixelRowId = cb.cbHeight + cb.padding.Top - 1
	cb.penultimatePixelRowId = cb.lastPixelRowId - 1
}

func (cb *CheckBox) AddToggledHandler(f CheckBoxToggledHandlerFunc) *CheckBox {
	cb.ToggledEvent.AddHandler(func(args interface{}) { f(args.(*CheckBoxToggledEventArgs)) })

	return cb
}

// SetLabel sets the label of the checkbox and adjusts the checkbox's dimensions accordingly.
func (cb *CheckBox) SetLabel(label *Label) {
	label.setContainer(cb)
	cb.label = label
	cb.label.horizontalAlignment = option.AlignmentRight
	cb.label.verticalAlignment = option.AlignmentCenteredVertically

	if cb.label.padding.Left == 0 {
		cb.label.SetPaddingLeft(2)
	}

	width := cb.width
	if width <= cb.cbWidth+cb.label.widthWithPadding {
		width = cb.cbWidth + cb.label.widthWithPadding
	}

	height := cb.height // cb.cbHeight
	if height <= cb.label.height {
		height = cb.label.height
	}

	cb.SetDimensions(width, height)

	cb.label.align()
}

func (cb *CheckBox) Set(checked bool) {
	prevState := cb.checked
	cb.checked = checked
	if prevState != cb.checked {
		cb.eventManager.Fire(cb.ToggledEvent, &CheckBoxToggledEventArgs{
			CheckBox: cb,
		})
	}
}

func (cb *CheckBox) Checked() bool {
	return cb.checked
}

func (cb *CheckBox) Toggle() {
	cb.checked = !cb.checked
	cb.eventManager.Fire(cb.ToggledEvent, &CheckBoxToggledEventArgs{
		CheckBox: cb,
	})
}

func (cb *CheckBox) SetPosition(posX, posY float64) {
	cb.component.SetPosition(posX, posY)
	if cb.label != nil {
		cb.label.RecalculateAbsPosition()
	}
}

func (cb *CheckBox) RecalculateAbsPosition() {
	cb.component.RecalculateAbsPosition()
	if cb.label != nil {
		cb.label.RecalculateAbsPosition()
	}
}

func (cb *CheckBox) SetBackgroundColor(color color.RGBA) {
	cb.container.SetBackgroundColor(color)
}

func (cb *CheckBox) GetBackgroundColor() color.RGBA {
	return cb.container.GetBackgroundColor()
}

func (cb *CheckBox) FireEvents() {
	if cb.label != nil {
		cb.label.FireEvents()
	}

	cb.component.FireEvents()
}

func (cb *CheckBox) Draw() *ebiten.Image {
	if cb.hidden {
		return cb.image
	}

	cb.drawer.Draw(cb)

	if cb.label != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(cb.label.Position())
		cb.image.DrawImage(cb.label.Draw(), op)
	}

	cb.component.Draw()

	return cb.image
}
