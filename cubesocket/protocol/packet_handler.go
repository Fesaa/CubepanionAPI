package protocol

import (
	"github.com/Fesaa/CubepanionAPI/cubesocket/database"
	"github.com/Fesaa/CubepanionAPI/cubesocket/prometheus"
	"github.com/Fesaa/CubepanionAPI/cubesocket/protocol/packets"
	"github.com/go-netty/go-netty"
	"log/slog"
	"time"
)

type PacketHandler struct {
	db database.Database
}

func NewPacketHandler(db database.Database) *PacketHandler {
	return &PacketHandler{
		db: db,
	}
}

func (h *PacketHandler) HandleActive(ctx netty.ActiveContext) {
}

func (h *PacketHandler) HandleInactive(ctx netty.InactiveContext, ex netty.Exception) {
	slog.Info(format(FORMAT_DISCONNECT, ctx.Channel()), "id", ctx.Channel().ID(), "reason", ex)
	conn, ok := clients.Get(ctx.Channel().ID())
	if ok {

		if conn.protocol != 0 {
			prometheus.VersionDec(conn.protocol)
		}

		if conn.state == CONNECTED {
			prometheus.EndSession()
			if err := h.db.RemovePlayerLocation(conn.UUID()); err != nil {
				slog.Warn("Unable to remove player location", "uuid", conn.UUID(), "error", err)
			}
			prometheus.SessionDuration(time.Since(conn.start))
		}
	}

	clients.Remove(ctx.Channel().ID())
	idMapping.RemoveByValue(ctx.Channel().ID())
}

func (h *PacketHandler) HandleRead(ctx netty.InboundContext, msg netty.Message) {
	h.Handle(ctx, msg.(packets.Packet))
}

func (h *PacketHandler) Handle(ctx netty.InboundContext, packet packets.Packet) {
	if err := packet.Handle(ctx, h); err != nil {
		ctx.Close(err)
	}
}
