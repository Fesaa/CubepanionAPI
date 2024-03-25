package protocol

import (
	"log/slog"

	"github.com/Fesaa/CubepanionAPI/cubesocket/database"
	"github.com/Fesaa/CubepanionAPI/cubesocket/protocol/packets"
	"github.com/go-netty/go-netty"
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
	conn, ok := clients.Get(ctx.Channel().ID())
	if ok {
		if err := h.db.RemovePlayerLocation(conn.UUID()); err != nil {
			slog.Warn("Unable to remove player location", "uuid", conn.UUID(), "error", err)
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
		conn, ok := clients.Get(ctx.Channel().ID())
		if ok {
			h.db.RemovePlayerLocation(conn.UUID())
		}

		clients.Remove(ctx.Channel().ID())
		idMapping.RemoveByValue(ctx.Channel().ID())
		ctx.Close(err)
	}
}
