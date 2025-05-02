package rest

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v3"
	"gomod.usaken.org/uvcounter/spine"
	"time"
)

func RunServer() {
	server := fiber.New()

	server.Get("/", func(c fiber.Ctx) error {
		return c.SendString("healthy")
	})

	go func() {
		// start server
		spine.SystemGroup.Add(1)
		defer spine.SystemGroup.Done()

		reason := <-spine.C.Done()
		fmt.Printf("api server shutdown started due to %s\n", reason)
		// 5분은 휴리스틱하게 정해진 시간이다.
		// API서버를 내리기전에 이미 충분하게 요청이 들어오지 않은 상태이겠지만,
		// 혹시 5분이상 실행중인 요청이 있다면 실패하도록한다.
		// timeout값보다 크게 하여 최대한 보수적으로 잡는다.
		context.WithTimeout(context.Background(), time.Minute*5)
		err := server.Shutdown()
		if err != nil {
			fmt.Printf("api server shutdown failed %e\n", err)
		}
		fmt.Printf("api server successfully shutdown\n")
	}()

	go func() {
		// add graceful shutdown hook
		spine.SystemGroup.Add(1)
		defer spine.SystemGroup.Done()

		err := server.Listen(":8080")
		if err != nil {
			err = fmt.Errorf("api server run failed: %e", err)
			spine.Cancel(err)
			fmt.Printf("api server error: %e", err)
		}
		fmt.Printf("api server shutdown end\n")
	}()
}
