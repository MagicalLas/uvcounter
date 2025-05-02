package main

import (
	"fmt"
	"gomod.usaken.org/uvcounter/monitor"
	"gomod.usaken.org/uvcounter/spine"
)

func main() {
	fmt.Println("uv counter system started.")

	monitor.RunPprofServer()
	monitor.RunPrometheusServer()

	spine.WaitUntilSystemShutdown()

	fmt.Println("uv counter system ended.")
}
