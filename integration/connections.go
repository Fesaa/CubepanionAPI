package integration

import (
	"fmt"
	"log/slog"

	"github.com/Fesaa/CubepanionAPI/models"
	ws "github.com/gofiber/contrib/websocket"
)

var (
	errWrongMessageType []byte = ws.FormatCloseMessage(ws.CloseUnsupportedData, "Unsupported message type")
)

type Connection struct {
	Uuid string
	c    *ws.Conn

	holder models.Holder
}

func (c *Connection) Start() {
	var (
		mt  int
		msg []byte
		err error
	)
wsLoop:
	for {
		mt, msg, err = c.ReadMessage()
		if err != nil {
			slog.Error(fmt.Sprintf("error reading message: %v - %d", err, mt))
			break
		}

		switch mt {
		case ws.TextMessage:
			err = c.WriteMessage(ws.CloseMessage, errWrongMessageType)
		case ws.BinaryMessage:
			err = c.handlePacket(mt, msg)
		case ws.PingMessage:
			err = c.WriteMessage(ws.PongMessage, msg)
		case ws.PongMessage:
		case ws.CloseMessage:
			_ = c.WriteMessage(mt, msg)
			break wsLoop
		}

		if err != nil {
			slog.Error(fmt.Sprintf("error writing message: %v", err))
			break
		}
	}
}

func (c *Connection) ReadMessage() (int, []byte, error) {
	return c.c.ReadMessage()
}

func (c *Connection) WriteMessage(mt int, msg []byte) error {
	return c.c.WriteMessage(mt, msg)
}

func (c *Connection) Close() error {
	return c.c.Close()
}
