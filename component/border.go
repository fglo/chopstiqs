package component

import imgColor "image/color"

type Border struct {
	Width int
	Color imgColor.RGBA
}

var DefaultBorder Border = Border{
	Width: 0,
}
