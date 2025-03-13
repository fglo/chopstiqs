package main

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"image/color"
	"log"
	"sort"
	"strconv"

	"github.com/fglo/chopstiqs"
	"github.com/fglo/chopstiqs/component"
	"github.com/fglo/chopstiqs/debug"
	"github.com/fglo/chopstiqs/input"
	"github.com/fglo/chopstiqs/option"
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
	gui *chopstiqs.GUI

	bgColorToggled bool

	screenWidth  int
	screenHeight int

	quitIsPressed        bool
	showBordersIsPressed bool
	showPaddingIsPressed bool

	backgroundColor color.RGBA
}

var lblPressedKeys *component.Label

// New generates a new Game object.
func NewGame() *Game {
	gui := chopstiqs.NewGUI().WithOptions(chopstiqs.GUIOptions{
		HorizontalAlignment: option.AlignmentCenteredHorizontally,
		VerticalAlignment:   option.AlignmentTop,
	})

	g := &Game{
		gui:             gui,
		screenWidth:     220,
		screenHeight:    260,
		backgroundColor: color.RGBA{32, 32, 32, 255},
	}

	ebiten.SetWindowSize(g.getWindowSize())
	ebiten.SetWindowTitle("chopstiqs demo")

	// component.SetDefaultPadding(2, 2, 2, 2)

	rootContainer := gui.NewContainer(&component.ContainerOptions{
		Padding: &component.Padding{Left: 5, Right: 5, Top: 5, Bottom: 5},
		Layout:  &component.VerticalListLayout{RowGap: 5}})

	gui.SetRootContainer(rootContainer)

	lblTitle := gui.NewLabel("chopstiqs demo", &component.LabelOptions{Color: color.RGBA{120, 190, 100, 255}, VerticalAlignment: option.AlignmentTop})

	lblInstructions := gui.NewLabel("b - show borders\np - show padding\nq - quit", &component.LabelOptions{Color: color.RGBA{120, 120, 120, 255}, VerticalAlignment: option.AlignmentTop})
	lblPressedKeys = gui.NewLabel("[]", &component.LabelOptions{Color: color.RGBA{120, 120, 120, 255}, VerticalAlignment: option.AlignmentTop})

	cbOpts := &component.CheckBoxOptions{
		Drawer: component.DefaultCheckBoxDrawer{
			Color: color.RGBA{255, 100, 50, 255},
		},
		Width: option.Int(15),
	}
	cb := gui.NewCheckBox(cbOpts)
	cb.Toggle()

	btn := gui.NewButton(&component.ButtonOptions{
		Label: gui.NewLabel("toggle background", &component.LabelOptions{Color: color.RGBA{50, 50, 50, 255}}),
	})
	btn.AddClickedHandler(func(args *component.ButtonClickedEventArgs) { g.toggleBackground() })

	btn2 := gui.NewButton(&component.ButtonOptions{
		Drawer: &component.DefaultButtonDrawer{
			Color:         color.RGBA{100, 180, 90, 255},
			ColorPressed:  color.RGBA{90, 160, 80, 255},
			ColorHovered:  color.RGBA{120, 190, 100, 255},
			ColorDisabled: color.RGBA{80, 100, 70, 255},
		},
	})

	sliderLabel := gui.NewLabel("4", &component.LabelOptions{
		Color: color.RGBA{230, 230, 230, 255},
		Padding: &component.Padding{
			Top: 4,
		},
	})

	slider := gui.NewSlider(&component.SliderOptions{
		Min:          option.Float(0.),
		Max:          option.Float(10.),
		Step:         option.Float(1.),
		DefaultValue: option.Float(4.),
		Width:        option.Int(100),
		Height:       option.Int(15),
	})

	slider.AddSlidedHandler(func(args *component.SliderSlidedEventArgs) {
		sliderLabel.SetText(strconv.Itoa(int(args.Value)))
		newBackgroundG := int8(g.backgroundColor.G) + int8(args.Change)*5
		if newBackgroundG > 0 {
			g.backgroundColor.G = uint8(newBackgroundG)
		}
	})

	sliderContainer := gui.NewContainer(&component.ContainerOptions{Layout: &component.HorizontalListLayout{ColumnGap: 5}})
	sliderContainer.AddComponent(slider)
	sliderContainer.AddComponent(sliderLabel)

	slider2TextInput := component.NewTextInput(&component.TextInputOptions{
		Width: option.Int(25),
		InputValidationFunc: func(s string) (bool, string) {
			val, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return false, ""
			}

			return true, fmt.Sprintf("%.2f", val)
		},
		OnSubmitFunc: func(s string) string {
			val, err := strconv.ParseFloat(s, 64)
			switch {
			case err != nil:
				return "0.50"
			case val < 0:
				return "0.00"
			case val > 1:
				return "1.00"
			default:
				return fmt.Sprintf("%.2f", val)
			}
		},
		SubmitOnUnfocus: true,
	})
	slider2TextInput.SetValue("0.50")

	slider2 := gui.NewSlider(&component.SliderOptions{
		Min:          option.Float(0.),
		Max:          option.Float(1.),
		Step:         option.Float(.05),
		DefaultValue: option.Float(.5),
		Width:        option.Int(100),
		Height:       option.Int(15),
	})

	slider2.AddSlidedHandler(func(args *component.SliderSlidedEventArgs) {
		slider2TextInput.SetValue(fmt.Sprintf("%.2f", args.Value))
	})

	slider2TextInput.AddSubmittedHandler(func(args *component.TextInputSubmittedEventArgs) {
		val, err := strconv.ParseFloat(slider2TextInput.Value(), 64)
		if err != nil {
			val = 0
		}
		slider2.Set(val)
	})

	sliderContainer2 := gui.NewContainer(&component.ContainerOptions{Layout: &component.HorizontalListLayout{ColumnGap: 5}})
	sliderContainer2.AddComponent(slider2)
	sliderContainer2.AddComponent(slider2TextInput)

	textInput := component.NewTextInput(&component.TextInputOptions{Width: option.Int(100)})
	textInput.SetValue("Lorem Ipsum dolor sit amet")

	cb2Opts := &component.CheckBoxOptions{
		Label: gui.NewLabel("disable components", &component.LabelOptions{Color: color.RGBA{230, 230, 230, 255}}),
	}
	cb2 := gui.NewCheckBox(cb2Opts)
	cb2.AddToggledHandler(func(args *component.CheckBoxToggledEventArgs) {
		btn.SetDisabled(args.CheckBox.Checked())
		btn2.SetDisabled(args.CheckBox.Checked())
		sliderContainer.SetDisabled(args.CheckBox.Checked())
		sliderContainer2.SetDisabled(args.CheckBox.Checked())
		textInput.SetDisabled(args.CheckBox.Checked())
	})

	checkBoxContainer := gui.NewContainer(&component.ContainerOptions{Layout: &component.HorizontalListLayout{ColumnGap: 5}})
	checkBoxContainer.AddComponent(cb)
	checkBoxContainer.AddComponent(cb2)

	abcContainer := gui.NewContainer(&component.ContainerOptions{
		Padding: &component.Padding{Left: 5, Right: 5, Top: 5, Bottom: 5},
		Layout:  &component.GridLayout{Columns: 4, ColumnGap: 5, Rows: 3, RowGap: 5}})

	abcContainer.AddComponent(gui.NewLabel("a", nil))
	abcContainer.AddComponent(gui.NewLabel("b", nil))
	abcContainer.AddComponent(gui.NewLabel("c", nil))
	abcContainer.AddComponent(gui.NewLabel("d", nil))
	abcContainer.AddComponent(gui.NewLabel("e", nil))
	abcContainer.AddComponent(gui.NewLabel("f", nil))
	abcContainer.AddComponent(gui.NewLabel("g", nil))
	abcContainer.AddComponent(gui.NewLabel("h", nil))
	abcContainer.AddComponent(gui.NewLabel("i", nil))
	abcContainer.AddComponent(gui.NewLabel("j", nil))
	abcContainer.AddComponent(gui.NewLabel("k", nil))
	abcContainer.AddComponent(gui.NewLabel("l", nil))

	img, _, _ := ebitenutil.NewImageFromReader(bytes.NewReader(chopstiqsLogo))
	sprite := component.NewSprite(img, nil)

	rootContainer.AddComponent(sprite)
	lblTitle.SetPosition(5, 5)
	rootContainer.AddComponent(lblTitle)
	lblInstructions.SetPosition(5, 15)
	rootContainer.AddComponent(lblInstructions)
	rootContainer.AddComponent(lblPressedKeys)
	checkBoxContainer.SetPosition(5, 45)
	rootContainer.AddComponent(checkBoxContainer)
	btn.SetPosition(5, 60)
	rootContainer.AddComponent(btn)
	btn2.SetPosition(5, 75)
	rootContainer.AddComponent(btn2)
	rootContainer.AddComponent(sliderContainer)
	rootContainer.AddComponent(sliderContainer2)
	rootContainer.AddComponent(textInput)
	rootContainer.AddComponent(abcContainer)

	return g
}

