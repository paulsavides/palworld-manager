package rcon

import (
	"log/slog"

	"github.com/paulsavides/palworld-manager/internal/clients"
)

type RconOptions struct {
	RconClient clients.RconClientOptions
	Command    string
}

func Execute(options RconOptions) {
	logger := slog.Default()

	client, err := clients.Rcon(options.RconClient)
	if err != nil {
		logger.Error("Rcon connection not established successfully", "err", err.Error())
		return
	}

	defer client.Close()

	out, err := client.Execute(options.Command)
	if err != nil {
		logger.Error("Rcon command not run successfully", "err", err)
	} else {
		logger.Info(out)
	}
}
