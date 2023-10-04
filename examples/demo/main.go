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
	guiImage *ebiten.Image

	gui   *chopstiqs.Gui
	mouse *input.Mouse

	backgroundColor imgColor.RGBA
	bgColorToggled  bool

	screenWidth  int
	screenHeight int

	quitIsPressed bool
}

// New generates a new Game object.
func NewGame() *Game {
	g := &Game{
		gui:             &chopstiqs.Gui{},
		mouse:           input.NewMouse(),
		screenWidth:     200,
		screenHeight:    200,
		backgroundColor: imgColor.RGBA{9, 32, 42, 255},
	}

	ebiten.SetWindowSize(g.getWindowSize())
	ebiten.SetWindowTitle("chopstiqs demo")

	btnOpts := &widget.ButtonOptions{}
	btn := widget.NewButton(5, 20, btnOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
		g.toggleBgColor()
	}).Text(5, 20, "toggle background", color.RGBA{25, 25, 25, 255}))
	g.gui.AddComponent(btn)

	cbOpts := &widget.CheckBoxOptions{}
	cb := widget.NewCheckBox(5, 5, cbOpts.ToggledHandler(func(args *widget.CheckBoxToggledEventArgs) {
		btn.SetDisabled(args.CheckBox.Checked)
	}))
	g.gui.AddComponent(cb)

	lblOpts := &widget.LabelOptions{}
	lbl := widget.NewLabel(5, 40, "label", color.RGBA{230, 230, 230, 255}, lblOpts.Left())
	g.gui.AddComponent(lbl)

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

func (g *Game) toggleBgColor() {
	if !g.bgColorToggled {
		g.backgroundColor = imgColor.RGBA{32, 32, 32, 255}
	} else {
		g.backgroundColor = imgColor.RGBA{9, 32, 42, 255}
	}

	g.bgColorToggled = !g.bgColorToggled
}

// Draw draws the current game to the given screen.
func (g *Game) Draw(screen *ebiten.Image) {
	g.mouse.Draw()

	if g.guiImage == nil {
		g.guiImage = ebiten.NewImage(g.screenWidth, g.screenHeight)
	}

	screen.Fill(g.backgroundColor)

	g.gui.Update(g.mouse)
	g.gui.Draw(g.guiImage, g.mouse)

	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(g.guiImage, op)
}
