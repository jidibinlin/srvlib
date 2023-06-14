package utils

import (
	"fmt"

	"github.com/gzjjyz/srvlib/logger"
)

func SafeLogErr(err error, printWhileLoggerNoReady bool) {
	if !logger.HasInit() && printWhileLoggerNoReady {
		fmt.Println(err)
		return
	}

	logger.Errorf(err.Error())
}

func SafeLogWarn(printWhileLoggerNoReady bool, format string, args ...interface{}) {
	if !logger.HasInit() && printWhileLoggerNoReady {
		fmt.Printf(format, args...)
		return
	}

	logger.Warn(format, args...)
}
