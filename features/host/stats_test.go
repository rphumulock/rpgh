package host

import "testing"

// TestUptime pins the format shown in the footer: coarsest unit first, never
// more than two, so a box up for a week does not report its seconds.
func TestUptime(t *testing.T) {
	for _, tc := range []struct {
		seconds uint64
		want    string
	}{
		{45, "45s"},
		{90, "1m"},
		{3600, "1h 0m"},
		{3661, "1h 1m"},
		{86400, "1d 0h"},
		{270000, "3d 3h"},
	} {
		if got := uptime(tc.seconds); got != tc.want {
			t.Errorf("uptime(%d) = %q, want %q", tc.seconds, got, tc.want)
		}
	}
}

func TestBytes(t *testing.T) {
	for _, tc := range []struct {
		n    uint64
		want string
	}{
		{512, "512 B"},
		{2048, "2.0 KB"},
		{5 << 20, "5.0 MB"},
		{3 << 30, "3.0 GB"},
	} {
		if got := bytes(tc.n); got != tc.want {
			t.Errorf("bytes(%d) = %q, want %q", tc.n, got, tc.want)
		}
	}
}

// TestCollectNeverBlank guards the footer: a collector that fails should show
// the unknown marker, never an empty span that renders as a gap.
func TestCollectNeverBlank(t *testing.T) {
	s := Collect()
	for name, v := range map[string]string{
		"uptime": s.Uptime, "load": s.Load, "mem": s.Mem, "cpu": s.CPU,
	} {
		if v == "" {
			t.Errorf("%s is empty", name)
		}
	}
}

// TestCollectIsCached checks the shared pass: every open stream calls Collect
// on its own ticker, and they must not each walk /proc.
func TestCollectIsCached(t *testing.T) {
	first := Collect()
	if second := Collect(); second != first {
		t.Errorf("two calls inside the TTL returned different stats:\n%+v\n%+v", first, second)
	}
}
