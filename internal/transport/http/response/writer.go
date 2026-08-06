package response

import (
	"net/http"

	"github.com/gin-gonic/gin"

	swerrors "github.com/egovernments/sw-services-go/internal/domain/errors"
)

// accessDeniedCodes lists domain error codes that should map to HTTP 403
// (Forbidden) rather than 400 (Bad Request), so clients and proxies can
// correctly distinguish authorization failures from validation errors.
var accessDeniedCodes = map[string]bool{
	"EG_SW_ACCESS_DENIED": true,
}

// WriteError maps a service-layer error to the DIGIT standard error envelope.
// EG_SW_ACCESS_DENIED → 403; other CustomExceptions → 400; unexpected → 500.
func WriteError(c *gin.Context, err error) {
	if custom, ok := err.(*swerrors.CustomException); ok {
		status := http.StatusBadRequest
		if accessDeniedCodes[custom.Code] {
			status = http.StatusForbidden
		}
		c.JSON(status, custom.Response())
		return
	}
	c.JSON(http.StatusInternalServerError, swerrors.New("EG_SW_UNKNOWN_ERROR", err.Error()).Response())
}
