package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v3"

	"gomod.usaken.org/uvcounter/application"
)

func GetUVCounter(c fiber.Ctx) error {
	counterID := c.Get("counterID", "0")

	data := make([]int, 10)
	service := application.UVCounterService{}
	count := service.GetUVCounter(counterID)
	count += int64(len(data))

	return c.Send([]byte(fmt.Sprintf("count: %d", count)))
}
