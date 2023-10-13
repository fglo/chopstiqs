package component

var (
	DefaultLeftPadding   = 0
	DefaultRightPadding  = 0
	DefaultTopPadding    = 0
	DefaultBottomPadding = 0
)

func SetDefaultPadding(left, right, top, bottom int) {
	DefaultLeftPadding = left
	DefaultRightPadding = right
	DefaultTopPadding = top
	DefaultBottomPadding = bottom
}

func SetDefaultHorizontalPadding(padding int) {
	DefaultLeftPadding = padding
	DefaultRightPadding = padding
}

func SetDefaultVerticalPadding(padding int) {
	DefaultTopPadding = padding
	DefaultBottomPadding = padding
}
