package monitor

import (
	"fmt"
	"log/slog"
	"regexp"
	"time"

	"github.com/paulsavides/palworld-manager/internal/clients"
	"github.com/paulsavides/palworld-manager/internal/utility"
)

type MonitorOptions struct {
	RconClient      clients.RconClientOptions
	MemoryThreshold int
	ServiceName     string
	DryRun          bool
}

func Execute(options MonitorOptions) error {
	logger := slog.Default()

	if options.DryRun {
		logger.Info("Running monitor command with dry run enabled, will not take any actions")
	}

	rconClient, err := clients.Rcon(options.RconClient)
	if err != nil {
		return err
	}

	defer rconClient.Close()

	started, err := startServiceIfNeeded(options.ServiceName, options.DryRun)
	if started {
		logger.Info("Service was started, skipping rest of run")
		return nil
	}

	if err != nil {
		return err
	}

	sysmem := utility.GetSysMemory()
	inuse := (1 - float64(sysmem.MemAvailable)/float64(sysmem.MemTotal)) * 100

	logger.Info(fmt.Sprintf("Memory in use = %f", inuse))

	if inuse > float64(options.MemoryThreshold) {
		logger.Info("Determined that a restart of palworld is needed... restarting now")
		broadcastRestartWarning(rconClient, int(inuse), options.MemoryThreshold, options.DryRun)
		if err := restartService(options.ServiceName, options.DryRun); err != nil {
			return err
		}
	}

	return nil
}

func startServiceIfNeeded(serviceName string, dryRun bool) (bool, error) {
	logger := slog.Default()

	res := utility.Shell("systemctl", "status", serviceName)
	if !res.Success {
		return false, res.AsError()
	}

	matched, err := regexp.Match("Active: active \\(running\\)", []byte(res.Stdout))
	if err != nil {
		return false, err
	}

	if !matched {
		if dryRun {
			logger.Info("Dryrun enabled, will not run start")
		} else {
			logger.Info("Found service not to be running, starting service now")
			utility.Shell("systemctl", "start", "palworld")
		}
		return true, nil
	}

	return false, nil
}

func restartService(serviceName string, dryRun bool) error {
	logger := slog.Default()

	if dryRun {
		logger.Info("Dryrun enabled, will not restart")
		return nil
	}

	res := utility.Shell("systemctl", "restart", "palworld")
	if !res.Success {
		return res.AsError()
	} else {
		return nil
	}
}

func broadcastRestartWarning(rcon clients.RconClient, memUsage int, memThreshold int, dryRun bool) {
	if dryRun {
		rcon.Broadcast("[Announcement] This is a dry run, a restart will not actually be triggered")
		time.Sleep(1 * time.Second)
	}

	rcon.Broadcast(fmt.Sprintf("[Announcement] Found server memory usage to be %d, above threshold of %d", memUsage, memThreshold))
	time.Sleep(1 * time.Second)

	rcon.Broadcast("Restarting server in 30 seconds")
	time.Sleep(10 * time.Second)
	rcon.Broadcast("Restarting server in 20 seconds")
	time.Sleep(10 * time.Second)

	for i := 10; i > 0; i-- {
		rcon.Broadcast(fmt.Sprintf("Restarting server in %d", i))
		time.Sleep(1 * time.Second)
	}

	rcon.Broadcast("Triggering restart now")
}
