package v2

import (
	"os"
	"os/signal"
)

func OnSign(fn func(os.Signal), signList ...os.Signal) {
	if fn == nil {
		return
	}

	signChan := make(chan os.Signal)
	signal.Notify(signChan, signList...)
	go func() {
		for {
			fn(<-signChan)
		}
	}()
}
