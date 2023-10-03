package chopstiqs

import (
	"github.com/fglo/chopstiqs/pkg/event"
	"github.com/fglo/chopstiqs/pkg/input"
	"github.com/fglo/chopstiqs/pkg/widget"
	ebiten "github.com/hajimehoshi/ebiten/v2"
)

type Gui struct {
	components []widget.Widget
}

func New() *Gui {
	return &Gui{
		components: make([]widget.Widget, 0),
	}
}

func (gui *Gui) AddComponent(component widget.Widget) {
	gui.components = append(gui.components, component)
}

func (gui *Gui) Update(mouse *input.Mouse) {
	for _, component := range gui.components {
		component.FireEvents(mouse)
	}
}

func (gui *Gui) Draw(guiImage *ebiten.Image, mouse *input.Mouse) {
	event.ExecuteDeferred()

	for _, component := range gui.components {
		component.FireEvents(mouse)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(component.Position())
		guiImage.DrawImage(component.Draw(), op)
	}
}
