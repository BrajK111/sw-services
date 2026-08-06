package model

// RoadCuttingInfo maps to eg_sw_roadcuttinginfo.
type RoadCuttingInfo struct {
	ID              string  `gorm:"column:id;primaryKey"`
	TenantID        string  `gorm:"column:tenantid"`
	SWID            string  `gorm:"column:swid;index"`
	Active          bool    `gorm:"column:active"`
	RoadType        string  `gorm:"column:roadtype"`
	RoadCuttingArea float64 `gorm:"column:roadcuttingarea"`
	AuditDetails    `gorm:"embedded"`
}

func (RoadCuttingInfo) TableName() string { return "eg_sw_roadcuttinginfo" }
