package pipeline

import (
	"bufio"
	"fmt"

	"github.com/Fesaa/CubepanionAPI/cubesocket/protocol/buf"
	"github.com/go-netty/go-netty"
	"github.com/go-netty/go-netty/utils"
)

const (
	FORMAT_WRONG_LENGTH      = "Expected %d bytes, got %d. Closing connection."
	FORMAT_MAX_PACKET_LENGTH = "Packet length of %d is larger than maximum of %d. Closing connection."

	MAX_PACKET_LENGTH = 2048
)

type PacketSplitter struct {
}

func (p *PacketSplitter) HandleRead(ctx netty.InboundContext, msg netty.Message) {
	reader := utils.MustToReader(msg)
	r := bufio.NewReader(reader)

	length := buf.ReadVarInt(r)
	if length > MAX_PACKET_LENGTH {
		panic(fmt.Sprintf(FORMAT_MAX_PACKET_LENGTH, length, MAX_PACKET_LENGTH))
	}

	bs := make([]byte, length)
	n, err := r.Read(bs)
	utils.Assert(err)
	if n != length {
		panic(fmt.Sprintf(FORMAT_WRONG_LENGTH, length, n))
	}
	ctx.HandleRead(buf.NewPacketBuffer(bs))
}
