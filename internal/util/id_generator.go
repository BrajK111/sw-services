package util

import (
	"crypto/rand"
	"fmt"
)

// LocalApplicationNumber is the fallback generator used until Phase 2 wires a
// real egov-idgen client. It mirrors the exact fallback format the engineering
// report specifies IDGen client uses when the peer service is unresponsive
// (SW-APP-xxxxxxxx), so switching to a real client later is a drop-in swap.
func LocalApplicationNumber() string {
	return fmt.Sprintf("SW-APP-%s", randomHex(8))
}

func randomHex(n int) string {
	bytes := make([]byte, n/2)
	if _, err := rand.Read(bytes); err != nil {
		return "00000000"
	}
	return fmt.Sprintf("%x", bytes)
}
