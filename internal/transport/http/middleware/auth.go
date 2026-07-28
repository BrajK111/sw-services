// Package middleware provides Gin middleware for role-based authorization.
// It mirrors the DIGIT access-control pattern: roles are carried inside
// RequestInfo.userInfo.roles[] (full DIGIT format) OR userInfo.type (simple
// format used in local testing). Both are supported so Postman tests work
// with either style of request body.
package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/egovernments/sw-services-go/internal/domain/dto"
)

// SW role codes — sourced from pb-state-ACCESSCONTROL-ROLES.json in the MDMS data.
const (
	RoleCitizen           = "CITIZEN"
	RoleEmployee          = "EMPLOYEE"
	RoleSWCEMP            = "SW_CEMP"           // Counter Employee: handles walk-in creation
	RoleSWFieldInspector  = "SW_FIELD_INSPECTOR" // Field Inspector: does site visit & advances workflow
	RoleSWApprover        = "SW_APPROVER"        // Approver: final approval authority
	RoleSuperUser         = "SUPERUSER"          // Admin: unrestricted access
	RoleSystem            = "SYSTEM"             // Internal: service-to-service calls
)

// contextKey is the Gin context key where extracted roles are stored.
const contextKey = "sw_user_roles"

// extractRoles reads RequestInfo from the JSON body without consuming it
// (body is restored so the actual handler can still bind it). Returns the
// deduplicated set of role codes for the caller.
func extractRoles(c *gin.Context) []string {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return nil
	}
	// Restore body for downstream handler binding.
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	var wrapper struct {
		RequestInfo dto.RequestInfo `json:"RequestInfo"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil
	}

	seen := make(map[string]bool)
	var roles []string

	ui := wrapper.RequestInfo.UserInfo
	if ui == nil {
		return nil
	}

	// Support simple type field ("CITIZEN", "EMPLOYEE", "SYSTEM").
	if ui.Type != "" {
		code := ui.Type
		if !seen[code] {
			seen[code] = true
			roles = append(roles, code)
		}
	}
	// Support full DIGIT roles array.
	for _, r := range ui.Roles {
		if r.Code != "" && !seen[r.Code] {
			seen[r.Code] = true
			roles = append(roles, r.Code)
		}
	}
	return roles
}

// hasAnyRole returns true if the caller holds at least one of the allowed roles.
func hasAnyRole(callerRoles []string, allowed []string) bool {
	allowedSet := make(map[string]bool, len(allowed))
	for _, r := range allowed {
		allowedSet[r] = true
	}
	for _, r := range callerRoles {
		if allowedSet[r] {
			return true
		}
	}
	return false
}

// RequireRole returns a Gin middleware that allows the request only if the
// caller holds at least one of the specified role codes. On denial it writes
// a DIGIT-format 403 response and aborts the chain.
func RequireRole(allowed ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roles := extractRoles(c)
		c.Set(contextKey, roles)

		if !hasAnyRole(roles, allowed) {
			c.JSON(http.StatusForbidden, gin.H{
				"Errors": []gin.H{{
					"code":    "EG_SW_ACCESS_DENIED",
					"message": "Access denied: insufficient role for this operation",
				}},
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// GetRoles retrieves the caller's roles from the Gin context (set by RequireRole).
func GetRoles(c *gin.Context) []string {
	val, exists := c.Get(contextKey)
	if !exists {
		return nil
	}
	roles, _ := val.([]string)
	return roles
}
