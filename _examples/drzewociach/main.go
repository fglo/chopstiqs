package main

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"image/color"
	"log"

	"github.com/fglo/chopstiqs"
	"github.com/fglo/chopstiqs/component"
	"github.com/fglo/chopstiqs/debug"
	"github.com/fglo/chopstiqs/option"
	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

//go:embed img/Tree1_idle_ns.png
var tree1Idle []byte

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
	gui *chopstiqs.GUI

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
	gui := chopstiqs.NewGUI(&chopstiqs.GUIOptions{
		HorizontalAlignment: option.AlignmentLeft,
		VerticalAlignment:   option.AlignmentTop,
	})

	g := &Game{
		gui:             gui,
		screenWidth:     385,
		screenHeight:    455,
		backgroundColor: color.RGBA{32, 32, 32, 255},
	}

	ebiten.SetWindowSize(g.getWindowSize())
	ebiten.SetWindowTitle("chopstiqs demo")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	component.SetDefaultPadding(1, 1, 1, 1)

	rootContainer := gui.NewContainer(&component.ContainerOptions{
		Padding: &component.Padding{Left: 5, Right: 5, Top: 5, Bottom: 5},
		Layout:  &component.HorizontalListLayout{ColumnGap: 5}})

	gui.SetRootContainer(rootContainer)

	drzewoContainer := gui.NewContainer(&component.ContainerOptions{
		Padding:             &component.Padding{Left: 5, Right: 5, Top: 5, Bottom: 5},
		HorizontalAlignment: option.StretchedHorizontally,
		Layout:              &component.VerticalListLayout{RowGap: 5}})

	ciached := 0

	label := gui.NewLabel("DRZEWO: 0", &component.LabelOptions{
		Color:               color.RGBA{120, 190, 100, 255},
		VerticalAlignment:   option.AlignmentCenteredVertically,
		HorizontalAlignment: option.AlignmentCenteredHorizontally,
	})

	drzewoContainer.AddComponent(label)

	img, _, _ := ebitenutil.NewImageFromReader(bytes.NewReader(tree1Idle))
	sprite := gui.NewSprite(img, nil)

	drzewoContainer.AddComponent(sprite)

	ciachButton := gui.NewButton(&component.ButtonOptions{
		Label: gui.NewLabel("CIACH", &component.LabelOptions{Color: color.RGBA{50, 50, 50, 255}}),
	})
	ciachButton.AddClickedHandler(func(args *component.ButtonClickedEventArgs) {
		ciached += 1
		label.SetText(fmt.Sprintf("DRZEWO: %d", ciached))
	})

	drzewoContainer.AddComponent(ciachButton)

	upgradesContainer := gui.NewContainer(&component.ContainerOptions{
		Padding:             &component.Padding{Left: 5, Right: 5, Top: 5, Bottom: 5},
		HorizontalAlignment: option.StretchedHorizontally,
		VerticalAlignment:   option.AlignmentCenteredVertically,
		Layout:              &component.VerticalListLayout{RowGap: 5}})

	labelUpgrades := gui.NewLabel("ULEPSZENIA", &component.LabelOptions{
		Color:               color.RGBA{120, 190, 100, 255},
		VerticalAlignment:   option.AlignmentTop,
		HorizontalAlignment: option.AlignmentCenteredHorizontally,
	})

	upgradesContainer.AddComponent(labelUpgrades)

	upgradeSpeedButton := gui.NewButton(&component.ButtonOptions{
		Label:               gui.NewLabel("SZYBKOŚĆ", &component.LabelOptions{Color: color.RGBA{50, 50, 50, 255}}),
		HorizontalAlignment: option.AlignmentRight,
	})

	upgradesContainer.AddComponent(upgradeSpeedButton)

	upgradeOutputButton := gui.NewButton(&component.ButtonOptions{
		Label:               gui.NewLabel("JAKOŚĆ", &component.LabelOptions{Color: color.RGBA{50, 50, 50, 255}}),
		HorizontalAlignment: option.StretchedHorizontally,
	})

	upgradesContainer.AddComponent(upgradeOutputButton)

	rootContainer.AddComponent(drzewoContainer)
	rootContainer.AddComponent(upgradesContainer)

	return g
}

func (g *Game) getWindowSize() (int, int) {
	return g.screenWidth, g.screenHeight
}

// Layout implements ebiten.Game's Layout.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	w, h := ebiten.WindowSize()
	if w == 0 || h == 0 {
		return int(float64(g.screenWidth) / 2), int(float64(g.screenHeight) / 2)
	}

	return int(float64(w) / 2), int(float64(h) / 2) // TODO: proper scaling
}

// Update updates the current game state.
func (g *Game) Update() error {
	g.gui.Update()

	g.checkShowBordersButton()
	g.checkShowPaddingButton()

	return g.checkQuitButton()
}

func (g *Game) textInputHasFocus() bool {
	_, ok := g.gui.FocusedComponent().(*component.TextInput)
	return ok
}

func (g *Game) checkQuitButton() error {
	if g.textInputHasFocus() {
		return nil
	}
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
	if g.textInputHasFocus() {
		return
	}
	if !g.showBordersIsPressed && inpututil.IsKeyJustPressed(ebiten.KeyB) {
		g.showBordersIsPressed = true
	}
	if g.showBordersIsPressed && inpututil.IsKeyJustReleased(ebiten.KeyB) {
		g.showBordersIsPressed = false
		debug.ShowComponentBorders = !debug.ShowComponentBorders
	}
}

func (g *Game) checkShowPaddingButton() {
	if g.textInputHasFocus() {
		return
	}
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
