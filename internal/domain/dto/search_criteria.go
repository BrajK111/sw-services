package dto

// SearchCriteria mirrors Java's SearchCriteria. The audit flagged the prior Go
// attempt for only supporting tenantId+applicationNumber; this carries the
// full Java filter surface. mobileNumber/ownerName require a live user-service
// join which isn't wired in until Phase 2 — the field is accepted but not yet
// applied by the Phase 1 query builder (internal/repository/postgres).
type SearchCriteria struct {
	TenantID            string   `form:"tenantId" json:"tenantId,omitempty"`
	IDs                 []string `form:"ids" json:"ids,omitempty"`
	ApplicationNumber   []string `form:"applicationNumber" json:"applicationNumber,omitempty"`
	ConnectionNumber    []string `form:"connectionNumber" json:"connectionNumber,omitempty"`
	OldConnectionNumber string   `form:"oldConnectionNumber" json:"oldConnectionNumber,omitempty"`
	Status              string   `form:"status" json:"status,omitempty"`
	ApplicationStatus   []string `form:"applicationStatus" json:"applicationStatus,omitempty"`
	PropertyID          string   `form:"propertyId" json:"propertyId,omitempty"`
	ApplicationType     string   `form:"applicationType" json:"applicationType,omitempty"`
	Locality            string   `form:"locality" json:"locality,omitempty"`
	MobileNumber        string   `form:"mobileNumber" json:"mobileNumber,omitempty"`
	OwnerName           string   `form:"ownerName" json:"ownerName,omitempty"`
	DoorNo              string   `form:"doorNo" json:"doorNo,omitempty"`
	FromDate            int64    `form:"fromDate" json:"fromDate,omitempty"`
	ToDate              int64    `form:"toDate" json:"toDate,omitempty"`
	Offset              int      `form:"offset" json:"offset,omitempty"`
	Limit               int      `form:"limit" json:"limit,omitempty"`
	IsCountCall         bool     `json:"-"`
}
