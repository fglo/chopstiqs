package colorutils

import imgColor "image/color"

func ColorToRGBA(color imgColor.Color) imgColor.RGBA {
	r, g, b, a := color.RGBA()
	return imgColor.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
}
