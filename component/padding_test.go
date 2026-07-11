package component

import (
	"testing"

	"github.com/matryer/is"
)

func TestNewPadding(t *testing.T) {
	is := is.New(t)

	p := NewPadding(1, 2, 3, 4)

	is.Equal(p.Top, 1)
	is.Equal(p.Right, 2)
	is.Equal(p.Bottom, 3)
	is.Equal(p.Left, 4)
}

func TestNewPadding_NegativeValues_ReplacedWithDefaults(t *testing.T) {
	is := is.New(t)

	origDefault := DefaultPadding
	defer func() { DefaultPadding = origDefault }()

	SetDefaultPadding(10, 20, 30, 40)

	p := NewPadding(-1, -2, -3, -4)

	is.Equal(p.Top, 30)
	is.Equal(p.Right, 20)
	is.Equal(p.Bottom, 40)
	is.Equal(p.Left, 10)
}

func TestPadding_Validate(t *testing.T) {
	is := is.New(t)

	origDefault := DefaultPadding
	defer func() { DefaultPadding = origDefault }()

	SetDefaultPadding(5, 6, 7, 8)

	p := &Padding{Top: -1, Right: -2, Bottom: -3, Left: -4}
	p.Validate()

	is.Equal(p.Top, 7)
	is.Equal(p.Right, 6)
	is.Equal(p.Bottom, 8)
	is.Equal(p.Left, 5)
}

func TestPadding_Validate_PositiveValuesUnchanged(t *testing.T) {
	is := is.New(t)

	p := &Padding{Top: 10, Right: 20, Bottom: 30, Left: 40}
	p.Validate()

	is.Equal(p.Top, 10)
	is.Equal(p.Right, 20)
	is.Equal(p.Bottom, 30)
	is.Equal(p.Left, 40)
}

func TestSetDefaultPadding(t *testing.T) {
	is := is.New(t)

	origDefault := DefaultPadding
	defer func() { DefaultPadding = origDefault }()

	SetDefaultPadding(1, 2, 3, 4)

	is.Equal(DefaultPadding.Left, 1)
	is.Equal(DefaultPadding.Right, 2)
	is.Equal(DefaultPadding.Top, 3)
	is.Equal(DefaultPadding.Bottom, 4)
}

func TestSetDefaultHorizontalPadding(t *testing.T) {
	is := is.New(t)

	origDefault := DefaultPadding
	defer func() { DefaultPadding = origDefault }()

	SetDefaultPadding(0, 0, 0, 0)
	SetDefaultHorizontalPadding(15)

	is.Equal(DefaultPadding.Left, 15)
	is.Equal(DefaultPadding.Right, 15)
}

func TestSetDefaultVerticalPadding(t *testing.T) {
	is := is.New(t)

	origDefault := DefaultPadding
	defer func() { DefaultPadding = origDefault }()

	SetDefaultPadding(0, 0, 0, 0)
	SetDefaultVerticalPadding(25)

	is.Equal(DefaultPadding.Top, 25)
	is.Equal(DefaultPadding.Bottom, 25)
}
