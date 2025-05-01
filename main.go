package main

import (
	"fmt"
	"gomod.usaken.org/uvcounter/monitor"
	"gomod.usaken.org/uvcounter/spine"
)

func main() {
	fmt.Println("impression counter system started.")

	monitor.RunPprofServer()
	spine.WaitUntilSystemShutdown()
}
