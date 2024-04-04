package main

import (
	"errors"
	"image/color"
	"log"

	"github.com/fglo/chopstiqs/component"
	"github.com/fglo/chopstiqs/gui"
	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Game encapsulates game logic
type Game struct {
	gui            *gui.GUI
	bgColorToggled bool
	screenWidth    int
	screenHeight   int
}

// New generates a new Game object.
func NewGame() *Game {
	g := &Game{
		gui:          gui.New(),
		screenWidth:  130,
		screenHeight: 70,
	}

	ebiten.SetWindowSize(int(float64(g.screenWidth)*1.75), int(float64(g.screenHeight)*1.75))
	ebiten.SetWindowTitle("chopstiqs demo")

	rootContainer := g.gui.NewContainer(&component.ContainerOptions{
		Padding: &component.Padding{Left: 5, Right: 5, Top: 5, Bottom: 5},
		Layout:  &component.VerticalListLayout{RowGap: 5}})

	g.gui.SetRootContainer(rootContainer)

	lblTitle := g.gui.NewLabel("chopstiqs demo", &component.LabelOptions{Color: color.RGBA{120, 190, 100, 255}})

	btn := g.gui.NewButton(&component.ButtonOptions{
		Label: g.gui.NewLabel("toggle background", &component.LabelOptions{Color: color.RGBA{50, 50, 50, 255}}),
	})

	btn.AddClickedHandler(func(args *component.ButtonClickedEventArgs) { g.bgColorToggled = !g.bgColorToggled })

	cbOpts := &component.CheckBoxOptions{
		Label: g.gui.NewLabel("disable components", &component.LabelOptions{Color: color.RGBA{230, 230, 230, 255}}),
	}
	cb := g.gui.NewCheckBox(cbOpts)
	cb.AddToggledHandler(func(args *component.CheckBoxToggledEventArgs) {
		btn.SetDisabled(args.CheckBox.Checked())
	})

	rootContainer.AddComponent(lblTitle)
	rootContainer.AddComponent(cb)
	rootContainer.AddComponent(btn)

	return g
}

// Layout implements ebiten.Game's Layout.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.screenWidth, g.screenHeight
}

// Update updates the current game state.
func (g *Game) Update() error {
	g.gui.Update()
	return g.checkQuitButton()
}

// Draw draws the current game to the given screen.
func (g *Game) Draw(screen *ebiten.Image) {
	if !g.bgColorToggled {
		screen.Fill(color.RGBA{32, 32, 32, 255})
	} else {
		screen.Fill(color.RGBA{10, 50, 50, 255})
	}
	g.gui.Draw(screen)
}

func (g *Game) checkQuitButton() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		return Terminated
	}

	return nil
}

var Terminated = errors.New("terminated")

func main() {
	if err := ebiten.RunGame(NewGame()); err != nil {
		if err == Terminated {
			return
		}

		log.Fatal(err)
	}
}
