package chopstiqs

import (
	"github.com/fglo/chopstiqs/pkg/component"
	"github.com/fglo/chopstiqs/pkg/event"
	"github.com/fglo/chopstiqs/pkg/input"
	ebiten "github.com/hajimehoshi/ebiten/v2"
)

// GUI represents an entire graphical user interface.
// Only a single GUI instance should exist for each running program.
type GUI struct {
	containers []*component.Container
}

// New creates a new gui
func New() *GUI {
	return &GUI{
		containers: make([]*component.Container, 0),
	}
}

// AddContainer adds a container to the gui
func (gui *GUI) AddContainer(container *component.Container) {
	gui.containers = append(gui.containers, container)
}

// Update updates registered in the gui containers.
// It should be called in the Ebiten Game's Update function.
func (gui *GUI) Update() {
	input.Update()

	for _, container := range gui.containers {
		container.Update()
	}
}

// Draw draws registered in the gui containers to the guiImage
// It should be called in the Ebiten Game's Draw function.
func (gui *GUI) Draw(guiImage *ebiten.Image) {
	event.HandleFired()

	for _, container := range gui.containers {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(container.Position())
		guiImage.DrawImage(container.Draw(), op)
	}
}
