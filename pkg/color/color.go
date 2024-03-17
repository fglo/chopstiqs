package color

import (
	"encoding/xml"
	"fmt"
	color "image/color"
)

type Color color.Color
type RGBASerializable color.RGBA

func ColorToRGBA(c color.Color) color.RGBA {
	r, g, b, a := c.RGBA()
	return color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
}

func (rgba RGBASerializable) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	return xml.Attr{
		Name:  name,
		Value: ToHex(color.RGBA(rgba)),
	}, nil
}

func (rgba *RGBASerializable) UnmarshalXMLAttr(attr xml.Attr) error {
	c, err := ParseHex(attr.Value)
	if err != nil {
		return err // nolint
	}

	*rgba = RGBASerializable(c)

	return nil
}

func ToHex(color color.RGBA) string {
	return fmt.Sprintf("#%02x%02x%02x%02x", color.R, color.G, color.B, color.A)
}

func ParseHex(s string) (c color.RGBA, err error) {
	if s[0] == '#' {
		s = s[1:]
	}

	switch len(s) {
	case 8:
		_, err = fmt.Sscanf(s, "%02x%02x%02x%02x", &c.R, &c.G, &c.B, &c.A)
	case 6:
		_, err = fmt.Sscanf(s, "%02x%02x%02x", &c.R, &c.G, &c.B)
	case 4:
		_, err = fmt.Sscanf(s, "%1x%1x%1x%1x", &c.R, &c.G, &c.B, &c.A)
		// Double the hex digits:
		c.R *= 17
		c.G *= 17
		c.B *= 17
		c.A *= 17
	case 3:
		_, err = fmt.Sscanf(s, "%1x%1x%1x", &c.R, &c.G, &c.B)
		// Double the hex digits:
		c.R *= 17
		c.G *= 17
		c.B *= 17
	default:
		err = fmt.Errorf("invalid format")
	}

	return
}
