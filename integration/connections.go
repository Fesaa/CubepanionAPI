package integration

import (
	"fmt"
	"log/slog"

	"github.com/Fesaa/CubepanionAPI/models"
	"github.com/Fesaa/CubepanionAPI/proto/packets"
	ws "github.com/gofiber/contrib/websocket"
	"google.golang.org/protobuf/proto"
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
			c.cleanup()
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
		case ws.CloseMessage:
			c.cleanup()
			break wsLoop
		}

		if err != nil {
			slog.Error(fmt.Sprintf("error writing message: %v", err))
			break
		}
	}
}

func (c *Connection) cleanup() {
	err := c.Close()
	if err != nil {
		slog.Error(fmt.Sprintf("error closing connection: %v", err))
	}
	err = c.holder.GetPlayerLocationProvider().RemovePlayerLocation(c.Uuid)
	if err != nil {
		slog.Error(fmt.Sprintf("error removing player location: %v", err))
	}
	mux.Lock()
	delete(connections, c.Uuid)
	mux.Unlock()
}

func (c *Connection) ReadMessage() (int, []byte, error) {
	return c.c.ReadMessage()
}

func (c *Connection) WriteMessage(mt int, msg []byte) error {
	return c.c.WriteMessage(mt, msg)
}

func (c *Connection) WritePacket(packet *packets.S2CPacket) error {
	bytes, err := proto.Marshal(packet)
	if err != nil {
		return err
	}
	return c.WriteMessage(ws.BinaryMessage, bytes)
}

func (c *Connection) Close() error {
	return c.c.Close()
}
