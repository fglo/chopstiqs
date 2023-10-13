package debug

var (
	Debug = false
)

func TurnDebugOn() {
	Debug = true
}

func TurnDebugOff() {
	Debug = false
}
