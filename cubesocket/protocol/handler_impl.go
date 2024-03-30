package protocol

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/Fesaa/CubepanionAPI/cubesocket/prometheus"
	"github.com/Fesaa/CubepanionAPI/cubesocket/protocol/packets"
	"github.com/go-netty/go-netty"
)

const (
	FORMAT_LOGIN      = "[%s] [LOGIN]"
	FORMAT_HELLO      = "[%s] [HELLO]"
	FORMAT_DISCONNECT = "[%s] [DISCONNECT]"
)

func (h *PacketHandler) HandleHelloPing(ctx netty.InboundContext, packet *packets.PacketHelloPing) error {
	slog.Info(format(FORMAT_HELLO, ctx.Channel()), "id", ctx.Channel().ID(), "time", packet.Timestamp())
	ctx.Write(&packets.PacketHelloPong{})

	clients.Put(ctx.Channel().ID(), &Connection{
		ctx:   ctx,
		state: LOGIN,
	})

	go func(ctx netty.InboundContext) {
		<-time.After(5 * time.Second)
		conn := mustConnection(ctx.Channel())
		if conn.state == LOGIN {
			slog.Warn("Client didn't send login packet in time. Disconnecting...", "id", ctx.Channel().ID())
			ctx.Write(packets.DisconnectWithReason("Login timeout"))
			ctx.Close(nil)
		}
	}(ctx)

	return nil
}

func (h *PacketHandler) HandlePing(ctx netty.InboundContext, packet *packets.PacketPing) error {
	ctx.Write(&packets.PacketPong{})
	return nil
}

func (h *PacketHandler) HandlePerkUpdate(ctx netty.InboundContext, packet *packets.PacketPerkUpdate) error {
	conn := mustConnection(ctx.Channel())

	players, err := h.db.GetSharedPlayers(conn.UUID())
	if err != nil {
		slog.Warn("Unable to get shared players", "uuid", conn.UUID(), "error", err)
		// Returning nil here because we don't want to disconnect the client
		// The error has nothing to do with the client, or the connection
		return nil
	}

	for _, player := range players {
		channelId, ok := idMapping.Get(player)
		if !ok {
			slog.Warn("No connection found for player not found, but was in database. Removing...", "uuid", player)
			if err := h.db.RemovePlayerLocation(player); err != nil {
				slog.Error("Unable to remove player location", "uuid", player, "error", err)
			}
			continue
		}

		other, ok := clients.Get(channelId)
		if ok {
			slog.Debug("Sending perk update to player", "uuid", player)
			other.ctx.Write(packet)
		}
	}

	return nil
}

func (h *PacketHandler) HandleDisconnection(ctx netty.InboundContext, packet *packets.PacketDisconnection) error {
	ctx.Close(fmt.Errorf("Client disconnected %s", packet.Reason()))
	return nil
}

func (h *PacketHandler) HandleLocationUpdate(ctx netty.InboundContext, packet *packets.PacketLocationUpdate) error {
	conn := getConnection(ctx.Channel())
	if conn == nil {
		slog.Error("Received location update from unknown player", "id", ctx.Channel().ID())
		return fmt.Errorf("unknown connection")
	}

	err := h.db.SetPlayerLocation(conn.UUID(), packet.Location())
	if err != nil {
		slog.Error("Unable to set player location", "uuid", conn.UUID, "error", err)
		// The above failing means that the server isn't aware of the clients correct location.
		// This may lead to incorrect information being sent. So we disconnect the client.
		return fmt.Errorf("Unable to set player location.")
	}
	return nil
}

func (h *PacketHandler) HandleLogin(ctx netty.InboundContext, packet *packets.PacketLogin) error {
	conn := mustConnection(ctx.Channel())

	defer func() {
		prometheus.NewSessions()
		prometheus.StartSession()
	}()
	slog.Info(format(FORMAT_LOGIN, ctx.Channel()), "id", ctx.Channel().ID(), "uuid", packet.UUID())

	conn.uuid = packet.UUID()
	conn.start = time.Now()
	conn.state = CONNECTED

	idMapping.Put(packet.UUID(), ctx.Channel().ID())

	ctx.Write(&packets.PacketLoginComplete{})
	return nil
}

func format(format string, ch netty.Channel) string {
	return fmt.Sprintf(format, ch.RemoteAddr())
}
