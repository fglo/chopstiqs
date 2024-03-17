package component

import (
	"encoding/xml"
	"fmt"
)

type Padding struct {
	Top    int
	Right  int
	Bottom int
	Left   int
}

var DefaultPadding Padding = Padding{
	Top:    0,
	Bottom: 0,
	Left:   0,
	Right:  0,
}

func SetDefaultPadding(left, right, top, bottom int) {
	DefaultPadding = Padding{
		Top:    top,
		Bottom: bottom,
		Left:   left,
		Right:  right,
	}
}

func SetDefaultHorizontalPadding(padding int) {
	DefaultPadding.Left = padding
	DefaultPadding.Right = padding
}

func SetDefaultVerticalPadding(padding int) {
	DefaultPadding.Top = padding
	DefaultPadding.Bottom = padding
}

func (p *Padding) Validate() {
	if p.Top < 0 {
		p.Top = DefaultPadding.Top
	}

	if p.Bottom < 0 {
		p.Bottom = DefaultPadding.Bottom
	}

	if p.Left < 0 {
		p.Left = DefaultPadding.Left
	}

	if p.Right < 0 {
		p.Right = DefaultPadding.Right
	}
}

func NewPadding(top, right, bottom, left int) *Padding {
	p := &Padding{
		Top:    top,
		Right:  right,
		Bottom: bottom,
		Left:   left,
	}

	p.Validate()

	return p
}

func (p *Padding) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	if p.Top == 0 && p.Right == 0 && p.Bottom == 0 && p.Left == 0 {
		return xml.Attr{}, nil
	}

	return xml.Attr{
		Name:  name,
		Value: fmt.Sprintf("%d,%d,%d,%d", p.Top, p.Right, p.Bottom, p.Left),
	}, nil
}

func (p *Padding) UnmarshalXMLAttr(attr xml.Attr) error {
	_, err := fmt.Sscanf(attr.Value, "%d,%d,%d,%d", &p.Top, &p.Right, &p.Bottom, &p.Left)
	return err
}
