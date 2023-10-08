package gui

import (
	"github.com/fglo/chopstiqs/pkg/component"
	"github.com/fglo/chopstiqs/pkg/event"
	"github.com/fglo/chopstiqs/pkg/input"
	ebiten "github.com/hajimehoshi/ebiten/v2"
)

// mainContainer is the gui main mainContainer that contains all other components
var mainContainer component.Container

func SetMainContainer(container component.Container) {
	mainContainer = container
}

// Update updates containers.
// It should be called in the Ebiten Game's Update function.
func Update() {
	input.Update()
	mainContainer.FireEvents()
}

// Draw draws containers to the guiImage.
// It should be called in the Ebiten Game's Draw function.
func Draw(guiImage *ebiten.Image) {
	event.HandleFired()

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(mainContainer.Position())
	guiImage.DrawImage(mainContainer.Draw(), op)
}
