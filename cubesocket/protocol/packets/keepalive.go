package packets

import (
	"github.com/Fesaa/CubepanionAPI/cubesocket/protocol/buf"
	"github.com/go-netty/go-netty"
)

type PacketPing struct{}

func (p *PacketPing) Read(buf buf.PacketBuffer) {
}

func (p *PacketPing) Write(buf buf.PacketBuffer) {
}

func (p *PacketPing) Handle(ctx netty.InboundContext, handler Handler) error {
	return handler.HandlePing(ctx, p)
}

func (p *PacketPing) ID() uint8 {
	return 0
}

func (p *PacketPing) Name() string {
	return "PacketPing"
}

type PacketPong struct{}

func (p *PacketPong) Read(buf buf.PacketBuffer) {
}

func (p *PacketPong) Write(buf buf.PacketBuffer) {
}

func (p *PacketPong) Handle(ctx netty.InboundContext, handler Handler) error {
	return nil
}

func (p *PacketPong) ID() uint8 {
	return 1
}

func (p *PacketPong) Name() string {
	return "PacketPong"
}
