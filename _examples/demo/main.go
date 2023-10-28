package main

import (
	"errors"
	"image/color"
	"log"

	"github.com/fglo/chopstiqs/pkg/component"
	"github.com/fglo/chopstiqs/pkg/debug"
	"github.com/fglo/chopstiqs/pkg/gui"
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
	bgColorToggled bool

	screenWidth  int
	screenHeight int

	quitIsPressed  bool
	debugIsPressed bool

	backgroundColor color.RGBA
}

// New generates a new Game object.
func NewGame() *Game {
	g := &Game{
		screenWidth:     200,
		screenHeight:    200,
		backgroundColor: color.RGBA{32, 32, 32, 255},
	}

	ebiten.SetWindowSize(g.getWindowSize())
	ebiten.SetWindowTitle("chopstiqs demo")

	// component.SetDefaultPadding(3, 3, 3, 3)

	// rootContainer := component.NewContainer(&component.ContainerOptions{Width: to.Ptr(200), Height: to.Ptr(200)})
	rootContainer := component.NewContainer(&component.ContainerOptions{
		Padding: &component.Padding{Left: 5, Right: 5, Top: 5, Bottom: 5},
		Layout:  &component.VerticalListLayout{RowGap: 5}})
	// rootContainer := component.NewContainer(&component.ContainerOptions{
	// 	Padding: &component.Padding{Left: 5, Right: 5, Top: 5, Bottom: 5},
	// 	Layout:  &component.GridLayout{Columns: 2, ColumnGap: 5, Rows: 2, RowGap: 5}})

	lblTitle := component.NewLabel("chopstiqs demo", &component.LabelOptions{Color: color.RGBA{120, 190, 100, 255}, VerticalAlignment: component.AlignmentTop})

	lblInstructions := component.NewLabel("d - debug\nq - quit", &component.LabelOptions{Color: color.RGBA{120, 120, 120, 255}, VerticalAlignment: component.AlignmentTop})

	cbOpts := &component.CheckBoxOptions{
		Color: color.RGBA{255, 100, 50, 255},
	}
	cb := component.NewCheckBox(cbOpts)
	cb.Toggle()

	cb2Opts := &component.CheckBoxOptions{
		Color: color.RGBA{230, 230, 230, 255},
		Label: component.NewLabel("disable buttons", &component.LabelOptions{Color: color.RGBA{230, 230, 230, 255}}),
	}

	btn := component.NewButton(&component.ButtonOptions{
		Label: component.NewLabel("toggle background", &component.LabelOptions{Color: color.RGBA{25, 25, 25, 255}}),
	})
	btn.AddClickedHandler(func(args *component.ButtonClickedEventArgs) { g.toggleBackground() })

	btn2 := component.NewButton(&component.ButtonOptions{
		Color:         color.RGBA{100, 180, 90, 255},
		ColorPressed:  color.RGBA{90, 160, 80, 255},
		ColorHovered:  color.RGBA{120, 190, 100, 255},
		ColorDisabled: color.RGBA{80, 100, 70, 255},
	})

	cb2 := component.NewCheckBox(cb2Opts)
	cb2.AddToggledHandler(func(args *component.CheckBoxToggledEventArgs) {
		btn.SetDisabled(args.CheckBox.Checked)
		btn2.SetDisabled(args.CheckBox.Checked)
	})

	checkBoxContainer := component.NewContainer(&component.ContainerOptions{Layout: &component.HorizontalListLayout{ColumnGap: 5}})
	checkBoxContainer.AddComponent(cb)
	checkBoxContainer.AddComponent(cb2)

	lblTitle.SetPosision(5, 5)
	rootContainer.AddComponent(lblTitle)
	lblInstructions.SetPosision(5, 15)
	rootContainer.AddComponent(lblInstructions)
	checkBoxContainer.SetPosision(5, 45)
	rootContainer.AddComponent(checkBoxContainer)
	btn.SetPosision(5, 60)
	rootContainer.AddComponent(btn)
	btn2.SetPosision(5, 75)
	rootContainer.AddComponent(btn2)

	gui.SetRootContainer(rootContainer)

	return g
}

func (g *Game) getWindowSize() (int, int) {
	var factor float32 = 2
	return int(float32(g.screenWidth) * factor), int(float32(g.screenHeight) * factor)
}

// Layout implements ebiten.Game's Layout.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.screenWidth, g.screenHeight
}

// Update updates the current game state.
func (g *Game) Update() error {
	gui.Update()

	g.checkDebugButton()

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

func (g *Game) checkDebugButton() {
	if !g.debugIsPressed && inpututil.IsKeyJustPressed(ebiten.KeyD) {
		g.debugIsPressed = true
	}
	if g.debugIsPressed && inpututil.IsKeyJustReleased(ebiten.KeyD) {
		g.debugIsPressed = false
		debug.Debug = !debug.Debug
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
	gui.Draw(screen)
}
