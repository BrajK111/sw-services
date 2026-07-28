// Package handler maps 1:1 onto Java's SewarageController (@RequestMapping("/swc")).
package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/egovernments/sw-services-go/internal/domain/dto"
	swerrors "github.com/egovernments/sw-services-go/internal/domain/errors"
	"github.com/egovernments/sw-services-go/internal/service"
	"github.com/egovernments/sw-services-go/internal/transport/http/response"
)

type ConnectionHandler struct {
	service *service.ConnectionService
}

func NewConnectionHandler(svc *service.ConnectionService) *ConnectionHandler {
	return &ConnectionHandler{service: svc}
}

func (h *ConnectionHandler) responseInfo(req dto.RequestInfo) dto.ResponseInfo {
	return dto.ResponseInfoFromRequestInfo(&req, "successful")
}

func (h *ConnectionHandler) Create(c *gin.Context) {
	var req dto.SewerageConnectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteError(c, swerrors.New("EG_SW_INVALID_REQUEST", err.Error()))
		return
	}
	req.IsCreateCall = true

	created, err := h.service.CreateConnection(req)
	if err != nil {
		response.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.SewerageConnectionResponse{
		ResponseInfo:        h.responseInfo(req.RequestInfo),
		SewerageConnections: []dto.SewerageConnection{created},
	})
}

func (h *ConnectionHandler) Update(c *gin.Context) {
	var req dto.SewerageConnectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteError(c, swerrors.New("EG_SW_INVALID_REQUEST", err.Error()))
		return
	}
	// Inject roles from middleware context — never from the JSON body.
	if roles, ok := c.Get("sw_user_roles"); ok {
		req.CallerRoles, _ = roles.([]string)
	}

	updated, err := h.service.UpdateConnection(req)
	if err != nil {
		response.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.SewerageConnectionResponse{
		ResponseInfo:        h.responseInfo(req.RequestInfo),
		SewerageConnections: []dto.SewerageConnection{updated},
	})
}


func (h *ConnectionHandler) bindSearch(c *gin.Context) (dto.RequestInfoWrapper, dto.SearchCriteria, error) {
	var wrapper dto.RequestInfoWrapper
	if err := c.ShouldBindJSON(&wrapper); err != nil {
		return wrapper, dto.SearchCriteria{}, err
	}
	var criteria dto.SearchCriteria
	if err := c.ShouldBindQuery(&criteria); err != nil {
		return wrapper, dto.SearchCriteria{}, err
	}
	return wrapper, criteria, nil
}

func (h *ConnectionHandler) Search(c *gin.Context) {
	wrapper, criteria, err := h.bindSearch(c)
	if err != nil {
		response.WriteError(c, swerrors.New("EG_SW_INVALID_REQUEST", err.Error()))
		return
	}

	results, total, err := h.service.Search(criteria)
	if err != nil {
		response.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.SewerageConnectionResponse{
		ResponseInfo:        h.responseInfo(wrapper.RequestInfo),
		SewerageConnections: results,
		TotalCount:          total,
	})
}

func (h *ConnectionHandler) PlainSearch(c *gin.Context) {
	wrapper, criteria, err := h.bindSearch(c)
	if err != nil {
		response.WriteError(c, swerrors.New("EG_SW_INVALID_REQUEST", err.Error()))
		return
	}

	results, _, err := h.service.PlainSearch(criteria)
	if err != nil {
		response.WriteError(c, err)
		return
	}

	// Java sets totalCount = list size for _plainsearch (not a DB count query) — replicated as-is.
	c.JSON(http.StatusOK, dto.SewerageConnectionResponse{
		ResponseInfo:        h.responseInfo(wrapper.RequestInfo),
		SewerageConnections: results,
		TotalCount:          int64(len(results)),
	})
}

// EncryptOldData mirrors Java's own dead-code endpoint: the real batch logic
// is disabled there too (SewerageEncryptionService is never invoked), so this
// intentionally always returns the same disabled error rather than a stub
// implementation of something Java itself doesn't run.
func (h *ConnectionHandler) EncryptOldData(c *gin.Context) {
	response.WriteError(c, swerrors.New("EG_SW_ENC_OLD_DATA_ERROR", "Privacy disabled: old data encryption is not enabled for this environment"))
}

func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "UP", "checkedAt": time.Now().UTC()})
}
