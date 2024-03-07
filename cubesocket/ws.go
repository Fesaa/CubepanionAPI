package main

import (
	"sync"

	ws "github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

var (
	clients map[string]*Client = make(map[string]*Client)
	mux     *sync.RWMutex      = &sync.RWMutex{}
)

func RequireUpgrade(c *fiber.Ctx) error {
	if ws.IsWebSocketUpgrade(c) {
		c.Locals("allowed", true)
		return c.Next()
	}

	return c.SendStatus(fiber.StatusUpgradeRequired)
}

func Handler() fiber.Handler {
	return ws.New(func(c *ws.Conn) {
		uuid := utils.CopyString("uuid")
		client := &Client{
			UUID: uuid,
			c:    c,
		}

		mux.Lock()
		clients[uuid] = client
		mux.Unlock()

		client.Start()
	})
}
