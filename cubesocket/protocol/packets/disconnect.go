package packets

import (
	"github.com/Fesaa/CubepanionAPI/cubesocket/protocol/buf"
	"github.com/go-netty/go-netty"
)

type PacketDisconnection struct{}

func (p *PacketDisconnection) Read(buf buf.PacketBuffer) {
}

func (p *PacketDisconnection) Write(buf buf.PacketBuffer) {
}

func (p *PacketDisconnection) Handle(ctx netty.InboundContext, handler Handler) error {
	return nil
}

func (p *PacketDisconnection) ID() uint8 {
	return 6
}

func (p *PacketDisconnection) Name() string {
	return "Disconnection"
}
