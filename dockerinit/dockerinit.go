package main

import (
	_ "example.com/m/v2/daemon/execdriver/lxc"
	_ "example.com/m/v2/daemon/execdriver/native"
	"example.com/m/v2/reexec"
)

func main() {
	// Running in init mode
	reexec.Init()
}
