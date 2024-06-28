package mikrus

import (
	"errors"
	"fmt"
	"strings"
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

func ParseDiskSpace(s string) (DiskSpace, error) {
	var (
		fileSystem, size, used, avail, usage, mountedOn string
	)

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
	Time          string `json:"time"`
	DaysUp        int    `json:"days_up"`
	UsersLoggedIn int    `json:"users_logged_in"`
	CPUload1min   string `json:"load_average_1_min"`
	CPUload5min   string `json:"load_average_5_min"`
	CPUload15min  string `json:"load_average_15_min"`
}

func ParseUptime(s string) (_ Uptime, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("parsing `uptime` command output: %w", err)
		}
	}()

	var (
		daysUp, users           int
		time, cpu1, cpu5, cpu15 string
	)

	table := strings.Split(s, "\n")
	if len(table) < 1 {
		return Uptime{}, err
	}

	line := strings.TrimSpace(table[0])
	chunks := strings.Split(line, ",")
	if len(chunks) != 6 {
		return Uptime{}, err
	}

	_, err = fmt.Sscanf(strings.TrimSpace(chunks[0]), "%s up %d days", &time, &daysUp)
	if err != nil {
		return Uptime{}, err
	}
	_, err = fmt.Sscanf(strings.TrimSpace(chunks[2]), "%d users", &users)
	if err != nil {
		return Uptime{}, err
	}
	_, err = fmt.Sscanf(strings.TrimSpace(chunks[3]), "load average: %s", &cpu1)
	if err != nil {
		return Uptime{}, err
	}
	_, err = fmt.Sscanf(strings.TrimSpace(chunks[4]), "%s", &cpu5)
	if err != nil {
		return Uptime{}, err
	}
	_, err = fmt.Sscanf(strings.TrimSpace(chunks[5]), "%s", &cpu15)
	if err != nil {
		return Uptime{}, err
	}

	return Uptime{
		Time:          time,
		DaysUp:        daysUp,
		UsersLoggedIn: users,
		CPUload1min:   cpu1,
		CPUload5min:   cpu5,
		CPUload15min:  cpu15,
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

func ParseProcess(s string) (_ Process, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("parsing `ps -u` command output: %w", err)
		}
	}()

	var (
		user, pid, cpu, mem string
	)

	line := strings.TrimSpace(s)
	chunks := strings.Split(line, "")
	if len(chunks) != 11 {
		return Process{}, err
	}

	return Process{}, nil

}
