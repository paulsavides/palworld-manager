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
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		requireBroadcast, _ := cmd.Flags().GetBool("require-broadcast")

		opts := monitor.MonitorOptions{
			RconClient: clients.RconClientOptions{
				Host:     viper.GetString("rcon.host"),
				Port:     viper.GetString("rcon.port"),
				Password: viper.GetString("rcon.password"),
			},
			ServiceName:      service,
			MemoryThreshold:  memoryThreshold,
			DryRun:           dryRun,
			RequireBroadcast: requireBroadcast,
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
	monitorCmd.Flags().Bool("dry-run", false, "Just log intentions, don't actually run restarts")
	monitorCmd.Flags().Bool("require-broadcast", false, "Fail script if rcon broadcast of restart message fails")
}
