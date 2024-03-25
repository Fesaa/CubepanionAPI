package packets

import (
	"github.com/Fesaa/CubepanionAPI/core/models"
	"github.com/Fesaa/CubepanionAPI/cubesocket/protocol/buf"
	"github.com/go-netty/go-netty"
)

type PacketLocationUpdate struct {
	origin     string
	dest       string
	inPreLobby bool
}

func (p *PacketLocationUpdate) Location() models.Location {
	return models.Location{
		Current:    p.dest,
		Previous:   p.origin,
		InPreLobby: p.inPreLobby,
	}
}

func (p *PacketLocationUpdate) Read(buf buf.PacketBuffer) {
	p.origin = buf.ReadString()
	p.dest = buf.ReadString()
	p.inPreLobby = buf.ReadBool()
}

func (p *PacketLocationUpdate) Write(buf buf.PacketBuffer) {
}

func (p *PacketLocationUpdate) Handle(ctx netty.InboundContext, handler Handler) error {
	return handler.HandleLocationUpdate(ctx, p)
}

func (p *PacketLocationUpdate) ID() uint8 {
	return 2
}

func (p *PacketLocationUpdate) Name() string {
	return "LocationUpdate"
}
