package component

import (
	"testing"

	"github.com/fglo/chopstiqs/pkg/event"
	ebiten "github.com/hajimehoshi/ebiten/v2"
)

func leftMouseButtonClick(t *testing.T, c *component) {
	t.Helper()

	leftMouseButtonPress(t, c)
	leftMouseButtonRelease(t, c)
}

func leftMouseButtonPress(t *testing.T, c *component) {
	t.Helper()

	c.MouseButtonPressedEvent.Fire(&ComponentMouseButtonPressedEventArgs{
		Component: c,
		Button:    ebiten.MouseButtonLeft,
	})

	event.HandleFired()
}

func leftMouseButtonRelease(t *testing.T, c *component) {
	t.Helper()

	c.MouseButtonReleasedEvent.Fire(&ComponentMouseButtonReleasedEventArgs{
		Component: c,
		Button:    ebiten.MouseButtonLeft,
		Inside:    true,
	})

	event.HandleFired()
}
