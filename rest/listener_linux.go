//go:build linux

package rest

import (
	"golang.org/x/sys/unix"
	"log"
	"net"
	"syscall"
)

const (
	TCP_QUICKACK = 12
)

func Listener(addr string) net.Listener {
	lc := net.ListenConfig{
		Control: func(network, address string, c syscall.RawConn) error {
			return setSocketOptionsForPerformance(c)
		},
	}

	ln, err := lc.Listen(nil, "tcp", addr)
	if err != nil {
		log.Fatalf("Failed to create listener: %v", err)
	}

	return ln
}

func setSocketOptionsForPerformance(c syscall.RawConn) (err error) {
	// tcp nodalay와 tcp quickack을 활성화합니다.
	c.Control(func(fd uintptr) {
		err = syscall.SetsockoptInt(int(fd), syscall.IPPROTO_TCP, TCP_QUICKACK, 1)
		if err != nil {
			log.Printf("Failed to set TCP_QUICKACK: %v", err)
		}

		err = syscall.SetsockoptInt(int(fd), syscall.IPPROTO_TCP, unix.TCP_NODELAY, 1)
		if err != nil {
			log.Printf("Failed to set TCP_NODELAY: %v", err)
		}
	})
}
