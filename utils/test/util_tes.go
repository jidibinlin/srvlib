package main

import (
	"fmt"
	"github.com/gzjjyz/srvlib/utils"
)

func main() {
	utils.ParseCmdInput()
	fmt.Println(utils.IsDev())
}
