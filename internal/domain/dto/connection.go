package dto

// ProcessInstance represents the workflow action from the frontend (e.g. INITIATE, SUBMIT_APPLICATION)
type ProcessInstance struct {
	Action string `json:"action,omitempty"`
}

// SewerageConnection is the flat wire-format contract DIGIT-UI/Postman expect
// — it merges Java's Connection + SewerageConnection inheritance chain into
// one struct since Go has no class inheritance, and flattens the eg_sw_service
// child-table fields (ConnectionType, execution dates, toilet counts) back to
// top level to match the Java JSON shape exactly.
type SewerageConnection struct {
	ID                         string                 `json:"id,omitempty"`
	TenantID                   string                 `json:"tenantId"`
	PropertyID                 string                 `json:"propertyId"`
	ApplicationNo              string                 `json:"applicationNo,omitempty"`
	ApplicationStatus          string                 `json:"applicationStatus,omitempty"`
	Status                     string                 `json:"status,omitempty"`
	ConnectionNo               string                 `json:"connectionNo,omitempty"`
	OldConnectionNo            string                 `json:"oldConnectionNo,omitempty"`
	Documents                  []Document             `json:"documents,omitempty"`
	PlumberInfo                []PlumberInfo          `json:"plumberInfo,omitempty"`
	RoadType                   string                 `json:"roadType,omitempty"`
	RoadCuttingArea            float64                `json:"roadCuttingArea,omitempty"`
	RoadCuttingInfo            []RoadCuttingInfo      `json:"roadCuttingInfo,omitempty"`
	ConnectionExecutionDate    int64                  `json:"connectionExecutionDate,omitempty"`
	DisconnectionExecutionDate int64                  `json:"disconnectionExecutionDate,omitempty"`
	ConnectionType             string                 `json:"connectionType,omitempty"`
	NoOfWaterClosets           int                    `json:"noOfWaterClosets,omitempty"`
	NoOfToilets                int                    `json:"noOfToilets,omitempty"`
	ProposedWaterClosets       int                    `json:"proposedWaterClosets,omitempty"`
	ProposedToilets            int                    `json:"proposedToilets,omitempty"`
	AdditionalDetails          map[string]interface{} `json:"additionalDetails,omitempty"`
	AuditDetails               *AuditDetails          `json:"auditDetails,omitempty"`
	ConnectionHolders          []OwnerInfo            `json:"connectionHolders,omitempty"`
	ApplicationType            string                 `json:"applicationType,omitempty"`
	DateEffectiveFrom          int64                  `json:"dateEffectiveFrom,omitempty"`
	Locality                   string                 `json:"locality,omitempty"`
	OldApplication             bool                   `json:"oldApplication"`
	Channel                    string                 `json:"channel,omitempty"`
	IsDisconnectionTemporary   bool                   `json:"isDisconnectionTemporary"`
	DisconnectionReason        string                 `json:"disconnectionReason,omitempty"`
	ProcessInstance            *ProcessInstance       `json:"processInstance,omitempty"`
}
