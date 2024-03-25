package packets

import (
	"github.com/Fesaa/CubepanionAPI/cubesocket/protocol/buf"
	"github.com/go-netty/go-netty"
)

type PacketLogin struct {
	uuid string
}

func (p *PacketLogin) UUID() string {
	return p.uuid
}

func (p *PacketLogin) Read(buf buf.PacketBuffer) {
	p.uuid = buf.ReadString()
}

func (p *PacketLogin) Write(buf buf.PacketBuffer) {}

func (p *PacketLogin) Handle(ctx netty.InboundContext, handler Handler) error {
	return handler.HandleLogin(ctx, p)
}

func (p *PacketLogin) ID() uint8 {
	return 7
}

func (p *PacketLogin) Name() string {
	return "Login"
}

type PacketLoginComplete struct{}

func (p *PacketLoginComplete) Read(buf buf.PacketBuffer) {}

func (p *PacketLoginComplete) Write(buf buf.PacketBuffer) {}

func (p *PacketLoginComplete) Handle(ctx netty.InboundContext, handler Handler) error {
	return nil
}

func (p *PacketLoginComplete) ID() uint8 {
	return 8
}

func (p *PacketLoginComplete) Name() string {
	return "LoginComplete"
}
