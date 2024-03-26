package gui

import (
	"github.com/fglo/chopstiqs/pkg/component"
	"github.com/fglo/chopstiqs/pkg/event"
	"github.com/fglo/chopstiqs/pkg/input"
	ebiten "github.com/hajimehoshi/ebiten/v2"
)

type GUI struct {
	// rootContainer is the gui main rootContainer that contains all other components
	rootContainer *component.Container
	// eventManager is a queue of events by GUI components
	eventManager *event.Manager

	focusedComponent component.Component
}

func New() *GUI {
	return &GUI{
		eventManager: event.NewManager(),
	}
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
