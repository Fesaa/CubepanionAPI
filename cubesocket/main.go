package main

import (
	"log/slog"

	"github.com/Fesaa/CubepanionAPI/core"
	"github.com/Fesaa/CubepanionAPI/cubesocket/database"
)

func main() {
	config, err := core.LoadDefaultConfig("config.yaml")
	if err != nil {
		slog.Error("Error loading config: ", "error", err)
		return
	}

	ms, err := core.NewMicroService(config, database.Connect)
	if err != nil {
		slog.Error("Error creating microservice: ", "error", err)
		return
	}

	ms.UseDefaults()

	ms.Use("/ws")

	err = ms.Start()
	if err != nil {
		slog.Error("Error starting microservice: ", "error", err)
	}
}
