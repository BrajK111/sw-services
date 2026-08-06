package model

// ConnectionServiceDetail maps to eg_sw_service, the 1:1 sewerage-specific
// extension of eg_sw_connection (mirrors Java's SewerageConnection extends
// Connection split across two tables).
type ConnectionServiceDetail struct {
	ConnectionID               string `gorm:"column:connection_id;primaryKey"`
	ConnectionExecutionDate    int64  `gorm:"column:connectionexecutiondate"`
	DisconnectionExecutionDate int64  `gorm:"column:disconnectionexecutiondate"`
	NoOfWaterClosets           int    `gorm:"column:noofwaterclosets"`
	NoOfToilets                int    `gorm:"column:nooftoilets"`
	ConnectionType             string `gorm:"column:connectiontype"`
	ProposedWaterClosets       int    `gorm:"column:proposedwaterclosets"`
	ProposedToilets            int    `gorm:"column:proposedtoilets"`
	AppCreatedDate             int64  `gorm:"column:appcreateddate"`
	DetailsProvidedBy          string `gorm:"column:detailsprovidedby"`
	EstimationFileStoreID      string `gorm:"column:estimationfilestoreid"`
	SanctionFileStoreID        string `gorm:"column:sanctionfilestoreid"`
	EstimationLetterDate       int64  `gorm:"column:estimationletterdate"`
	AuditDetails               `gorm:"embedded"`
}

func (ConnectionServiceDetail) TableName() string { return "eg_sw_service" }
