package main

import (
	_ "embed"
	"encoding/xml"
	"errors"
	"image/color"
	"log"
	"os"

	"github.com/fglo/chopstiqs/pkg/debug"
	"github.com/fglo/chopstiqs/pkg/gui"
	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

//go:embed img/chopstiqs-logo-1-1x1.png
var chopstiqsLogo []byte

var Terminated = errors.New("terminated")

func main() {
	debug.Debug = true

	if err := ebiten.RunGame(NewGame()); err != nil {
		if err == Terminated {
			return
		}

		log.Fatal(err)
	}
}

// Game encapsulates game logic
type Game struct {
	gui *gui.GUI

	bgColorToggled bool

	screenWidth  int
	screenHeight int

	quitIsPressed        bool
	showBordersIsPressed bool
	showPaddingIsPressed bool

	backgroundColor color.RGBA
}

// New generates a new Game object.
func NewGame() *Game {
	g := &Game{
		gui:             gui.New(),
		screenWidth:     200,
		screenHeight:    200,
		backgroundColor: color.RGBA{32, 32, 32, 255},
	}

	ebiten.SetWindowSize(g.getWindowSize())
	ebiten.SetWindowTitle("chopstiqs demo")

	file, _ := os.ReadFile("unmarshal_test.xml")

	_ = xml.Unmarshal(file, g.gui)

	return g
}

func (g *Game) getWindowSize() (int, int) {
	var factor float64 = 2
	return int(float64(g.screenWidth) * factor), int(float64(g.screenHeight) * factor)
}

// Layout implements ebiten.Game's Layout.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.screenWidth, g.screenHeight
}

// Update updates the current game state.
func (g *Game) Update() error {
	g.gui.Update()

	g.checkShowBordersButton()
	g.checkShowPaddingButton()

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

func (g *Game) checkShowBordersButton() {
	if !g.showBordersIsPressed && inpututil.IsKeyJustPressed(ebiten.KeyB) {
		g.showBordersIsPressed = true
	}
	if g.showBordersIsPressed && inpututil.IsKeyJustReleased(ebiten.KeyB) {
		g.showBordersIsPressed = false
		debug.ShowComponentBorders = !debug.ShowComponentBorders
	}
}

func (g *Game) checkShowPaddingButton() {
	if !g.showPaddingIsPressed && inpututil.IsKeyJustPressed(ebiten.KeyP) {
		g.showPaddingIsPressed = true
	}
	if g.showPaddingIsPressed && inpututil.IsKeyJustReleased(ebiten.KeyP) {
		g.showPaddingIsPressed = false
		debug.ShowComponentPadding = !debug.ShowComponentPadding
	}
}

func (g *Game) toggleBackground() {
	if !g.bgColorToggled {
		g.backgroundColor = color.RGBA{9, 32, 42, 255}
	} else {
		g.backgroundColor = color.RGBA{32, 32, 32, 255}
	}

	g.bgColorToggled = !g.bgColorToggled
}

// Draw draws the current game to the given screen.
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(g.backgroundColor)
	g.gui.Draw(screen)
}
