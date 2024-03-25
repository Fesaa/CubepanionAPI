package packets

import (
	"github.com/Fesaa/CubepanionAPI/cubesocket/protocol/buf"
	"github.com/go-netty/go-netty"
)

type PacketPerkUpdate struct {
	uuid     string
	perks    []byte
	category string
}

func (p *PacketPerkUpdate) Perks() []byte {
	return p.perks
}

func (p *PacketPerkUpdate) Category() string {
	return p.category
}

func (p *PacketPerkUpdate) UUID() string {
	return p.uuid
}

func (p *PacketPerkUpdate) Read(buf buf.PacketBuffer) {
	p.category = buf.ReadString()
	p.uuid = buf.ReadString()
	perksLength := buf.ReadInt()
	perks := make([]byte, perksLength)
	buf.ReadBytes(perks)
	p.perks = perks
}

func (p *PacketPerkUpdate) Write(buf buf.PacketBuffer) {
	buf.WriteString(p.category)
	buf.WriteString(p.uuid)
	buf.WriteInt(len(p.perks))
	buf.WriteBytes(p.perks)
}

func (p *PacketPerkUpdate) Handle(ctx netty.InboundContext, handler Handler) error {
	return handler.HandlePerkUpdate(ctx, p)
}

func (p *PacketPerkUpdate) ID() uint8 {
	return 5
}

func (p *PacketPerkUpdate) Name() string {
	return "PerkUpdate"
}
