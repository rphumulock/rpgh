package router

import (
	"net/http"
	"strings"

	"rpgh/config"
	"rpgh/web/csp"
)

// policy allowlists exactly what a page here loads: the iconify web component
// and its icon API, and Google Fonts' stylesheet plus the font files it points
// at. Everything else -- the site's own CSS, JS and images -- is same-origin.
//
// The nonce is what lets script-src stay free of 'unsafe-eval'. Datastar reads
// it from the html element and compiles every data-* expression into a script
// element carrying it; without it, expressions go through the Function
// constructor and the policy would have to allow eval outright.
//
// A host allowlist still applies alongside a nonce, so iconify keeps loading.
// That would stop being true under 'strict-dynamic', which is why it is absent.
func policy(nonce string) string {
	script := "script-src 'self' https://code.iconify.design"
	if nonce != "" {
		script += " 'nonce-" + nonce + "'"
	} else {
		// No nonce means Datastar is on its Function path, and blocking it
		// there breaks every signal on the page.
		script += " 'unsafe-eval'"
	}

	return strings.Join([]string{
		"default-src 'self'",
		script,
		// covers the style attributes Datastar toggles on data-show elements
		"style-src 'self' 'unsafe-inline' https://fonts.googleapis.com",
		"font-src 'self' https://fonts.gstatic.com",
		// Iconify's public API is served from three interchangeable hosts and
		// the component fails over between them, so allowlisting only the
		// first one makes icons load or not depending on which it picks.
		"connect-src 'self' https://api.iconify.design https://api.simplesvg.com https://api.unisvg.com",
		"img-src 'self' data:",
		"object-src 'none'",
		"base-uri 'none'",
		"form-action 'none'",
		"frame-ancestors 'none'",
	}, "; ")
}

// SecurityHeaders mints the response's CSP nonce and sets the headers a static,
// form-free site can send unconditionally. There is no user input anywhere on
// the site, so this is defence in depth rather than a fix for a known hole.
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nonce := csp.NewNonce()
		h := w.Header()
		h.Set("Content-Security-Policy", policy(nonce))
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		h.Set("X-Frame-Options", "DENY")
		h.Set("Permissions-Policy", "geolocation=(), microphone=(), camera=(), interest-cohort=()")

		// A fresh nonce per response only holds if the response is not shared
		// between visitors, so keep pages out of any shared cache. Static
		// assets carry no nonce and are content-hashed in prod, so they are
		// left alone to be cached hard.
		if !strings.HasPrefix(r.URL.Path, "/static/") {
			h.Set("Cache-Control", "no-store")
		}

		// HSTS is only meaningful, and only safe to send, once the response is
		// known to have travelled over TLS. Behind a proxy that is the
		// forwarded scheme, which is trustworthy only if the proxy is.
		https := r.TLS != nil ||
			(config.Global.TrustProxy && r.Header.Get("X-Forwarded-Proto") == "https")
		if https {
			h.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}

		next.ServeHTTP(w, r.WithContext(csp.WithNonce(r.Context(), nonce)))
	})
}
