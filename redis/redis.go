package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"net"
	"syscall"
	"time"
)

var DefaultClient *redis.Client

func init() {
	DefaultClient = Client()
}

func Client() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:       "0.0.0.0:6379",
		ClientName: "uvcounter",
		Dialer: func(ctx context.Context, network, addr string) (net.Conn, error) {
			dialer := &net.Dialer{
				Timeout:   time.Second,
				KeepAlive: time.Minute,
				Control: func(network, address string, c syscall.RawConn) error {
					var controlErr error
					err := c.Control(func(fd uintptr) {
						controlErr = syscall.SetsockoptInt(int(fd), syscall.IPPROTO_TCP, syscall.TCP_NODELAY, 1)
						if controlErr != nil {
							return
						}
						controlErr = syscall.SetsockoptInt(int(fd), syscall.IPPROTO_TCP, 12, 1)
					})
					if err != nil {
						return err
					}
					return controlErr
				},
			}
			return dialer.DialContext(ctx, network, addr)

		},
		OnConnect:             nil,
		Protocol:              3,
		MaxRetries:            3,
		DialTimeout:           time.Second,
		ReadTimeout:           time.Second,
		WriteTimeout:          time.Second,
		ContextTimeoutEnabled: true,
		Limiter:               nil,
	})
	return rdb
}
