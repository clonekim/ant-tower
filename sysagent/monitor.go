package sysagent

import (
	"context"
	_ "context"
	"os/user"
	_ "syscall"
	"time"
	"twn-monitor/data"

	"github.com/rs/zerolog/log"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/process"
)

// MonitorService : Manages system monitoring logic
type MonitorService struct {
	broadcast chan interface{}
}

// NewMonitorService : Constructor
func NewMonitorService(broadcast chan interface{}) *MonitorService {
	return &MonitorService{
		broadcast,
	}
}

// Start StartMonitor Polling process
func (s *MonitorService) Start(ctx context.Context) {
	runningProcs := make(map[int32]data.ProcessLog)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	log.Info().Msg("MonitorService: Loop started")

	for {
		select {

		case <-ctx.Done():
			// Graceful Shutdown
			log.Info().Msg("MonitorService: Stopping loop (Context cancelled)")
			return

		case <-ticker.C:
			procs, err := process.Processes()
			if err != nil {
				continue
			}

			currentPIDs := make(map[int32]bool)

			for _, p := range procs {
				pid := p.Pid
				currentPIDs[pid] = true
				if _, exists := runningProcs[pid]; !exists {
					name, _ := p.Name()
					if name != "" {
						newLog := data.ProcessLog{
							PID:       pid,
							Name:      name,
							StartTime: getProcessStartTime(p),
							EndTime:   nil,
						}
						runningProcs[pid] = newLog
						// Broadcast with uptime
						uptime, _ := GetUptime()
						s.broadcast <- map[string]interface{}{
							"type":   "START",
							"data":   newLog,
							"uptime": uptime,
						}
					}
				}
			}

			// Detect End
			for pid, pLog := range runningProcs {
				if !currentPIDs[pid] {
					now := time.Now()
					pLog.EndTime = &now
					pLog.Duration = pLog.EndTime.Sub(pLog.StartTime).String()

					// Save to DB
					data.DB.Create(&pLog)
					delete(runningProcs, pid)

					// Broadcast with uptime
					uptime, _ := GetUptime()
					s.broadcast <- map[string]interface{}{
						"type":   "END",
						"data":   pLog,
						"uptime": uptime,
					}
				}
			}
		}
	}
}

func KillProcess(pid int32) error {
	p, err := process.NewProcess(pid)
	if err != nil {
		return err
	}
	return p.Kill()
}

func GetUptime() (uint64, error) {
	info, err := host.Info()
	if err != nil {
		return 0, err
	}
	return info.Uptime, nil
}

func GetProcessSnapshots() ([]data.ProcessLog, error) {
	procs, err := process.Processes()
	if err != nil {
		return nil, err
	}

	var result []data.ProcessLog
	for _, p := range procs {
		name, _ := p.Name()
		if name == "" {
			continue
		}

		result = append(result, data.ProcessLog{
			PID:       p.Pid,
			Name:      name,
			StartTime: getProcessStartTime(p),
			EndTime:   nil,
			Duration:  "",
		})
	}
	return result, nil
}

func getProcessStartTime(p *process.Process) time.Time {
	cTime, err := p.CreateTime() // Unix Milliseconds
	if err != nil {
		return time.Now()
	}
	return time.UnixMilli(cTime)
}

func GetCurrentUser() string {
	currentUser, err := user.Current()
	if err != nil {
		return "unknown"
	}
	return currentUser.Username
}
