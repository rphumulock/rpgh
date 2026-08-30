// Package csp carries the per-response Content-Security-Policy nonce from the
// middleware that mints it to the layout that has to echo it.
//
// Datastar needs the nonce because of how it compiles expressions. Without one
// it falls back to the Function constructor, which a policy has to allow with
// 'unsafe-eval'; with one it compiles each expression by appending a script
// element carrying the nonce, so the policy can drop 'unsafe-eval' entirely.
// The attribute is read once at startup and removed from the document.
package csp

import (
	"context"
	"crypto/rand"
	"encoding/base64"
)

type ctxKey struct{}

// NewNonce returns 128 bits of randomness, base64'd. A nonce is only worth
// anything while it is unpredictable, so this is minted per response and never
// reused -- which is also why a page carrying one must not be shared-cached.
func NewNonce() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		// crypto/rand does not fail on any platform this runs on; if it ever
		// did, serving a guessable nonce would be worse than not serving one.
		return ""
	}
	return base64.RawStdEncoding.EncodeToString(b)
}

// WithNonce stores the response's nonce for the layout to read.
func WithNonce(ctx context.Context, nonce string) context.Context {
	return context.WithValue(ctx, ctxKey{}, nonce)
}

// Nonce returns the nonce for this response, or "" if the request did not pass
// through the middleware. The layout omits the attribute when it is empty:
// Datastar throws on an empty data-nonce, and no attribute at least leaves it
// on its old code path rather than failing outright.
func Nonce(ctx context.Context) string {
	nonce, _ := ctx.Value(ctxKey{}).(string)
	return nonce
}
