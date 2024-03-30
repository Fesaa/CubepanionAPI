package packets

import (
	"github.com/Fesaa/CubepanionAPI/cubesocket/protocol/buf"
	"github.com/go-netty/go-netty"
)

type PacketDisconnection struct {
	reason string
}

func (p *PacketDisconnection) Read(buf buf.PacketBuffer) {
	p.reason = buf.ReadString()
}

func (p *PacketDisconnection) Write(buf buf.PacketBuffer) {
	buf.WriteString(p.reason)
}

func (p *PacketDisconnection) Handle(ctx netty.InboundContext, handler Handler) error {
	return handler.HandleDisconnection(ctx, p)
}

func (p *PacketDisconnection) Reason() string {
	return p.reason
}

func (p *PacketDisconnection) ID() uint8 {
	return 6
}

func (p *PacketDisconnection) Name() string {
	return "Disconnection"
}

func DisconnectWithReason(reason string) *PacketDisconnection {
	return &PacketDisconnection{
		reason: reason,
	}
}
