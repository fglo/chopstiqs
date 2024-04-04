package component

import (
	"testing"

	"github.com/fglo/chopstiqs/input"
	ebiten "github.com/hajimehoshi/ebiten/v2"
)

func leftMouseButtonClick(t *testing.T, c *component) {
	t.Helper()

	leftMouseButtonPress(t, c)
	leftMouseButtonRelease(t, c)
}

func leftMouseButtonPress(t *testing.T, c *component) {
	t.Helper()

	c.eventManager.Fire(c.MouseButtonPressedEvent, &ComponentMouseButtonPressedEventArgs{
		Component: c,
		Button:    ebiten.MouseButtonLeft,
	})

	c.eventManager.HandleFired()
}

func leftMouseButtonRelease(t *testing.T, c *component) {
	t.Helper()

	c.eventManager.Fire(c.MouseButtonReleasedEvent, &ComponentMouseButtonReleasedEventArgs{
		Component: c,
		Button:    ebiten.MouseButtonLeft,
		Inside:    true,
	})

	c.eventManager.HandleFired()
}

func keyPress(t *testing.T, key ebiten.Key) {
	t.Helper()

	input.KeyPressed[key] = true
}

func keyRelease(t *testing.T, key ebiten.Key) {
	t.Helper()

	input.KeyPressed[key] = false
}
