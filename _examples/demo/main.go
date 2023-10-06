package main

import (
	"errors"
	"image/color"
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

	container := widget.NewContainer(0, 0, 200, 200, color.RGBA{9, 32, 42, 255})

	btnLblOpts := &widget.LabelOptions{}
	btnOpts := &widget.ButtonOptions{}
	btn := widget.NewButton(btnOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
		if !g.bgColorToggled {
			container.SetBackgroundColor(color.RGBA{32, 32, 32, 255})
		} else {
			container.SetBackgroundColor(color.RGBA{9, 32, 42, 255})
		}

		g.bgColorToggled = !g.bgColorToggled
	}).Label("toggle background", btnLblOpts.Color(color.RGBA{25, 25, 25, 255})))
	container.AddComponent(5, 20, btn)

	btn2Opts := &widget.ButtonOptions{}
	btn2 := widget.NewButton(btn2Opts.Color(color.RGBA{100, 180, 90, 255}, color.RGBA{90, 160, 80, 255}, color.RGBA{120, 190, 100, 255}, color.RGBA{80, 100, 70, 255}))
	container.AddComponent(5, 52, btn2)

	cbOpts := &widget.CheckBoxOptions{}
	cb := widget.NewCheckBox(cbOpts.Color(color.RGBA{255, 100, 50, 255}))
	cb.Toggle()
	container.AddComponent(5, 5, cb)

	cb2LblOpts := &widget.LabelOptions{}
	cb2Opts := &widget.CheckBoxOptions{}
	cb2Opts = cb2Opts.
		Label("disable buttons", cb2LblOpts.Color(color.RGBA{230, 230, 230, 255})).
		ToggledHandler(func(args *widget.CheckBoxToggledEventArgs) {
			btn.SetDisabled(args.CheckBox.Checked)
			btn2.SetDisabled(args.CheckBox.Checked)
		})

	cb2 := widget.NewCheckBox(cb2Opts)
	container.AddComponent(20, 5, cb2)

	lblOpts := &widget.LabelOptions{}
	lbl := widget.NewLabel("label", lblOpts.Color(color.RGBA{120, 190, 100, 255}).Left())
	container.AddComponent(5, 40, lbl)

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
