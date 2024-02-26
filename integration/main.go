package integration

import (
	"github.com/Fesaa/CubepanionAPI/models"
	ws "github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
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
		holder, _ := c.Locals(models.HOLDER_KEY).(models.Holder)
		conn := &Connection{
			c:      c,
			Uuid:   c.Params("uuid"),
			holder: holder,
		}

		conn.Start()
	})
}
