package gui

import (
	"encoding/xml"
	"fmt"
	"os"

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

func (gui *GUI) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(gui.rootContainer, start)
}

func (gui *GUI) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	if gui.rootContainer == nil {
		gui.rootContainer = &component.Container{}
	}

	return d.DecodeElement(gui.rootContainer, &start)
}

func (gui *GUI) SaveToFile(filepath string) error {
	serialized, err := xml.MarshalIndent(gui, " ", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize GUI: %w", err)
	}

	if err = os.WriteFile(filepath, serialized, 0644); err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	return nil
}

func (gui *GUI) LoadFromFile(filepath string) error {
	file, _ := os.ReadFile("unmarshal_test.xml")
	if err := xml.Unmarshal(file, gui); err != nil {
		return fmt.Errorf("failed to parse XML GUI file: %w", err)
	}

	return nil
}
