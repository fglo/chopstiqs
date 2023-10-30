package debug

var (
	Debug                = false
	ShowComponentBorders = false
	ShowComponentPadding = false
)

func TurnDebugOn() {
	Debug = true
}

func TurnDebugOff() {
	Debug = false
}
