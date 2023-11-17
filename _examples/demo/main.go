package main

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"image/color"
	"log"
	"strconv"

	"github.com/fglo/chopstiqs/pkg/component"
	"github.com/fglo/chopstiqs/pkg/debug"
	"github.com/fglo/chopstiqs/pkg/gui"
	"github.com/fglo/chopstiqs/pkg/to"
	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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
		screenWidth:     200,
		screenHeight:    200,
		backgroundColor: color.RGBA{32, 32, 32, 255},
	}

	ebiten.SetWindowSize(g.getWindowSize())
	ebiten.SetWindowTitle("chopstiqs demo")

	// component.SetDefaultPadding(2, 2, 2, 2)

	// rootContainer := component.NewContainer(&component.ContainerOptions{Width: to.Ptr(200), Height: to.Ptr(200)})
	rootContainer := component.NewContainer(&component.ContainerOptions{
		Padding: &component.Padding{Left: 5, Right: 5, Top: 5, Bottom: 5},
		Layout:  &component.VerticalListLayout{RowGap: 5}})
	// rootContainer := component.NewContainer(&component.ContainerOptions{
	// 	Padding: &component.Padding{Left: 5, Right: 5, Top: 5, Bottom: 5},
	// 	Layout:  &component.GridLayout{Columns: 2, ColumnGap: 5, Rows: 2, RowGap: 5}})

	lblTitle := component.NewLabel("chopstiqs demo", &component.LabelOptions{Color: color.RGBA{120, 190, 100, 255}, VerticalAlignment: component.AlignmentTop})

	lblInstructions := component.NewLabel("b - show borders\np - show padding\nq - quit", &component.LabelOptions{Color: color.RGBA{120, 120, 120, 255}, VerticalAlignment: component.AlignmentTop})

	cbOpts := &component.CheckBoxOptions{
		Drawer: component.DefaultCheckBoxDrawer{
			Color: color.RGBA{255, 100, 50, 255},
		},
	}
	cb := component.NewCheckBox(cbOpts)
	cb.Toggle()

	btn := component.NewButton(&component.ButtonOptions{
		Label: component.NewLabel("toggle background", &component.LabelOptions{Color: color.RGBA{25, 25, 25, 255}}),
	})
	btn.AddClickedHandler(func(args *component.ButtonClickedEventArgs) { g.toggleBackground() })

	btn2 := component.NewButton(&component.ButtonOptions{
		Drawer: component.DefaultButtonDrawer{
			Color:         color.RGBA{100, 180, 90, 255},
			ColorPressed:  color.RGBA{90, 160, 80, 255},
			ColorHovered:  color.RGBA{120, 190, 100, 255},
			ColorDisabled: color.RGBA{80, 100, 70, 255},
		},
	})

	sliderLabel := component.NewLabel("4", &component.LabelOptions{
		Color: color.RGBA{230, 230, 230, 255},
		Padding: &component.Padding{
			Top: 4,
		},
	})

	slider := component.NewSlider(&component.SliderOptions{
		Min:          to.Ptr(0.),
		Max:          to.Ptr(10.),
		Step:         to.Ptr(1.),
		DefaultValue: to.Ptr(4.),
		Width:        to.Ptr(100),
		Height:       to.Ptr(15),
	})

	slider.AddSlidedHandler(func(args *component.SliderSlidedEventArgs) {
		sliderLabel.SetText(strconv.Itoa(int(args.Value)))
	})

	sliderContainer := component.NewContainer(&component.ContainerOptions{Layout: &component.HorizontalListLayout{ColumnGap: 5}})
	sliderContainer.AddComponent(slider)
	sliderContainer.AddComponent(sliderLabel)

	sliderLabel2 := component.NewLabel("0.5", &component.LabelOptions{
		Color: color.RGBA{230, 230, 230, 255},
		Padding: &component.Padding{
			Top: 4,
		},
	})

	slider2 := component.NewSlider(&component.SliderOptions{
		Min:          to.Ptr(0.),
		Max:          to.Ptr(1.),
		Step:         to.Ptr(.05),
		DefaultValue: to.Ptr(.5),
		Width:        to.Ptr(100),
		Height:       to.Ptr(15),
	})

	slider2.AddSlidedHandler(func(args *component.SliderSlidedEventArgs) {
		sliderLabel2.SetText(fmt.Sprintf("%.2f", args.Value))
	})

	sliderContainer2 := component.NewContainer(&component.ContainerOptions{Layout: &component.HorizontalListLayout{ColumnGap: 5}})
	sliderContainer2.AddComponent(slider2)
	sliderContainer2.AddComponent(sliderLabel2)

	cb2Opts := &component.CheckBoxOptions{
		Label: component.NewLabel("disable buttons", &component.LabelOptions{Color: color.RGBA{230, 230, 230, 255}}),
	}
	cb2 := component.NewCheckBox(cb2Opts)
	cb2.AddToggledHandler(func(args *component.CheckBoxToggledEventArgs) {
		btn.SetDisabled(args.CheckBox.Checked)
		btn2.SetDisabled(args.CheckBox.Checked)
		sliderContainer.SetDisabled(args.CheckBox.Checked)
		sliderContainer2.SetDisabled(args.CheckBox.Checked)
	})

	checkBoxContainer := component.NewContainer(&component.ContainerOptions{Layout: &component.HorizontalListLayout{ColumnGap: 5}})
	checkBoxContainer.AddComponent(cb)
	checkBoxContainer.AddComponent(cb2)

	img, _, _ := ebitenutil.NewImageFromReader(bytes.NewReader(chopstiqsLogo))
	sprite := component.NewSprite(img, nil)

	rootContainer.AddComponent(sprite)
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
	rootContainer.AddComponent(sliderContainer)
	rootContainer.AddComponent(sliderContainer2)

	gui.SetRootContainer(rootContainer)

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
	gui.Update()

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
	gui.Draw(screen)
}
