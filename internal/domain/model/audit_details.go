package model

// AuditDetails is embedded (not a child table) — column names are identical
// across every table that carries it, per the Java schema inventory.
type AuditDetails struct {
	CreatedBy        string `gorm:"column:createdby"`
	LastModifiedBy   string `gorm:"column:lastmodifiedby"`
	CreatedTime      int64  `gorm:"column:createdtime"`
	LastModifiedTime int64  `gorm:"column:lastmodifiedtime"`
}
