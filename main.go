package main

import (
	"fmt"

	"gomod.usaken.org/uvcounter/monitor"
	"gomod.usaken.org/uvcounter/rest"
	"gomod.usaken.org/uvcounter/spine"
)

func main() {
	fmt.Println("uv counter system started.")

	ballast := make([]byte, 10<<30)
	_ = ballast

	monitor.RunPprofServer()
	monitor.RunPrometheusServer()

	rest.RunServer()

	spine.WaitUntilSystemShutdown()

	fmt.Println("uv counter system ended.")
}
