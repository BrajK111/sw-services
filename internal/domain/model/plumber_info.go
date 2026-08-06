package model

// PlumberInfo maps to eg_sw_plumberinfo.
type PlumberInfo struct {
	ID                    string `gorm:"column:id;primaryKey"`
	TenantID              string `gorm:"column:tenantid"`
	Name                  string `gorm:"column:name"`
	LicenseNo             string `gorm:"column:licenseno"`
	MobileNumber          string `gorm:"column:mobilenumber"`
	Gender                string `gorm:"column:gender"`
	FatherOrHusbandName   string `gorm:"column:fatherorhusbandname"`
	CorrespondenceAddress string `gorm:"column:correspondenceaddress"`
	Relationship          string `gorm:"column:relationship"`
	SWID                  string `gorm:"column:swid;index"`
	AuditDetails          `gorm:"embedded"`
}

func (PlumberInfo) TableName() string { return "eg_sw_plumberinfo" }
