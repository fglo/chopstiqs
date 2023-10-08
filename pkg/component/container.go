package component

import (
	imgColor "image/color"
)

type Container interface {
	Component
	AddComponent(Component)
	SetBackgroundColor(imgColor.RGBA)
	GetBackgroundColor() imgColor.RGBA
}
