package monitor

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"runtime"
)

var EnableProfiling = true

func RunPprofServer() {
	if !EnableProfiling {
		return
	}

	runtime.SetBlockProfileRate(5)

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
}
