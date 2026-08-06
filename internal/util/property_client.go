package util

import (
	"fmt"
	"net/url"

	"github.com/egovernments/sw-services-go/internal/domain/dto"
)

type requestInfoWrapper struct {
	RequestInfo dto.RequestInfo `json:"RequestInfo"`
}

type propertyResponse struct {
	ResponseInfo dto.ResponseInfo `json:"ResponseInfo"`
	Properties   []Property       `json:"Properties"`
}

// Property is the subset of property-services' contract this validator needs.
type Property struct {
	PropertyID string `json:"propertyId"`
	TenantID   string `json:"tenantId"`
	Status     string `json:"status"`
}

type PropertyClient struct {
	host       string
	searchPath string
}

func NewPropertyClient(host, searchPath string) *PropertyClient {
	return &PropertyClient{host: host, searchPath: searchPath}
}

func (c *PropertyClient) SearchByPropertyID(requestInfo dto.RequestInfo, tenantID, propertyID string) (*Property, error) {
	q := url.Values{"tenantId": {tenantID}, "propertyIds": {propertyID}}
	var res propertyResponse
	body := requestInfoWrapper{RequestInfo: requestInfo}
	if err := PostJSON(JoinURL(c.host, c.searchPath), q, body, &res); err != nil {
		return nil, fmt.Errorf("property search: %w", err)
	}
	if len(res.Properties) == 0 {
		return nil, fmt.Errorf("no property found for propertyId %s", propertyID)
	}
	return &res.Properties[0], nil
}
