package mikrus_test

import (
	"slices"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/qba73/mikrus"
)

func TestRenderMemoryStats(t *testing.T) {
	t.Parallel()
	ts := newTestServer("/stats", []byte(statResponse), t)
	defer ts.Close()

	client := mikrus.New("dummyAPIKey", "dummySrv")
	client.HTTPClient = ts.Client()
	client.URL = ts.URL
}

func TestParseMemoryUsage_ParsesCommandOutputOnValidInput(t *testing.T) {
	t.Parallel()
	freeCmdOutput := "total        used        free      shared  buff/cache   available\nMem:           1024          43         816           0         164         980\nSwap:             0           0           0"
	got, err := mikrus.ParseMemoryUsage(freeCmdOutput)
	if err != nil {
		t.Fatal(err)
	}

	want := mikrus.Memory{
		Total:     1024,
		Used:      43,
		Free:      816,
		Shared:    0,
		Cache:     164,
		Available: 980,
		SwapTotal: 0,
		SwapUsed:  0,
		SwapFree:  0,
	}

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestParseDiskSpace_ParsesCommandOutputOnValidInput(t *testing.T) {
	t.Parallel()
	dfCmdOutput := "Filesystem                        Size  Used Avail Use% Mounted on\n/dev/mapper/pve-vm--230--disk--0  9.8G  2.7G  6.7G  29% /\nudev                               63G     0   63G   0% /dev/net"
	got, err := mikrus.ParseDiskSpace(dfCmdOutput)
	if err != nil {
		t.Fatal(err)
	}
	want := mikrus.DiskSpace{
		Filesystem: "/dev/mapper/pve-vm--230--disk--0",
		Size:       "9.8G",
		Used:       "2.7G",
		Available:  "6.7G",
		Usage:      "29%",
		MountedOn:  "/",
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestParseUptime_ParsesUptimeCommandOutput(t *testing.T) {
	t.Parallel()
	uptimeCmdOutput := "16:32:02 up 6 days,  8:33,  0 users,  load average: 0.10, 1.00, 0.50"
	wantUptime, err := time.ParseDuration("152h33m0s")
	if err != nil {
		t.Fatal(err)
	}
	want := mikrus.Uptime{
		Uptime:       wantUptime,
		Users:        0,
		CPUload1min:  0.1,
		CPUload5min:  1.0,
		CPUload15min: 0.5,
	}
	got, err := mikrus.ParseUptime(uptimeCmdOutput)
	if err != nil {
		t.Fatal(err)
	}
	if want != got {
		t.Error(cmp.Diff(want, got))
	}
}

func TestParseUptime_ErrorsForInvalidInput(t *testing.T) {
	t.Parallel()
	_, err := mikrus.ParseUptime("bogus")
	if err == nil {
		t.Fatal("want error for invalid input, got nil")
	}
}

func TestParsePS_ParsesPSCommandOutput(t *testing.T) {
	t.Parallel()
	psCmdOutput := "USER       PID %CPU %MEM    VSZ   RSS TTY      STAT START   TIME COMMAND\nroot     21605  0.0  0.3   9504  3368 ?        S    16:32   0:00 bash -c cat | sh\nroot     21607  0.0  0.0   2608   596 ?        S    16:32   0:00  \\_ sh\n"
	want := []mikrus.ProcessInfo{
		{
			User:              "root",
			PID:               21605,
			CPUPercent:        0.0,
			MemoryPercent:     0.3,
			VirtualMemorySize: 9504,
			ResidentSetSize:   3368,
			TTY:               "?",
			State:             "S",
			Start:             "16:32",
			CPUTime:           "0:00",
			Command:           "bash -c cat | sh",
		},
		{
			User:              "root",
			PID:               21607,
			CPUPercent:        0.0,
			MemoryPercent:     0.0,
			VirtualMemorySize: 2608,
			ResidentSetSize:   596,
			TTY:               "?",
			State:             "S",
			Start:             "16:32",
			CPUTime:           "0:00",
			Command:           "\\_ sh",
		},
	}
	got, err := mikrus.ParsePS(psCmdOutput)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestParsePS_ErrorsForInvalidInput(t *testing.T) {
	t.Parallel()
	_, err := mikrus.ParsePS("bogus")
	if err == nil {
		t.Fatal("want error for invalid input, got nil")
	}
}

var statResponse = `{
    "free": "total        used        free      shared  buff/cache   available\nMem:           1024          43         816           0         164         980\nSwap:             0           0           0",
    "df": "Filesystem                        Size  Used Avail Use% Mounted on\n/dev/mapper/pve-vm--230--disk--0  9.8G  2.7G  6.7G  29% /\nudev                               63G     0   63G   0% /dev/net",
    "uptime": "16:32:02 up 6 days,  8:33,  0 users,  load average: 0.00, 0.00, 0.00\nsh: 1: echo",
    "ps": ": not found\nUSER       PID %CPU %MEM    VSZ   RSS TTY      STAT START   TIME COMMAND\nroot     21605  0.0  0.3   9504  3368 ?        S    16:32   0:00 bash -c cat | sh\nroot     21607  0.0  0.0   2608   596 ?        S    16:32   0:00  \\_ sh\nroot     21612  0.0  0.3  11420  3264 ?        R    16:32   0:00      \\_ ps auxf\nroot         1  0.0  1.0 169412 10748 ?        Ss   Jun05   0:04 /sbin/init\nroot        48  0.0  6.0 141148 63436 ?        Ss   Jun05   0:20 /lib/systemd/systemd-journald\nsystemd+    73  0.0  0.7  18376  7616 ?        Ss   Jun05   0:00 /lib/systemd/systemd-networkd\nsystemd+    89  0.0  1.1  23924 12128 ?        Ss   Jun05   0:14 /lib/systemd/systemd-resolved\nroot       103  0.0  0.6 238088  7264 ?        Ssl  Jun05   0:06 /usr/lib/accountsservice/accounts-daemon\nroot       104  0.0  0.2   9344  2796 ?        Ss   Jun05   0:00 /usr/sbin/cron -f\nmessage+   105  0.0  0.3   7404  4100 ?        Ss   Jun05   0:00 /usr/bin/dbus-daemon --system --address=systemd: --nofork --nopidfile --systemd-activation --syslog-only\nroot       108  0.0  1.7  31804 18404 ?        Ss   Jun05   0:00 /usr/bin/python3 /usr/bin/networkd-dispatcher --run-startup-triggers\nsyslog     110  0.0  0.4 154708  4248 ?        Ssl  Jun05   0:01 /usr/sbin/rsyslogd -n -iNONE\nroot       113  0.0  0.5  16440  6124 ?        Ss   Jun05   0:00 /lib/systemd/systemd-logind\nroot       122  0.0  0.2   8132  2144 console  Ss+  Jun05   0:00 /sbin/agetty -o -p -- \\u --noclear --keep-baud console 115200,38400,9600 linux\nroot       123  0.0  0.2   8132  2248 pts/0    Ss+  Jun05   0:00 /sbin/agetty -o -p -- \\u --noclear --keep-baud tty1 115200,38400,9600 linux\nroot       124  0.0  0.2   8132  2132 pts/1    Ss+  Jun05   0:00 /sbin/agetty -o -p -- \\u --noclear --keep-baud tty2 115200,38400,9600 linux\nroot       126  0.0  0.6  12172  7124 ?        Ss   Jun05   0:03 sshd: /usr/sbin/sshd -D [listener] 0 of 10-100 startups"
}`
