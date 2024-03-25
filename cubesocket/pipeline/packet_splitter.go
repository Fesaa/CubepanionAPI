package pipeline

import (
	"bufio"

	"github.com/Fesaa/CubepanionAPI/cubesocket/protocol/buf"
	"github.com/go-netty/go-netty"
	"github.com/go-netty/go-netty/utils"
)

type PacketSplitter struct {
}

func (p *PacketSplitter) HandleRead(ctx netty.InboundContext, msg netty.Message) {
	reader := utils.MustToReader(msg)
	r := bufio.NewReader(reader)

	length := buf.ReadVarInt(r)
	bs := make([]byte, length)
	_, err := r.Read(bs)
	utils.Assert(err)
	ctx.HandleRead(buf.NewPacketBuffer(bs))
}
