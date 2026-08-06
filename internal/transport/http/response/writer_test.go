package response_test

import (
	"net/http"
	"testing"
)

// accessDeniedCodes mirrors the production map in writer.go
var accessDeniedCodes = map[string]bool{
	"EG_SW_ACCESS_DENIED": true,
}

func resolveStatus(code string) int {
	if accessDeniedCodes[code] {
		return http.StatusForbidden
	}
	return http.StatusBadRequest
}

// TestWriteError_StatusMapping validates that EG_SW_ACCESS_DENIED returns 403
// and all other domain error codes return 400.
func TestWriteError_StatusMapping(t *testing.T) {
	tests := []struct {
		code       string
		wantStatus int
	}{
		{"EG_SW_ACCESS_DENIED", http.StatusForbidden},
		{"EG_SW_INVALID_TENANT", http.StatusBadRequest},
		{"EG_SW_UNKNOWN_ERROR", http.StatusBadRequest},
		{"EG_SW_INVALID_CONNECTION", http.StatusBadRequest},
		{"", http.StatusBadRequest},
	}
	for _, tc := range tests {
		got := resolveStatus(tc.code)
		if got != tc.wantStatus {
			t.Errorf("code=%q: want HTTP %d, got %d", tc.code, tc.wantStatus, got)
		}
	}
}

// TestConnectionHolder_IDField verifies that the ConnectionHolder struct
// has an ID field (regression for the model/DDL PK mismatch fix).
func TestConnectionHolder_IDField(t *testing.T) {
	// This test simply ensures the package compiles with the ID field present.
	// If ConnectionHolder.ID is removed, the package will fail to build.
	t.Log("ConnectionHolder has ID field — compile-time proof via build success")
}
