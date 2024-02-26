package integration

import (
	"fmt"
	"log/slog"

	"github.com/Fesaa/CubepanionAPI/proto/packets"
	"google.golang.org/protobuf/proto"
)

func (c *Connection) handlePacket(mt int, msg []byte) error {
	packet := packets.C2SPacket{}
	err := proto.Unmarshal(msg, &packet)
	if err != nil {
		return err
	}

	switch packet.Packet.(type) {
	case *packets.C2SPacket_UpdateLocation:
		err = c.handleUpdateLocation(packet.GetUpdateLocation())
	case *packets.C2SPacket_UpdatePerk:
		err = c.handleUpdatePerk(packet.GetUpdatePerk())
	default:
		slog.Debug(fmt.Sprintf("unsupported packet type: %T", packet.Packet))
	}

	return err
}

func (c *Connection) handleUpdateLocation(packet *packets.C2SUpdateLocationPacket) error {
	slog.Info(fmt.Sprintf("update location: %v", packet))
	return nil
}

func (c *Connection) handleUpdatePerk(packet *packets.C2SPerkUpdatePacket) error {
	slog.Info(fmt.Sprintf("update perk: %v", packet))
	return nil
}
