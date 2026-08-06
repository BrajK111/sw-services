package validator

import (
	"strings"

	"github.com/egovernments/sw-services-go/internal/domain/dto"
	swerrors "github.com/egovernments/sw-services-go/internal/domain/errors"
	"github.com/egovernments/sw-services-go/internal/util"
)

const (
	moduleWSServicesMasters = "ws-services-masters"
	masterConnectionType    = "connectionType"
)

type mdmsCodeEntry struct {
	Code   string `json:"code"`
	Active bool   `json:"active"`
}

// MDMSValidator mirrors Java's MDMSValidator: validates that the
// connectionType on a request is a real, active MDMS-configured code —
// issues doc §7 flagged this validator as entirely missing.
type MDMSValidator struct {
	client *util.MdmsClient
}

func NewMDMSValidator(client *util.MdmsClient) *MDMSValidator {
	return &MDMSValidator{client: client}
}

func (v *MDMSValidator) ValidateConnectionType(requestInfo dto.RequestInfo, tenantID, connectionType string) error {
	if connectionType == "" {
		return nil
	}

	res, err := v.client.Search(requestInfo, tenantID, []util.ModuleDetail{
		{ModuleName: moduleWSServicesMasters, MasterDetails: []util.MasterDetail{{Name: masterConnectionType}}},
	})
	if err != nil {
		return swerrors.New("EG_SW_MDMS_UNREACHABLE", "could not validate connectionType against MDMS: "+err.Error())
	}

	var entries []mdmsCodeEntry
	if err := res.Master(moduleWSServicesMasters, masterConnectionType, &entries); err != nil {
		return swerrors.New("EG_SW_MDMS_PARSE_ERROR", "could not parse connectionType master: "+err.Error())
	}

	for _, e := range entries {
		if strings.EqualFold(e.Code, connectionType) && e.Active {
			return nil
		}
	}
	return swerrors.New("EG_SW_INVALID_CONNECTIONTYPE", "connectionType "+connectionType+" is not a valid active MDMS master value")
}
