package color

import imgColor "image/color"

func ToRGBA(color imgColor.Color) imgColor.RGBA {
	r, g, b, a := color.RGBA()
	return imgColor.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
}

func Invert(color imgColor.Color) imgColor.RGBA {
	r, g, b, a := color.RGBA()
	return imgColor.RGBA{uint8(255 - r), uint8(255 - g), uint8(255 - b), uint8(a)}
}
