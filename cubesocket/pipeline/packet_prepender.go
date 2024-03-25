package pipeline

import (
	"bytes"
	"fmt"

	"github.com/Fesaa/CubepanionAPI/cubesocket/protocol/buf"
	"github.com/go-netty/go-netty"
)

type PacketPrepender struct{}

func (p *PacketPrepender) HandleWrite(ctx netty.OutboundContext, msg netty.Message) {
	packetBuffer := msg.(buf.PacketBuffer)
	length := packetBuffer.ReadableBytes()
	varInt := buf.GetVarIntSize(length)
	if varInt > 3 {
		panic(fmt.Errorf("Packet size is too big: %d. Cannot fit %d in 3 bytes", length, varInt))
	}

	buffer := bytes.NewBuffer(nil)
	buffer.Grow(length + varInt)
	buf.WriteVarInt(buffer, length)
	buffer.Write(packetBuffer.GetBytes())
	ctx.HandleWrite(buffer.Bytes())
}
