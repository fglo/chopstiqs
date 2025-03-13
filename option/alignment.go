package option

type HorizontalAlignment int

const (
	AlignmentLeft HorizontalAlignment = iota
	AlignmentCenteredHorizontally
	AlignmentRight
)

type VerticalAlignment int

const (
	AlignmentTop VerticalAlignment = iota
	AlignmentCenteredVertically
	AlignmentBottom
)
