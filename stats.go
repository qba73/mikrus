package mikrus

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Stats struct {
	Memory    Memory    `json:"memory"`
	DiskSpace string    `json:"disk_space"`
	Uptime    string    `json:"uptime"`
	Processes []Process `json:"processes"`
}

type Memory struct {
	Total     int `json:"total"`
	Used      int `json:"used"`
	Free      int `json:"free"`
	Shared    int `json:"shared"`
	Cache     int `json:"cache"`
	Available int `json:"available"`
	SwapTotal int `json:"swap_total"`
	SwapUsed  int `json:"swap_used"`
	SwapFree  int `json:"swap_free"`
}

func ParseMemoryUsage(s string) (Memory, error) {
	var (
		total, used, free, shared, cache, available int
		swapTotal, swapUsed, swapFree               int
	)

	table := strings.Split(s, "\n")
	if len(table) < 3 {
		return Memory{}, errors.New("parsing `free` command output")
	}
	memory := strings.TrimSpace(table[1])
	swap := strings.TrimSpace(table[2])

	_, err := fmt.Sscanf(memory, "Mem: %d %d %d %d %d %d", &total, &used, &free, &shared, &cache, &available)
	if err != nil {
		return Memory{}, fmt.Errorf("incorrect input data for `free` command: %w", err)
	}
	_, err = fmt.Sscanf(swap, "Swap: %d %d %d", &swapTotal, &swapUsed, &swapFree)
	if err != nil {
		return Memory{}, fmt.Errorf("incorrect input data for `free` command: %w", err)
	}
	return Memory{
		Total:     total,
		Used:      used,
		Free:      free,
		Shared:    shared,
		Cache:     cache,
		Available: available,
		SwapFree:  swapFree,
		SwapUsed:  swapUsed,
		SwapTotal: swapTotal,
	}, nil
}

type DiskSpace struct {
	Filesystem string `json:"filesystem"`
	Size       string `json:"size"`
	Used       string `json:"used"`
	Available  string `json:"available"`
	Usage      string `json:"usage"`
	MountedOn  string `json:"mounted_on"`
}

type ProcessInfo struct {
	User              string
	PID               uint64
	CPUPercent        float64
	MemoryPercent     float64
	VirtualMemorySize uint64
	ResidentSetSize   uint64
	TTY               string
	State             string
	Start             string
	CPUTime           string
	Command           string
}

// Format: USER  PID     %CPU       %MEM         VSZ    RSS    TTY  STAT   START   TIME      COMMAND
var psRE = regexp.MustCompile(`^(\w+) +(\d+) +(\d+\.\d+) +(\d+\.\d+) +(\d+) +(\d+) +(.) +(\w+) +([\d:]+) +([\d:]+) +(.+)$`)

func ParsePS(s string) ([]ProcessInfo, error) {
	lines := strings.Split(s, "\n")
	list := make([]ProcessInfo, 0, len(lines)-1)
	const (
		USER = iota + 1
		PID
		CPUPERCENT
		MEMPERCENT
		VSZ
		RSS
		TTY
		STAT
		START
		TIME
		COMMAND
	)
	for _, line := range lines[1 : len(lines)-1] {
		matches := psRE.FindStringSubmatch(line)
		// if len(matches) != 12 {
		// 	fmt.Printf("%#v)\n", matches)
		// 	return nil, fmt.Errorf("parsing %q", line)
		// }
		pid, err := strconv.ParseUint(matches[PID], 10, 64)
		if err != nil {
			return nil, err
		}
		cpuPercent, err := strconv.ParseFloat(matches[CPUPERCENT], 64)
		if err != nil {
			return nil, err
		}
		memPercent, err := strconv.ParseFloat(matches[MEMPERCENT], 64)
		if err != nil {
			return nil, err
		}
		vsz, err := strconv.ParseUint(matches[VSZ], 10, 64)
		if err != nil {
			return nil, err
		}
		rss, err := strconv.ParseUint(matches[RSS], 10, 64)
		if err != nil {
			return nil, err
		}

		list = append(list, ProcessInfo{
			User:              matches[USER],
			PID:               pid,
			CPUPercent:        cpuPercent,
			MemoryPercent:     memPercent,
			VirtualMemorySize: vsz,
			ResidentSetSize:   rss,
			TTY:               matches[TTY],
			State:             matches[STAT],
			Start:             matches[START],
			CPUTime:           matches[TIME],
			Command:           matches[COMMAND],
		})
	}
	return list, nil
}

