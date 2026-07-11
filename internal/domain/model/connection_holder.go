package model

// ConnectionHolder maps to eg_sw_connectionholder — the reference row linking a
// connection to a user-service UUID plus holder-specific attributes. Full owner
// profile fields (name, mobile number, etc.) are not stored here; Java's real
// design fetches those live from the user service. Two Java schema bugs are
// fixed here rather than replicated: (1) connectionid had a UNIQUE constraint,
// silently limiting every connection to one holder despite the model
// supporting a list — dropped; (2) HoldershipPercentage lives here as numeric,
// not the Java table's varchar, matching the actual OwnerShipPercentage(Double) model field.
type ConnectionHolder struct {
	TenantID             string  `gorm:"column:tenantid"`
	ConnectionID         string  `gorm:"column:connectionid;primaryKey"`
	Status               string  `gorm:"column:status"`
	UserID               string  `gorm:"column:userid"`
	IsPrimaryHolder      bool    `gorm:"column:isprimaryholder"`
	ConnectionHolderType string  `gorm:"column:connectionholdertype"`
	HoldershipPercentage float64 `gorm:"column:holdershippercentage"`
	Relationship         string  `gorm:"column:relationship"`
	AuditDetails         `gorm:"embedded"`
}

func (ConnectionHolder) TableName() string { return "eg_sw_connectionholder" }
