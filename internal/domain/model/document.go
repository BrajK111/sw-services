package model

// Document maps to eg_sw_applicationdocument.
type Document struct {
	ID           string `gorm:"column:id;primaryKey"`
	TenantID     string `gorm:"column:tenantid"`
	DocumentType string `gorm:"column:documenttype"`
	FileStoreID  string `gorm:"column:filestoreid"`
	DocumentUID  string `gorm:"column:documentuid"`
	SWID         string `gorm:"column:swid;index"`
	Active       bool   `gorm:"column:active"`
	AuditDetails `gorm:"embedded"`
}

func (Document) TableName() string { return "eg_sw_applicationdocument" }
