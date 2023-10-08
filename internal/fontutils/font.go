package fontutils

import (
	_ "embed"
	"os"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

//go:embed fonts/Minecraftia-Regular.ttf
var tffFile []byte

var DefaultFontFace, _ = LoadFont(tffFile, 8)

func LoadFont(fontData []byte, size float64) (font.Face, error) {
	ttfFont, err := truetype.Parse(fontData)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(ttfFont, &truetype.Options{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	}), nil
}

func LoadFontFromFile(path string, size float64) (font.Face, error) {
	fontData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	ttfFont, err := truetype.Parse(fontData)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(ttfFont, &truetype.Options{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	}), nil
}