func (g *Game) getWindowSize() (int, int) {
	var factor float64 = 1.75
	return int(float64(g.screenWidth) * factor), int(float64(g.screenHeight) * factor)
}

// Layout implements ebiten.Game's Layout.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.screenWidth, g.screenHeight
}

var pressedKeysStr string
var justPressedKeysStr string

// Update updates the current game state.
func (g *Game) Update() error {
	g.gui.Update()

	g.checkShowBordersButton()
	g.checkShowPaddingButton()

	_pressedKeys := make([]string, 0)
	_justPressedKeys := make([]string, 0)

	for k := ebiten.Key(0); k <= ebiten.KeyMax; k++ {
		if inpututil.IsKeyJustPressed(k) {
			_justPressedKeys = append(_justPressedKeys, k.String())
		}
		// if ebiten.IsKeyPressed(k) {
		// 	_pressedKeys = append(_pressedKeys, k.String())
		// }
	}

	for key, pressed := range input.KeyPressed {
		if pressed {
			_pressedKeys = append(_pressedKeys, ebiten.Key(key).String())
		}
	}

	sort.Strings(_pressedKeys)
	sort.Strings(_justPressedKeys)

	_pressedKeysStr := fmt.Sprintf("%v", _pressedKeys)
	_justPressedKeysStr := fmt.Sprintf("%v", _justPressedKeys)

	if _pressedKeysStr != pressedKeysStr {
		pressedKeysStr = _pressedKeysStr
	}

	if _justPressedKeysStr != justPressedKeysStr {
		justPressedKeysStr = _justPressedKeysStr
	}

	lblPressedKeys.SetText(_pressedKeysStr)

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
