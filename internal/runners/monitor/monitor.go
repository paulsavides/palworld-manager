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
}

func Execute(options MonitorOptions) error {
	logger := slog.Default()

	rconClient, err := clients.Rcon(options.RconClient)
	if err != nil {
		return err
	}

	defer rconClient.Close()

	restarted, err := startServiceIfNeeded(options.ServiceName)
	if restarted {
		logger.Info("Found service not to be running... starting service now and skip running rest of job")
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
		broadcastRestartWarning(rconClient, int(inuse), options.MemoryThreshold)
		restartService(options.ServiceName)
	}

	return nil
}

func startServiceIfNeeded(serviceName string) (bool, error) {
	res := utility.Shell("systemctl", "status", serviceName)
	if !res.Success {
		return false, res.AsError()
	}

	matched, err := regexp.Match("Active: active \\(running\\)", []byte(res.Stdout))
	if err != nil {
		return false, err
	}

	if !matched {
		utility.Shell("systemctl", "start", "palworld")
		return true, nil
	}

	return false, nil
}

func restartService(serviceName string) error {
	res := utility.Shell("systemctl", "palworld", "restart")
	if !res.Success {
		return res.AsError()
	} else {
		return nil
	}
}

func broadcastRestartWarning(rcon clients.RconClient, memUsage int, memThreshold int) {
	rcon.Broadcast(fmt.Sprintf("[Announcement] Found server memory usage to be %d, above threshold of %d", memUsage, memThreshold))
	time.Sleep(1 * time.Second)
	rcon.Broadcast("Restarting server in 30 seconds")
	time.Sleep(10 * time.Second)
	rcon.Broadcast("Restarting server in 20 seconds")
	time.Sleep(10 * time.Second)
	rcon.Broadcast("Restarting server in 10")
	time.Sleep(1 * time.Second)
	rcon.Broadcast("Restarting server in 9")
	time.Sleep(1 * time.Second)
	rcon.Broadcast("Restarting server in 8")
	time.Sleep(1 * time.Second)
	rcon.Broadcast("Restarting server in 7")
	time.Sleep(1 * time.Second)
	rcon.Broadcast("Restarting server in 6")
	time.Sleep(1 * time.Second)
	rcon.Broadcast("Restarting server in 5")
	time.Sleep(1 * time.Second)
	rcon.Broadcast("Restarting server in 4")
	time.Sleep(1 * time.Second)
	rcon.Broadcast("Restarting server in 3")
	time.Sleep(1 * time.Second)
	rcon.Broadcast("Restarting server in 2")
	time.Sleep(1 * time.Second)
	rcon.Broadcast("Restarting server in 1")
	time.Sleep(1 * time.Second)
	rcon.Broadcast("Triggering restart now")
}
