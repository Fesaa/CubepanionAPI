package pipeline

import (
	"fmt"
	"log/slog"

	"github.com/Fesaa/CubepanionAPI/cubesocket/protocol/buf"
	"github.com/Fesaa/CubepanionAPI/cubesocket/protocol/packets"
	"github.com/go-netty/go-netty"
	"github.com/go-netty/go-netty/utils"
)

type PacketEncoder struct {
}

func (e *PacketEncoder) HandleWrite(ctx netty.OutboundContext, msg netty.Message) {
	switch packet := msg.(type) {
	case packets.Packet:
		if packet.ID() != 0 && packet.ID() != 1 {
			slog.Debug(encodeFormat(ctx, packet))
		}

		buffer := buf.NewPacketBuffer(nil)
		buffer.WriteVarInt(int(packet.ID()))
		packet.Write(buffer)
		ctx.HandleWrite(buffer)
	default:
		utils.Assert(fmt.Errorf("unexpected message type: %T", packet))
	}
}

const encodeFormatString = "[%s] [OUT] %d %s"

func encodeFormat(ctx netty.OutboundContext, p packets.Packet) string {
	return fmt.Sprintf(encodeFormatString, ctx.Channel().RemoteAddr(), p.ID(), p.Name())
}
