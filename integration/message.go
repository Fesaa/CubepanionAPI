package integration

import (
	"fmt"
	"log/slog"

	"github.com/Fesaa/CubepanionAPI/models"
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
	case *packets.C2SPacket_Ping:
		err = c.handlePing(packet.GetPing())
	case *packets.C2SPacket_Disconnect:
		err = c.handleDisconnect(packet.GetDisconnect())
	default:
		slog.Debug(fmt.Sprintf("unsupported packet type: %T", packet.Packet))
	}

	return err
}

func (c *Connection) handlePing(packet *packets.C2SPingPacket) error {
	var out = packets.S2CPacket{
		Packet: &packets.S2CPacket_Ping{
			Ping: &packets.S2CPingPacket{},
		},
	}

	return c.WritePacket(&out)
}

func (c *Connection) handleDisconnect(packet *packets.C2SDisconnectPacket) error {
	c.cleanup()
	return nil
}

func (c *Connection) handleUpdateLocation(packet *packets.C2SUpdateLocationPacket) error {
	c.holder.GetPlayerLocationProvider().SetPlayerLocation(c.Uuid, models.Location{
		Current:    packet.Destination,
		Previous:   packet.Origin,
		InPreLobby: packet.PreLobby,
	})
	return nil
}

func (c *Connection) handleUpdatePerk(packet *packets.C2SPerkUpdatePacket) error {
	var p = packets.S2CPerkUpdatePacket{
		Category: packet.Category,
		Perks:    packet.Perks,
		Uuid:     c.Uuid,
	}

	var out = packets.S2CPacket{
		Packet: &packets.S2CPacket_UpdatePerk{
			UpdatePerk: &p,
		},
	}

	players, err := c.holder.GetPlayerLocationProvider().GetSharedPlayers(c.Uuid)
	if err != nil {
		return err
	}

	for _, player := range players {
		mux.RLock()
		conn, ok := connections[player]
		mux.RUnlock()

		if ok {
			err = conn.WritePacket(&out)
			if err != nil {
				slog.Error(fmt.Sprintf("Failed to send perk update to %v: %v", player, err))
			}
		}
	}
	return nil
}
