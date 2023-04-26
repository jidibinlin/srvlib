package logger

import "testing"

func TestLog(t *testing.T) {
	InitLogger("gamesrv")
	Debug("hello world")
	Flush()
}
