package pipeline

import (
	"io"
	"log/slog"

	"github.com/go-netty/go-netty"
)

type EOFFilter struct{}

func (f *EOFFilter) HandleException(ctx netty.ExceptionContext, ex netty.Exception) {
	if ex == io.EOF {
		slog.Debug("Connection closed by peer, ignoring EOF", "addr", ctx.Channel().RemoteAddr())
		ctx.Close(io.EOF)
		return
	}

	ctx.HandleException(ex)
}
