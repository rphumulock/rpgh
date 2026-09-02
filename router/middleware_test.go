package router

import (
	"strings"
	"testing"
)

// directives splits a policy into the sources named under each directive, so a
// test can ask about one without matching a host that happens to appear under
// another -- 'self' is in most of them, and cloudflareinsights.com would be a
// different thing entirely under connect-src.
func directives(policy string) map[string][]string {
	out := map[string][]string{}
	for _, d := range strings.Split(policy, "; ") {
		if name, sources, ok := strings.Cut(strings.TrimSpace(d), " "); ok {
			out[name] = strings.Fields(sources)
		} else {
			out[strings.TrimSpace(d)] = nil
		}
	}
	return out
}

func allows(policy, directive, source string) bool {
	for _, s := range directives(policy)[directive] {
		if s == source {
			return true
		}
	}
	return false
}

// TestBeaconIsAllowed guards a failure that reads as something else entirely.
// Cloudflare injects the analytics beacon at the edge, after this policy has
// been written, so it cannot carry our nonce and is here by host or not at
// all. Dropped from script-src it is refused by the browser, and the console
// error looks exactly like DNS null-routing the host -- so the first place
// anyone looks is the network, not this line.
func TestBeaconIsAllowed(t *testing.T) {
	p := policy("abc123")

	if !allows(p, "script-src", "https://static.cloudflareinsights.com") {
		t.Errorf("script-src does not allow the beacon: %q", directives(p)["script-src"])
	}

	// Under automatic injection the beacon reports to this origin's
	// /cdn-cgi/rum. Installing the snippet by hand instead posts to
	// cloudflareinsights.com, which is a connect-src change -- so if that host
	// ever turns up here, the comment above script-src has gone stale.
	if !allows(p, "connect-src", "'self'") {
		t.Errorf("connect-src does not allow the beacon's own origin: %q", directives(p)["connect-src"])
	}
}

// TestNonceReplacesUnsafeEval is the trade the nonce exists to make: Datastar
// compiles its expressions either into script elements carrying our nonce or
// through the Function constructor, and only the second needs eval allowed.
// Shipping both would keep the hole open for no gain.
func TestNonceReplacesUnsafeEval(t *testing.T) {
	with := policy("abc123")
	if !allows(with, "script-src", "'nonce-abc123'") {
		t.Errorf("the nonce never reaches script-src: %q", directives(with)["script-src"])
	}
	if allows(with, "script-src", "'unsafe-eval'") {
		t.Error("script-src allows eval even with a nonce, which is what the nonce is for")
	}

	// No nonce means Datastar is on its Function path, and a policy that
	// blocks it there breaks every signal on the page rather than hardening it.
	without := policy("")
	if !allows(without, "script-src", "'unsafe-eval'") {
		t.Errorf("without a nonce, script-src blocks Datastar outright: %q", directives(without)["script-src"])
	}
}

// TestNothingCanFrameOrPostToUs covers the directives with no legitimate use
// here at all, which are the ones most likely to be loosened by accident while
// chasing something else.
func TestNothingCanFrameOrPostToUs(t *testing.T) {
	d := directives(policy("abc123"))
	for _, want := range []struct{ directive, source string }{
		{"object-src", "'none'"},
		{"base-uri", "'none'"},
		{"form-action", "'none'"},
		{"frame-ancestors", "'none'"},
	} {
		if got := d[want.directive]; len(got) != 1 || got[0] != want.source {
			t.Errorf("%s = %q, want [%s]", want.directive, got, want.source)
		}
	}
}
