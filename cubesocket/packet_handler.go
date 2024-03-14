package main

import (
	"fmt"
	"log/slog"

	"github.com/Fesaa/CubepanionAPI/core/models"
	"github.com/Fesaa/CubepanionAPI/core/proto/packets"
	"google.golang.org/protobuf/proto"
)

func (c *Client) handlePacket(mt int, msg []byte) error {
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
		slog.Error(fmt.Sprintf("unsupported packet type: %T", packet.Packet))
	}

	return err
}

func (c *Client) handlePing(packet *packets.C2SPingPacket) error {
	var out = packets.S2CPacket{
		Packet: &packets.S2CPacket_Ping{
			Ping: &packets.S2CPingPacket{},
		},
	}

	return c.WritePacket(&out)
}

func (c *Client) handleDisconnect(packet *packets.C2SDisconnectPacket) error {
	c.cleanup()
	return nil
}

func (c *Client) handleUpdateLocation(packet *packets.C2SUpdateLocationPacket) error {
	err := c.ms.DB().SetPlayerLocation(c.UUID, models.Location{
		Current:    packet.Destination,
		Previous:   packet.Origin,
		InPreLobby: packet.PreLobby,
	})
	if err != nil {
		slog.Error("Unable to set player location", "uuid", c.UUID, "error", err)
		// The above failing means that the server isn't aware of the clients correct location.
		// This may lead to incorrect information being sent. So we disconnect the client.
		return fmt.Errorf("Unable to set player location.")
	}
	return nil
}

func (c *Client) handleUpdatePerk(packet *packets.C2SPerkUpdatePacket) error {
	var p = packets.S2CPerkUpdatePacket{
		Category: packet.Category,
		Perks:    packet.Perks,
		Uuid:     c.UUID,
	}

	var out = packets.S2CPacket{
		Packet: &packets.S2CPacket_UpdatePerk{
			UpdatePerk: &p,
		},
	}

	players, err := c.ms.DB().GetSharedPlayers(c.UUID)
	if err != nil {
		slog.Error("Failed to get shared players", "uuid", c.UUID, "error", err)
		// Returning nil here because we don't want to disconnect the client
		// The error has nothing to do with the client, or the connection
		return nil
	}

	for _, player := range players {
		mux.RLock()
		conn, ok := clients[player]
		mux.RUnlock()

		if ok {
			err = conn.WritePacket(&out)
			if err != nil {
				slog.Error("Failed to send perk update", "s_uuid", c.UUID, "rec_uuid", player, "error", err)
			}
		}
	}
	return nil
}
