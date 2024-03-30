package pipeline

import (
	"fmt"
	"log/slog"

	"github.com/Fesaa/CubepanionAPI/cubesocket/protocol"
	"github.com/Fesaa/CubepanionAPI/cubesocket/protocol/buf"
	"github.com/Fesaa/CubepanionAPI/cubesocket/protocol/packets"
	"github.com/go-netty/go-netty"
)

const (
	FORMAT_UNKNOWN_ID_COUNT_REACHED = "Reached unknown ID count limit. Closing connection."

	UNKNOWN_ID_LIMIT = 5
)

type PacketDecoder struct {
	unknownIDCount int
}

func (p *PacketDecoder) HandleRead(ctx netty.InboundContext, msg netty.Message) {
	buffer := msg.(buf.PacketBuffer)
	id := buffer.ReadVarInt()
	packet := protocol.PacketFromId(int(id))
	if packet == nil {
		slog.Warn("Ignoring packet", "id", id)
		p.unknownIDCount++
		if p.unknownIDCount > UNKNOWN_ID_LIMIT {
			panic(FORMAT_UNKNOWN_ID_COUNT_REACHED)
		}
		return
	}

	if packet.ID() != 0 && packet.ID() != 1 {
		slog.Debug(format(ctx, packet))
	}

	packet.Read(buffer)
	ctx.HandleRead(packet)
}

const formatString = "[%s] [IN] %d %s"

func format(ctx netty.InboundContext, p packets.Packet) string {
	return fmt.Sprintf(formatString, ctx.Channel().RemoteAddr(), p.ID(), p.Name())
}
