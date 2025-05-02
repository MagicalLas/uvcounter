package handlers

import (
	"github.com/gofiber/fiber/v3"
	"gomod.usaken.org/uvcounter/application"
	"strconv"
	"sync"
)

var bufPool = sync.Pool{
	New: func() interface{} {
		b := make([]byte, 0, 64)
		return &b
	},
}

var countPrefix = []byte("count: ")

func GetUVCounter(c fiber.Ctx) error {
	counterID := c.Get("counterID", "0")

	service := application.UVCounterService{}
	count := service.GetUVCounter(counterID)

	bufPtr := bufPool.Get().(*[]byte)
	buf := (*bufPtr)[:0]

	buf = append(buf, countPrefix...)
	buf = strconv.AppendInt(buf, count, 10)

	err := c.Send(buf)

	bufPool.Put(bufPtr)

	return err
}
