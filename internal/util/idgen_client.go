package util

import (
	"fmt"

	"github.com/egovernments/sw-services-go/internal/domain/dto"
)

type idGenerationRequest struct {
	RequestInfo dto.RequestInfo `json:"RequestInfo"`
	IDRequests  []idRequest     `json:"idRequests"`
}

type idRequest struct {
	IDName   string `json:"idName"`
	TenantID string `json:"tenantId"`
	Format   string `json:"format,omitempty"`
}

type idGenerationResponse struct {
	IDResponses []struct {
		ID string `json:"id"`
	} `json:"idResponses"`
}

type IdGenClient struct {
	host         string
	generatePath string
}

func NewIdGenClient(host, generatePath string) *IdGenClient {
	return &IdGenClient{host: host, generatePath: generatePath}
}

// GenerateApplicationNumber calls egov-idgen for the "sewerageservice.application.id"
// format. Callers should fall back to LocalApplicationNumber() on error rather
// than fail the citizen's request — mirrors the engineering-report's documented
// resilience pattern for an unresponsive IDGen service.
func (c *IdGenClient) GenerateApplicationNumber(requestInfo dto.RequestInfo, tenantID string) (string, error) {
	req := idGenerationRequest{
		RequestInfo: requestInfo,
		IDRequests: []idRequest{
			{IDName: "sewerageservice.application.id", TenantID: tenantID},
		},
	}

	var res idGenerationResponse
	if err := PostJSON(JoinURL(c.host, c.generatePath), nil, req, &res); err != nil {
		return "", fmt.Errorf("idgen generate: %w", err)
	}
	if len(res.IDResponses) == 0 || res.IDResponses[0].ID == "" {
		return "", fmt.Errorf("idgen returned no id")
	}
	return res.IDResponses[0].ID, nil
}
