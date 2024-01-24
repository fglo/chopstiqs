package component

import (
	"image/color"

	"github.com/fglo/chopstiqs/pkg/event"
	"github.com/fglo/chopstiqs/pkg/option"
	ebiten "github.com/hajimehoshi/ebiten/v2"
)

type CheckBox struct {
	component
	Checked bool

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

func NewCheckBox(options *CheckBoxOptions) *CheckBox {

	cb := &CheckBox{
		Checked:      false,
		ToggledEvent: &event.Event{},

		cbWidth:  10,
		cbHeight: 10,

		drawer: DefaultCheckBoxDrawer{
			Color: color.RGBA{230, 230, 230, 255},
		},
	}

	cb.component.width = cb.cbWidth
	cb.component.height = cb.cbHeight

	if options != nil {
		if options.Width.IsSet() {
			cb.cbWidth = options.Width.Val()
			cb.component.width = cb.cbWidth
		}

		if options.Height.IsSet() {
			cb.cbHeight = options.Height.Val()
			cb.component.height = cb.cbHeight
		}

		if options.Label != nil {
			cb.SetLabel(options.Label)
		}

		if options.Drawer != nil {
			cb.drawer = options.Drawer
		}
	}

	cb.setUpComponent(options)

	cb.setDrawingDimensions()

	return cb
}

func (cb *CheckBox) setUpComponent(options *CheckBoxOptions) {
	var componentOptions ComponentOptions

	if options != nil {
		componentOptions = ComponentOptions{
			Padding: options.Padding,
		}
	}

	cb.component.setUpComponent(&componentOptions)

	cb.component.AddMouseButtonReleasedHandler(func(args *ComponentMouseButtonReleasedEventArgs) {
		if !cb.disabled && args.Inside {
			cb.Checked = !cb.Checked
			cb.ToggledEvent.Fire(&CheckBoxToggledEventArgs{
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

	cb.drawer.Draw(cb)

	if cb.label != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(cb.label.Position())
		cb.image.DrawImage(cb.label.Draw(), op)
	}

	cb.component.Draw()

	return cb.image
}
