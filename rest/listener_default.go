//go:build !linux

package rest

import (
	"log"
	"net"
)

func Listener(addr string) net.Listener {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}

	return ln
}
