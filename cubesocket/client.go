package main

import (
	"fmt"
	"log/slog"

	"github.com/Fesaa/CubepanionAPI/core"
	"github.com/Fesaa/CubepanionAPI/core/proto/packets"
	"github.com/Fesaa/CubepanionAPI/cubesocket/database"
	ws "github.com/gofiber/contrib/websocket"
	"google.golang.org/protobuf/proto"
)

var (
	errWrongMessageType []byte = ws.FormatCloseMessage(ws.CloseUnsupportedData, "Unsupported message type")
)

type Client struct {
	UUID string
	c    *ws.Conn
	ms   core.MicroService[core.MicroServiceConfig, database.Database]
}

func (c *Client) Start() {
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
			c.cleanup()
			slog.Error("Error in client loop", "uuid", c.UUID, "error", err)
			break
		}
	}

	slog.Info("Client closed", "uuid", c.UUID)
}

func (c *Client) cleanup() {
	err := c.Close()
	if err != nil {
		slog.Error(fmt.Sprintf("error closing connection: %v", err))
	}
	err = c.ms.DB().RemovePlayerLocation(c.UUID)
	if err != nil {
		slog.Error(fmt.Sprintf("error removing player location: %v", err))
	}
	mux.Lock()
	delete(clients, c.UUID)
	mux.Unlock()
}

func (c *Client) ReadMessage() (int, []byte, error) {
	return c.c.ReadMessage()
}

func (c *Client) WriteMessage(mt int, msg []byte) error {
	return c.c.WriteMessage(mt, msg)
}

func (c *Client) WritePacket(packet *packets.S2CPacket) error {
	bytes, err := proto.Marshal(packet)
	if err != nil {
		return err
	}
	return c.WriteMessage(ws.BinaryMessage, bytes)
}

func (c *Client) Close() error {
	return c.c.Close()
}
