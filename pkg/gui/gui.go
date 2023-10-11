package gui

import (
	"github.com/fglo/chopstiqs/pkg/component"
	"github.com/fglo/chopstiqs/pkg/event"
	"github.com/fglo/chopstiqs/pkg/input"
	ebiten "github.com/hajimehoshi/ebiten/v2"
)

// rootContainer is the gui main rootContainer that contains all other components
var rootContainer component.Container

// SetRootContainer sets the gui root container.
func SetRootContainer(container component.Container) {
	rootContainer = container
}

// Update updates containers.
// It should be called in the Ebiten Game's Update function.
func Update() {
	input.Update()
	rootContainer.FireEvents()
}

// Draw draws containers to the guiImage.
// It should be called in the Ebiten Game's Draw function.
func Draw(guiImage *ebiten.Image) {
	event.HandleFired()

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(rootContainer.Position())
	guiImage.DrawImage(rootContainer.Draw(), op)
}
