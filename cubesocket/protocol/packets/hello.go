package packets

import (
	"github.com/Fesaa/CubepanionAPI/cubesocket/protocol/buf"
	"github.com/go-netty/go-netty"
)

type PacketHelloPing struct {
	timestamp int64
}

func (p *PacketHelloPing) Timestamp() int64 {
	return p.timestamp
}

func (p *PacketHelloPing) Read(buf buf.PacketBuffer) {
	buf.ReadableBytes()
	p.timestamp = buf.ReadLong()
}

func (p *PacketHelloPing) Write(buf buf.PacketBuffer) {
}

func (p *PacketHelloPing) Handle(ctx netty.InboundContext, handler Handler) error {
	return handler.HandleHelloPing(ctx, p)
}

func (p *PacketHelloPing) ID() uint8 {
	return 2
}

func (p *PacketHelloPing) Name() string {
	return "PacketHelloPing"
}

type PacketHelloPong struct {
}

func (p *PacketHelloPong) Read(buf buf.PacketBuffer) {
}

func (p *PacketHelloPong) Write(buf buf.PacketBuffer) {
}

func (p *PacketHelloPong) Handle(ctx netty.InboundContext, handler Handler) error {
	return nil
}

func (p *PacketHelloPong) ID() uint8 {
	return 3
}

func (p *PacketHelloPong) Name() string {
	return "PacketHelloPong"
}
