/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log/slog"
	"strings"

	"github.com/paulsavides/palworld-manager/internal/clients"
	"github.com/paulsavides/palworld-manager/internal/runners/rcon"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rconCmd represents the rcon command
var rconCmd = &cobra.Command{
	Use:   "rcon",
	Short: "Run rcon commands using connection configured in config file. Uses all arguments as input command",
	Run: func(cmd *cobra.Command, args []string) {
		var command string

		if len(args) == 0 {
			slog.Default().Info("No arguments given, defaulting to Info command")
			command = "Info"
		} else {
			command = strings.Join(args, " ")
		}

		opts := rcon.RconOptions{
			RconClient: clients.RconClientOptions{
				Host:     viper.GetString("rcon.host"),
				Port:     viper.GetString("rcon.port"),
				Password: viper.GetString("rcon.password"),
			},
			Command: command,
		}

		rcon.Execute(opts)
	},
}

func init() {
	rootCmd.AddCommand(rconCmd)
}
