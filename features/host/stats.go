// Package host reports what the machine serving this site is doing. It is a
// flourish rather than telemetry: the point is to show the site really is
// running on a desktop in a house, not to monitor anything.
//
// Everything here is read from inside the container. Proxmox mounts lxcfs, so
// /proc is namespaced -- uptime, load and memory describe the LXC, not the
// OptiPlex hosting it. That is why Chassis is a constant: the hardware is a
// fact about the box, and the live numbers are honestly labelled as the
// container's own.
package host

import (
	"fmt"
	"sync"
	"time"

	"rpgh/config"

	"github.com/shirou/gopsutil/v4/cpu"
	gohost "github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/load"
	"github.com/shirou/gopsutil/v4/mem"
)

// Chassis is the machine under the stairs, as named by the deployment. It is
// stated rather than detected because dmidecode needs hardware access an
// unprivileged container does not have -- so it is only true where it is set,
// and the footer leaves the segment out entirely when it is not.
func Chassis() string {
	return config.Global.HostChassis
}

// unknown is what a field shows when its collector failed. A stat that cannot
// be read should look absent, not look like zero.
const unknown = "--"

// cacheTTL is how long one collection is served to every open stream. Readers
// share it, so a hundred viewers still cost one pass over /proc per second
// rather than a hundred.
const cacheTTL = time.Second

// Stats is the payload, already formatted for display. The json names are the
// Datastar signals the footer binds to.
type Stats struct {
	Uptime string `json:"hostUptime"`
	Load   string `json:"hostLoad"`
	Mem    string `json:"hostMem"`
	CPU    string `json:"hostCpu"`
}

// Pending is what the page renders before the first event arrives.
func Pending() Stats {
	return Stats{Uptime: "...", Load: "...", Mem: "...", CPU: "..."}
}

var (
	mu     sync.Mutex
	cached Stats
	taken  time.Time
)

// Collect returns the current stats, reusing the last pass if it is still
// fresh. Safe to call from every open stream at once.
func Collect() Stats {
	mu.Lock()
	defer mu.Unlock()

	if !taken.IsZero() && time.Since(taken) < cacheTTL {
		return cached
	}

	cached = collect()
	taken = time.Now()
	return cached
}

func collect() Stats {
	s := Stats{Uptime: unknown, Load: unknown, Mem: unknown, CPU: unknown}

	if up, err := gohost.Uptime(); err == nil {
		s.Uptime = uptime(up)
	}
	if avg, err := load.Avg(); err == nil {
		s.Load = fmt.Sprintf("%.2f", avg.Load1)
	}
	if vm, err := mem.VirtualMemory(); err == nil {
		s.Mem = fmt.Sprintf("%s / %s", bytes(vm.Used), bytes(vm.Total))
	}
	// Interval 0 measures against the previous call rather than sampling in
	// line, which is what keeps this off the request path. It is why init
	// primes it: the very first call has nothing to compare against and
	// reports the average since boot instead.
	if pct, err := cpu.Percent(0, false); err == nil && len(pct) > 0 {
		s.CPU = fmt.Sprintf("%.1f%%", pct[0])
	}

	return s
}

func init() {
	// Discarded on purpose: this is the throwaway baseline for the deltas
	// every later cpu.Percent call reports against.
	_, _ = cpu.Percent(0, false)
}

// uptime renders seconds the way `uptime` does, coarsest unit first and never
// more than two, since nobody reads the seconds on a box that has been up for
// a week.
func uptime(seconds uint64) string {
	d := time.Duration(seconds) * time.Second
	switch {
	case d >= 24*time.Hour:
		return fmt.Sprintf("%dd %dh", int(d.Hours())/24, int(d.Hours())%24)
	case d >= time.Hour:
		return fmt.Sprintf("%dh %dm", int(d.Hours()), int(d.Minutes())%60)
	case d >= time.Minute:
		return fmt.Sprintf("%dm", int(d.Minutes()))
	default:
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
}

// bytes renders a byte count in the largest unit that leaves it above one.
func bytes(n uint64) string {
	const unit = 1024
	if n < unit {
		return fmt.Sprintf("%d B", n)
	}
	div, exp := uint64(unit), 0
	for n/div >= unit && exp < 3 {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(n)/float64(div), "KMGT"[exp])
}

var (
	hardwareOnce sync.Once
	hardware     string
)

// Hardware is the CPU line beside Chassis. lxcfs virtualises the count to what
// the container was given but leaves the model string alone, so this names the
// real silicon and the share of it this container can use. The count is of
// logical CPUs -- threads, not cores, which is what cpu.Counts(true) returns.
func Hardware() string {
	hardwareOnce.Do(func() {
		model := unknown
		if info, err := cpu.Info(); err == nil && len(info) > 0 && info[0].ModelName != "" {
			model = info[0].ModelName
		}
		threads, err := cpu.Counts(true)
		if err != nil || threads <= 0 {
			hardware = model
			return
		}
		hardware = fmt.Sprintf("%s · %d threads", model, threads)
	})
	return hardware
}
