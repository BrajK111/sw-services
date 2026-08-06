package model

// SewerageConnection maps to eg_sw_connection, the root aggregate. Child
// tables are loaded/persisted via GORM's Preload/Association mechanism
// instead of Java's single-big-LEFT-JOIN-then-dedupe-in-app-code pattern
// (SWQueryBuilder + SewerageRowMapper) — cleaner and avoids that bug class.
type SewerageConnection struct {
	ID                       string  `gorm:"column:id;primaryKey"`
	TenantID                 string  `gorm:"column:tenantid;index"`
	PropertyID               string  `gorm:"column:property_id;index"`
	ApplicationNo            string  `gorm:"column:applicationno;index"`
	ApplicationStatus        string  `gorm:"column:applicationstatus;index"`
	Status                   string  `gorm:"column:status"`
	ConnectionNo             string  `gorm:"column:connectionno;index"`
	OldConnectionNo          string  `gorm:"column:oldconnectionno;index"`
	RoadType                 string  `gorm:"column:roadtype"`
	RoadCuttingArea          float64 `gorm:"column:roadcuttingarea"`
	AdhocRebate              float64 `gorm:"column:adhocrebate"`
	AdhocPenalty             float64 `gorm:"column:adhocpenalty"`
	AdhocPenaltyReason       string  `gorm:"column:adhocpenaltyreason"`
	AdhocPenaltyComment      string  `gorm:"column:adhocpenaltycomment"`
	AdhocRebateReason        string  `gorm:"column:adhocrebatereason"`
	AdhocRebateComment       string  `gorm:"column:adhocrebatecomment"`
	ApplicationType          string  `gorm:"column:applicationtype"`
	DateEffectiveFrom        int64   `gorm:"column:dateeffectivefrom"`
	Locality                 string  `gorm:"column:locality"`
	OldApplication           bool    `gorm:"column:isoldapplication"`
	AdditionalDetails        JSONB   `gorm:"column:additionaldetails;type:jsonb"`
	Channel                  string  `gorm:"column:channel"`
	IsDisconnectionTemporary bool    `gorm:"column:isdisconnectiontemporary"`
	DisconnectionReason      string  `gorm:"column:disconnectionreason"`
	AuditDetails             `gorm:"embedded"`

	ServiceDetail     *ConnectionServiceDetail `gorm:"foreignKey:ConnectionID;references:ID"`
	Documents         []Document               `gorm:"foreignKey:SWID;references:ID"`
	PlumberInfo       []PlumberInfo            `gorm:"foreignKey:SWID;references:ID"`
	RoadCuttingInfo   []RoadCuttingInfo        `gorm:"foreignKey:SWID;references:ID"`
	ConnectionHolders []ConnectionHolder       `gorm:"foreignKey:ConnectionID;references:ID"`
}

func (SewerageConnection) TableName() string { return "eg_sw_connection" }
