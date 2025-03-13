package chopstiqs

import (
	"image"

	"github.com/fglo/chopstiqs/component"
	"github.com/fglo/chopstiqs/event"
	"github.com/fglo/chopstiqs/input"
	"github.com/fglo/chopstiqs/option"
	ebiten "github.com/hajimehoshi/ebiten/v2"
)

type GUI struct {
	// rootContainer is the gui main rootContainer that contains all other components
	rootContainer *component.Container
	// eventManager is a queue of events by GUI components
	eventManager *event.Manager

	focusedComponent component.Component

	horizontalAlignment option.HorizontalAlignment
	verticalAlignment   option.VerticalAlignment
}

type GUIOptions struct {
	HorizontalAlignment option.HorizontalAlignment
	VerticalAlignment   option.VerticalAlignment
}

func NewGUI(opt *GUIOptions) *GUI {
	gui := &GUI{
		eventManager: event.NewManager(),
	}

	if opt != nil {
		gui.horizontalAlignment = opt.HorizontalAlignment
		gui.verticalAlignment = opt.VerticalAlignment
	}

	return gui
}

// SetRootContainer sets the gui root container.
func (gui *GUI) SetRootContainer(container *component.Container) {
	container.SetEventManager(gui.eventManager)
	gui.rootContainer = container
	gui.rootContainer.AddFocusedHandler(gui.handleFocusEvent)
}

// Update updates containers.
// It should be called in the Ebiten Game's Update function.
func (gui *GUI) Update() {
	input.Update()
	gui.rootContainer.FireEvents()
}

// Draw draws containers to the guiImage.
// It should be called in the Ebiten Game's Draw function.
func (gui *GUI) Draw(guiImage *ebiten.Image) {
	input.Draw()
	defer input.AfterDraw()

	gui.eventManager.HandleFired()

	gui.alignRootContainerInBounds(guiImage.Bounds())

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(gui.rootContainer.Position())
	guiImage.DrawImage(gui.rootContainer.Draw(), op)
}

func (gui *GUI) NewContainer(options *component.ContainerOptions) *component.Container {
	c := component.NewContainer(options)
	c.SetEventManager(gui.eventManager)
	return c
}

func (gui *GUI) NewButton(options *component.ButtonOptions) *component.Button {
	b := component.NewButton(options)
	b.SetEventManager(gui.eventManager)
	return b
}

func (gui *GUI) NewCheckBox(options *component.CheckBoxOptions) *component.CheckBox {
	cb := component.NewCheckBox(options)
	cb.SetEventManager(gui.eventManager)
	return cb
}

func (gui *GUI) NewLabel(labelText string, options *component.LabelOptions) *component.Label {
	l := component.NewLabel(labelText, options)
	l.SetEventManager(gui.eventManager)
	return l
}

func (gui *GUI) NewSlider(options *component.SliderOptions) *component.Slider {
	s := component.NewSlider(options)
	s.SetEventManager(gui.eventManager)
	return s
}

func (gui *GUI) FocusedComponent() component.Component {
	return gui.focusedComponent
}

func (gui *GUI) handleFocusEvent(args *component.ComponentFocusedEventArgs) {
	if args.Focused {
		if gui.focusedComponent != nil && gui.focusedComponent != args.Component {
			gui.focusedComponent.SetFocused(false)
		}

		gui.focusedComponent = args.Component
	} else if gui.focusedComponent == args.Component {
		gui.focusedComponent = nil
	}
}

func (gui *GUI) alignRootContainerInBounds(bounds image.Rectangle) {
	w := gui.rootContainer.WidthWithPadding()
	h := gui.rootContainer.HeightWithPadding()

	var xTranslation int
	switch gui.horizontalAlignment {
	case option.AlignmentLeft:
		xTranslation = 0
	case option.AlignmentCenteredHorizontally:
		xTranslation = (bounds.Dx() - w) / 2
	case option.AlignmentRight:
		xTranslation = bounds.Dx() - w
	}

	var yTranslation int
	switch gui.verticalAlignment {
	case option.AlignmentTop:
		yTranslation = 0
	case option.AlignmentCenteredVertically:
		yTranslation = (bounds.Dy() - h) / 2
	case option.AlignmentBottom:
		yTranslation = bounds.Dy() - h
	}

	gui.rootContainer.SetPosition(float64(xTranslation), float64(yTranslation))
}
