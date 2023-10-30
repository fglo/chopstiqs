package debug

var (
	Debug                = false
	ShowComponentBorders = false
	ShowComponentPadding = false
	ShowGridCells        = false
)

func TurnDebugOn() {
	Debug = true
}

func TurnDebugOff() {
	Debug = false
}
