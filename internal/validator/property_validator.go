// Package validator mirrors Java's validator package: ValidateProperty and
// MDMSValidator, both invoked from the service layer before persisting a
// connection request. Issues doc §3 flagged these peer integrations as
// completely empty stubs in the prior Go attempt.
package validator

import (
	"github.com/egovernments/sw-services-go/internal/domain/dto"
	swerrors "github.com/egovernments/sw-services-go/internal/domain/errors"
	"github.com/egovernments/sw-services-go/internal/util"
)

type PropertyValidator struct {
	client *util.PropertyClient
}

func NewPropertyValidator(client *util.PropertyClient) *PropertyValidator {
	return &PropertyValidator{client: client}
}

// ValidatePropertyForConnection mirrors ValidateProperty.getOrValidateProperty
// + validatePropertyFields: the property must exist and be ACTIVE.
func (v *PropertyValidator) ValidatePropertyForConnection(requestInfo dto.RequestInfo, tenantID, propertyID string) error {
	prop, err := v.client.SearchByPropertyID(requestInfo, tenantID, propertyID)
	if err != nil {
		return swerrors.New("EG_SW_PROPERTY_NOT_FOUND", "could not validate propertyId "+propertyID+": "+err.Error())
	}
	if prop.Status != "ACTIVE" {
		return swerrors.New("EG_SW_PROPERTY_NOT_ACTIVE", "property "+propertyID+" is not ACTIVE (status: "+prop.Status+")")
	}
	return nil
}
