/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log/slog"

	"github.com/paulsavides/palworld-manager/internal/clients"
	"github.com/paulsavides/palworld-manager/internal/runners/monitor"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// monitorCmd represents the monitor command
var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Monitors palworld server and restarts if needed",
	// 	Long: `A longer description that spans multiple lines and likely contains examples
	// and usage of using your command. For example:

	// Cobra is a CLI library for Go that empowers applications.
	// This application is a tool to generate the needed files
	// to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		service, _ := cmd.Flags().GetString("service-name")
		memoryThreshold, _ := cmd.Flags().GetInt("memory-threshold")

		opts := monitor.MonitorOptions{
			RconClient: clients.RconClientOptions{
				Host:     viper.GetString("rcon.host"),
				Port:     viper.GetString("rcon.port"),
				Password: viper.GetString("rcon.password"),
			},
			ServiceName:     service,
			MemoryThreshold: memoryThreshold,
		}

		err := monitor.Execute(opts)
		if err != nil {
			slog.Default().Error("Monitor command failed", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(monitorCmd)

	monitorCmd.Flags().StringP("service-name", "s", "palworld", "The name of the systemd service")
	monitorCmd.Flags().IntP("memory-threshold", "m", 95, "Server memory threshold which will trigger a restart")
	monitorCmd.Flags().BoolP("broadcast-countdown", "b", true, "Broadcast countdown using rcon (requires rcon connection to be set in config)")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// monitorCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// monitorCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
