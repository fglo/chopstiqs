package main

import (
	"errors"
	"image/color"
	imgColor "image/color"
	"log"

	"github.com/fglo/chopstiqs"
	"github.com/fglo/chopstiqs/pkg/input"
	"github.com/fglo/chopstiqs/pkg/widget"
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
	gui   *chopstiqs.Gui
	mouse *input.Mouse

	bgColorToggled bool

	screenWidth  int
	screenHeight int

	quitIsPressed bool
}

// New generates a new Game object.
func NewGame() *Game {
	g := &Game{
		gui:          &chopstiqs.Gui{},
		mouse:        input.NewMouse(),
		screenWidth:  200,
		screenHeight: 200,
	}

	ebiten.SetWindowSize(g.getWindowSize())
	ebiten.SetWindowTitle("chopstiqs demo")

	container := widget.NewContainer(0, 0, 200, 200, imgColor.RGBA{9, 32, 42, 255})

	btnOpts := &widget.ButtonOptions{}
	btn := widget.NewButton(5, 20, btnOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
		if !g.bgColorToggled {
			container.SetBackgroundColor(imgColor.RGBA{32, 32, 32, 255})
		} else {
			container.SetBackgroundColor(imgColor.RGBA{9, 32, 42, 255})
		}

		g.bgColorToggled = !g.bgColorToggled
	}).Text("toggle background", color.RGBA{25, 25, 25, 255}))
	container.AddComponent(btn)

	cb := widget.NewCheckBox(5, 5, nil)
	cb.Toggle()
	container.AddComponent(cb)

	cb2Opts := &widget.CheckBoxOptions{}
	cb2Opts = cb2Opts.
		Text("disable button", color.RGBA{230, 230, 230, 255}).
		ToggledHandler(func(args *widget.CheckBoxToggledEventArgs) {
			btn.SetDisabled(args.CheckBox.Checked)
		})

	cb2 := widget.NewCheckBox(20, 5, cb2Opts)
	container.AddComponent(cb2)

	lblOpts := &widget.LabelOptions{}
	lbl := widget.NewLabel(5, 40, "label", color.RGBA{230, 230, 230, 255}, lblOpts.Left())
	container.AddComponent(lbl)

	btn2 := widget.NewButton(5, 52, nil)
	container.AddComponent(btn2)

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
	g.mouse.Draw()

	g.gui.Update(g.mouse)
	g.gui.Draw(screen, g.mouse)
}
