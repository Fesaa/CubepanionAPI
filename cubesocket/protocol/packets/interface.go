package packets

import (
	"github.com/Fesaa/CubepanionAPI/cubesocket/protocol/buf"
	"github.com/go-netty/go-netty"
)

type Packet interface {
	Read(buf.PacketBuffer)
	Write(buf.PacketBuffer)
	Handle(ctx netty.InboundContext, handler Handler) error
	ID() uint8
	Name() string
}

type Handler interface {
	HandlePing(ctx netty.InboundContext, packet *PacketPing) error
	HandleHelloPing(ctx netty.InboundContext, packet *PacketHelloPing) error
	HandleLocationUpdate(ctx netty.InboundContext, packet *PacketLocationUpdate) error
	HandlePerkUpdate(ctx netty.InboundContext, packet *PacketPerkUpdate) error
	HandleDisconnection(ctx netty.InboundContext, packet *PacketDisconnection) error
	HandleLogin(ctx netty.InboundContext, packet *PacketLogin) error
}