func ParseDiskSpace(s string) (DiskSpace, error) {
	var fileSystem, size, used, avail, usage, mountedOn string

	table := strings.Split(s, "\n")
	if len(table) < 2 {
		return DiskSpace{}, errors.New("parsing `df` command output")
	}
	values := strings.TrimSpace(table[1])

	_, err := fmt.Sscanf(values, "%s %s %s %s %s %s", &fileSystem, &size, &used, &avail, &usage, &mountedOn)
	if err != nil {
		return DiskSpace{}, fmt.Errorf("incorrect input data for `df` command: %w", err)
	}

	return DiskSpace{
		Filesystem: fileSystem,
		Size:       size,
		Used:       used,
		Available:  avail,
		Usage:      usage,
		MountedOn:  mountedOn,
	}, nil
}

type Uptime struct {
	Time         string        `json:"time"`
	Uptime       time.Duration `json:"days_up"`
	Users        int           `json:"users_logged_in"`
	CPUload1min  float64       `json:"load_average_1_min"`
	CPUload5min  float64       `json:"load_average_5_min"`
	CPUload15min float64       `json:"load_average_15_min"`
}

var uptimeRE = regexp.MustCompile(`up (\d+) days, +(\d+):(\d\d), +(\d+) users, +load average: (\d+\.\d\d), (\d+\.\d\d), (\d+\.\d\d)`)

func ParseUptime(s string) (Uptime, error) {
	matches := uptimeRE.FindStringSubmatch(s)
	if len(matches) != 8 {
		return Uptime{}, fmt.Errorf("parsing input %q", s)
	}
	const (
		UPDAYS = iota + 1
		UPHOURS
		UPMINUTES
		USERS
		LOAD1MIN
		LOAD5MIN
		LOAD15MIN
	)
	upDays, err := strconv.Atoi(matches[UPDAYS])
	if err != nil {
		return Uptime{}, err
	}
	upHours, err := strconv.Atoi(matches[UPHOURS])
	if err != nil {
		return Uptime{}, err
	}
	upMinutes, err := strconv.Atoi(matches[UPMINUTES])
	if err != nil {
		return Uptime{}, err
	}
	up := 24 * time.Hour * time.Duration(upDays)
	up += time.Duration(upHours) * time.Hour
	up += time.Duration(upMinutes) * time.Minute
	users, err := strconv.Atoi(matches[USERS])
	if err != nil {
		return Uptime{}, err
	}
	load1min, err := strconv.ParseFloat(matches[LOAD1MIN], 64)
	if err != nil {
		return Uptime{}, err
	}
	load5min, err := strconv.ParseFloat(matches[LOAD5MIN], 64)
	if err != nil {
		return Uptime{}, err
	}
	load15min, err := strconv.ParseFloat(matches[LOAD15MIN], 64)
	if err != nil {
		return Uptime{}, err
	}
	return Uptime{
		Uptime:       up,
		Users:        users,
		CPUload1min:  load1min,
		CPUload5min:  load5min,
		CPUload15min: load15min,
	}, nil
}

type Process struct {
	User    string `json:"user"`
	PID     string `json:"pid"`
	CPU     string `json:"cpu"`
	Memory  string `json:"memory"`
	VSZ     string `json:"vsz"`
	RSS     string `json:"rss"`
	TTY     string `json:"tty"`
	Stat    string `json:"stat"`
	Start   string `json:"start"`
	Time    string `json:"time"`
	Command string `json:"command"`
}

// func ParseProcess(s string) (_ Process, err error) {
// 	defer func() {
// 		if err != nil {
// 			err = fmt.Errorf("parsing `ps -u` command output: %w", err)
// 		}
// 	}()

// 	var (
// 		user, pid, cpu, mem string
// 	)

// 	line := strings.TrimSpace(s)
// 	chunks := strings.Split(line, "")
// 	if len(chunks) != 11 {
// 		return Process{}, err
// 	}

// 	return Process{}, nil

// }
