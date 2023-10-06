package main

import (
	"errors"
	"image/color"
	"log"

	"github.com/fglo/chopstiqs"
	"github.com/fglo/chopstiqs/pkg/component"
	"github.com/fglo/chopstiqs/pkg/input"
	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var Terminated = errors.New("terminated")

func main() {
	if err := ebiten.RunGame(NewGame()); err != nil {
		if err == Terminated {
			return
		}

		log.Fatal(err)
	}
}

// Game encapsulates game logic
type Game struct {
	gui   *chopstiqs.GUI
	mouse *input.Mouse

	bgColorToggled bool

	screenWidth  int
	screenHeight int

	quitIsPressed bool
}

// New generates a new Game object.
func NewGame() *Game {
	g := &Game{
		gui:          &chopstiqs.GUI{},
		mouse:        input.NewMouse(),
		screenWidth:  200,
		screenHeight: 200,
	}

	ebiten.SetWindowSize(g.getWindowSize())
	ebiten.SetWindowTitle("chopstiqs demo")

	container := component.NewContainer(0, 0, 200, 200, color.RGBA{32, 32, 32, 255})

	lbl := component.NewLabel("chopstiqs demo", &component.LabelOptions{Color: color.RGBA{120, 190, 100, 255}, VerticalAlignment: component.AlignmentTop})
	container.AddComponent(5, 5, lbl)

	btnOpts := &component.ButtonOptions{
		Label: component.NewLabel("toggle background", &component.LabelOptions{Color: color.RGBA{25, 25, 25, 255}}),
	}

	btnOpts.AddClickedHandler(func(args *component.ButtonClickedEventArgs) {
		if !g.bgColorToggled {
			container.SetBackgroundColor(color.RGBA{9, 32, 42, 255})
		} else {
			container.SetBackgroundColor(color.RGBA{32, 32, 32, 255})
		}

		g.bgColorToggled = !g.bgColorToggled
	})

	btn := component.NewButton(btnOpts)
	container.AddComponent(5, 35, btn)

	btn2 := component.NewButton(&component.ButtonOptions{
		Color:         color.RGBA{100, 180, 90, 255},
		ColorPressed:  color.RGBA{90, 160, 80, 255},
		ColorHovered:  color.RGBA{120, 190, 100, 255},
		ColorDisabled: color.RGBA{80, 100, 70, 255},
	})
	container.AddComponent(5, 55, btn2)

	cbOpts := &component.CheckBoxOptions{
		Color: color.RGBA{255, 100, 50, 255},
	}
	cb := component.NewCheckBox(cbOpts)
	cb.Toggle()
	container.AddComponent(5, 20, cb)

	cb2Opts := &component.CheckBoxOptions{
		Color: color.RGBA{230, 230, 230, 255},
		Label: component.NewLabel("disable buttons", &component.LabelOptions{Color: color.RGBA{230, 230, 230, 255}}),
	}

	cb2Opts.AddToggledHandler(func(args *component.CheckBoxToggledEventArgs) {
		btn.SetDisabled(args.CheckBox.Checked)
		btn2.SetDisabled(args.CheckBox.Checked)
	})

	cb2 := component.NewCheckBox(cb2Opts)
	container.AddComponent(20, 20, cb2)

	g.gui.AddContainer(container)

	return g
}

func (g *Game) getWindowSize() (int, int) {
	var factor float32 = 1.8
	return int(float32(g.screenWidth) * factor), int(float32(g.screenHeight) * factor)
}

// Layout implements ebiten.Game's Layout.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.screenWidth, g.screenHeight
}

// Update updates the current game state.
func (g *Game) Update() error {
	g.mouse.Update()
	g.gui.Update(g.mouse)

	return g.checkQuitButton()
}

func (g *Game) checkQuitButton() error {
	if !g.quitIsPressed && inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		g.quitIsPressed = true
	}
	if g.quitIsPressed && inpututil.IsKeyJustReleased(ebiten.KeyQ) {
		g.quitIsPressed = false
		return Terminated
	}
	return nil
}

// Draw draws the current game to the given screen.
func (g *Game) Draw(screen *ebiten.Image) {
	g.gui.Draw(screen, g.mouse)
}
