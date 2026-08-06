package dto

// OwnerInfo mirrors Java's OwnerInfo (extends web.models.users.User). Only a
// subset (UUID/OwnerType/OwnerShipPercentage/IsPrimaryOwner/Relationship/Status)
// is persisted in eg_sw_connectionholder in Phase 1 — the rest is carried
// through on the wire as-is until the Phase 2 user-service integration lands
// (Java resolves these live from egov-user, it doesn't own them either).
type OwnerInfo struct {
	ID                    string  `json:"id,omitempty"`
	UUID                  string  `json:"uuid,omitempty"`
	UserName              string  `json:"userName,omitempty"`
	Name                  string  `json:"name,omitempty"`
	Gender                string  `json:"gender,omitempty"`
	MobileNumber          string  `json:"mobileNumber,omitempty"`
	FatherOrHusbandName   string  `json:"fatherOrHusbandName,omitempty"`
	CorrespondenceAddress string  `json:"correspondenceAddress,omitempty"`
	IsPrimaryOwner        bool    `json:"isPrimaryOwner"`
	OwnerShipPercentage   float64 `json:"ownerShipPercentage,omitempty"`
	OwnerType             string  `json:"ownerType"`
	Relationship          string  `json:"relationship,omitempty"`
	Status                string  `json:"status,omitempty"`
}
