package util

import "strings"

// IsMasked mirrors Java's UnmaskingUtil convention: a string value containing
// the literal character '*' is treated as "masked/not actually provided by
// the client" rather than a real value.
func IsMasked(value string) bool {
	return strings.Contains(value, "*")
}

// ResolveIfMasked returns trusted when incoming looks masked, otherwise
// returns incoming unchanged — the "unmask" operation Java's UnmaskingUtil
// performs by substituting a previously-known-good value in place before
// persisting/validating.
func ResolveIfMasked(incoming, trusted string) string {
	if IsMasked(incoming) {
		return trusted
	}
	return incoming
}
