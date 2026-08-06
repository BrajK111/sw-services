package util

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/egovernments/sw-services-go/internal/domain/dto"
)

type ModuleDetail struct {
	ModuleName    string         `json:"moduleName"`
	MasterDetails []MasterDetail `json:"masterDetails"`
}

type MasterDetail struct {
	Name string `json:"name"`
}

type mdmsCriteria struct {
	TenantID      string         `json:"tenantId"`
	ModuleDetails []ModuleDetail `json:"moduleDetails"`
}

type mdmsRequest struct {
	RequestInfo  dto.RequestInfo `json:"RequestInfo"`
	MdmsCriteria mdmsCriteria    `json:"MdmsCriteria"`
}

// MdmsResponse holds the raw per-module/per-master JSON arrays so callers can
// unmarshal only the masters they need — same shape as sw-calculator-go's client.
type MdmsResponse struct {
	ResponseInfo dto.ResponseInfo                      `json:"ResponseInfo"`
	MdmsRes      map[string]map[string]json.RawMessage `json:"MdmsRes"`
}

func (r MdmsResponse) Master(module, master string, target interface{}) error {
	moduleData, ok := r.MdmsRes[module]
	if !ok {
		return nil
	}
	raw, ok := moduleData[master]
	if !ok {
		return nil
	}
	return json.Unmarshal(raw, target)
}

type MdmsClient struct {
	host       string
	searchPath string
}

func NewMdmsClient(host, searchPath string) *MdmsClient {
	return &MdmsClient{host: host, searchPath: searchPath}
}

func (c *MdmsClient) Search(requestInfo dto.RequestInfo, tenantID string, modules []ModuleDetail) (*MdmsResponse, error) {
	reqBody := mdmsRequest{
		RequestInfo: requestInfo,
		MdmsCriteria: mdmsCriteria{
			TenantID:      tenantID,
			ModuleDetails: modules,
		},
	}

	var res MdmsResponse
	if err := PostJSON(JoinURL(c.host, c.searchPath), url.Values{}, reqBody, &res); err != nil {
		return nil, fmt.Errorf("mdms search: %w", err)
	}
	return &res, nil
}
