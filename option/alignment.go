package option

type HorizontalAlignment int

const (
	AlignmentLeft HorizontalAlignment = iota
	AlignmentCenteredHorizontally
	AlignmentRight
	StretchedHorizontally
)

type VerticalAlignment int

const (
	AlignmentTop VerticalAlignment = iota
	AlignmentCenteredVertically
	AlignmentBottom
	StretchedVertically
)
