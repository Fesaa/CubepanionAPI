package protocol

import (
	"time"

	"github.com/Fesaa/CubepanionAPI/cubesocket/utils"
	"github.com/go-netty/go-netty"
)

var (
	clients   utils.Map[int64, Connection]   = utils.NewMap[int64, Connection]()
	idMapping utils.DoubleMap[string, int64] = utils.NewDoubleMap[string, int64]()
)

type Connection struct {
	ctx   netty.HandlerContext
	uuid  string
	start time.Time
}

func (c *Connection) Ctx() netty.HandlerContext {
	return c.ctx
}

func (c *Connection) UUID() string {
	return c.uuid
}
