package packets

import (
	"github.com/Fesaa/CubepanionAPI/core/models"
	"github.com/Fesaa/CubepanionAPI/cubesocket/protocol/buf"
	"github.com/go-netty/go-netty"
)

type PacketGameStatUpdate struct {
	game          string
	playerCount   int
	unixTimeStamp int64
}

func (p *PacketGameStatUpdate) GameStat() models.GameStat {
	return models.GameStat{
		Game:          p.game,
		PlayerCount:   p.playerCount,
		UnixTimeStamp: p.unixTimeStamp,
	}
}

func (p *PacketGameStatUpdate) Read(buf buf.PacketBuffer) {
	p.game = buf.ReadString()
	p.playerCount = buf.ReadInt()
	p.unixTimeStamp = buf.ReadLong()
}

func (p *PacketGameStatUpdate) Write(buf buf.PacketBuffer) {}

func (p *PacketGameStatUpdate) Handle(ctx netty.InboundContext, handler Handler) error {
	return handler.HandleGameStatUpdate(ctx, p)
}

func (p *PacketGameStatUpdate) ID() uint8 {
	return 10
}

func (p *PacketGameStatUpdate) Name() string {
	return "GameStatUpdate"
}
