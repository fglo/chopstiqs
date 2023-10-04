package chopstiqs

import (
	"github.com/fglo/chopstiqs/pkg/event"
	"github.com/fglo/chopstiqs/pkg/input"
	"github.com/fglo/chopstiqs/pkg/widget"
	ebiten "github.com/hajimehoshi/ebiten/v2"
)

type Gui struct {
	containers []*widget.Container
}

func New() *Gui {
	return &Gui{
		containers: make([]*widget.Container, 0),
	}
}

func (gui *Gui) AddContainer(container *widget.Container) {
	gui.containers = append(gui.containers, container)
}

func (gui *Gui) Update(mouse *input.Mouse) {
	for _, component := range gui.containers {
		component.Update(mouse)
	}
}

func (gui *Gui) Draw(guiImage *ebiten.Image, mouse *input.Mouse) {
	event.ExecuteDeferred()

	for _, component := range gui.containers {
		component.Update(mouse)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(component.Position())
		guiImage.DrawImage(component.Draw(mouse), op)
	}
}
