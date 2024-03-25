package main

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/Fesaa/CubepanionAPI/core"
	"github.com/Fesaa/CubepanionAPI/cubesocket/database"
	"github.com/Fesaa/CubepanionAPI/cubesocket/pipeline"
	"github.com/Fesaa/CubepanionAPI/cubesocket/protocol"
	"github.com/go-netty/go-netty"
)

func main() {
	config, err := core.LoadDefaultConfig("config.yaml")
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		return
	}

	if (len(os.Args) > 2 && os.Args[1] == "--debug") || os.Getenv("DEBUG") == "true" {
		slog.Info("Changing log level to debug")
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	db, err := database.Connect(config.Database())
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		return
	}

	childInitialize := func(channel netty.Channel) {
		channel.Pipeline().
			AddLast(pipeline.EOFFilter{}).
			AddLast(netty.ReadIdleHandler(time.Second * 30)).
			AddLast(netty.WriteIdleHandler(time.Second * 30)).
			AddLast(pipeline.PacketSplitter{}).
			AddLast(pipeline.PacketDecoder{}).
			AddLast(pipeline.PacketPrepender{}).
			AddLast(pipeline.PacketEncoder{}).
			AddLast(protocol.NewPacketHandler(db))
	}

	bootstrap := netty.NewBootstrap(netty.WithChildInitializer(childInitialize))

	address := fmt.Sprintf("%s:%d", config.Host(), config.Port())
	slog.Info("Starting server", "address", address)
	start := time.Now()

	err = bootstrap.Listen(address).Sync()
	if err != nil {
		panic(err)
	}

	select {
	case <-bootstrap.Context().Done():
	}

	slog.Info("Closing server", "uptime", time.Since(start))
}
