package main

import (
	"fmt"
	"runtime/debug"

	"gomod.usaken.org/uvcounter/monitor"
	"gomod.usaken.org/uvcounter/rest"
	"gomod.usaken.org/uvcounter/spine"
)

func main() {
	debug.SetMemoryLimit(16 << 30)

	fmt.Println("uv counter system started.")

	monitor.RunPprofServer()
	monitor.RunPrometheusServer()

	rest.RunServer()

	spine.WaitUntilSystemShutdown()

	fmt.Println("uv counter system ended.")
}
