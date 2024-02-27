package integration

import (
	"sync"

	"github.com/Fesaa/CubepanionAPI/models"
	ws "github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

var (
	connections map[string]*Connection = make(map[string]*Connection)
	mux         *sync.RWMutex          = &sync.RWMutex{}
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

		uuid := utils.CopyString(c.Params("uuid"))
		conn := &Connection{
			c:      c,
			Uuid:   uuid,
			holder: holder,
		}
		mux.Lock()
		connections[uuid] = conn
		mux.Unlock()

		conn.Start()
	})
}
