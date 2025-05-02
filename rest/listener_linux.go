//go:build linux

package rest

import (
	"log"
	"net"
	"syscall"

	"golang.org/x/sys/unix"
)

const (
	TCP_QUICKACK = 12
)

func Listener(addr string) net.Listener {
	ln, err := lc.Listen(nil, "tcp", addr)
	if err != nil {
		log.Fatalf("Failed to create listener: %v", err)
	}

	return &tcpListener{ln}
}

type tcpListener struct {
	net.Listener
}

func (l *tcpListener) Accept() (net.Conn, error) {
	c, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}

	// connection 단에서 setsockopt 적용
	if tcpConn, ok := c.(*net.TCPConn); ok {
		rawConn, err := tcpConn.SyscallConn()
		if err != nil {
			log.Printf("Failed to get raw conn: %v", err)
			return c, err
		}
		setSocketOptionsForPerformance(rawConn)
	}

	return c, nil
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
	return
}
