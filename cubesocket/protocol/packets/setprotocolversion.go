package packets

import (
	"github.com/Fesaa/CubepanionAPI/cubesocket/protocol/buf"
	"github.com/go-netty/go-netty"
)

type PacketSetProtocolVersion struct {
	protocolVersion int
}

func (p *PacketSetProtocolVersion) ProtocolVersion() int {
	return p.protocolVersion
}

func (p *PacketSetProtocolVersion) Read(buf buf.PacketBuffer) {
	p.protocolVersion = buf.ReadInt()
}

func (p *PacketSetProtocolVersion) Write(buf buf.PacketBuffer) {
	buf.WriteInt(p.protocolVersion)
}

func (p *PacketSetProtocolVersion) Handle(ctx netty.InboundContext, handler Handler) error {
	return handler.HandleSetProtocol(ctx, p)
}

func (p *PacketSetProtocolVersion) ID() uint8 {
	return 9
}

func (p *PacketSetProtocolVersion) Name() string {
	return "SetProtocolVersion"
}
