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

	cb.component.width = cb.cbWidth
	cb.component.height = cb.cbHeight

	if opt != nil {
		if opt.Width.IsSet() {
			cb.cbWidth = opt.Width.Val()
			cb.component.width = cb.cbWidth
		}

		if opt.Height.IsSet() {
			cb.cbHeight = opt.Height.Val()
			cb.component.height = cb.cbHeight
		}

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
	cb.label = label
	cb.label.alignHorizontally = cb.label.alignToLeft
	cb.label.alignVertically = cb.label.centerVertically

	if cb.label.padding.Left == 0 {
		cb.label.SetPaddingLeft(3)
	}

	cb.label.SetPosistion(float64(cb.cbWidth), float64(cb.cbHeight)/2)

	cb.SetDimensions(cb.cbWidth+cb.label.widthWithPadding, cb.cbHeight)
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
