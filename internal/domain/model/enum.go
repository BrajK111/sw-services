package model

// Consolidated enum set. Java carries two overlapping/inconsistent enum packages
// (web.models.* vs web.models.enums.*, with two conflicting Source enums) —
// the Go rewrite keeps a single authoritative set instead of replicating that split.
type Status string

const (
	StatusActive     Status = "ACTIVE"
	StatusInactive   Status = "INACTIVE"
	StatusInWorkflow Status = "INWORKFLOW"
)

// ApplicationStatus represents the sewerage connection workflow states.
// Transitions: INITIATED → PENDING_FOR_FIELD_INSPECTION → PENDING_FOR_APPROVAL → APPROVED → ACTIVE
type ApplicationStatus string

const (
	AppStatusInitiated                    ApplicationStatus = "INITIATED"
	AppStatusPendingForFieldInspection    ApplicationStatus = "PENDING_FOR_FIELD_INSPECTION"
	AppStatusPendingForApproval           ApplicationStatus = "PENDING_FOR_APPROVAL"
	AppStatusApproved                     ApplicationStatus = "APPROVED"
	AppStatusRejected                     ApplicationStatus = "REJECTED"
	AppStatusConnectionActivated          ApplicationStatus = "CONNECTION_ACTIVATED"
	AppStatusPendingApprovalForDisconnect ApplicationStatus = "PENDING_APPROVAL_FOR_DISCONNECTION"
)

type Channel string

const (
	ChannelSystem     Channel = "SYSTEM"
	ChannelCFCCounter Channel = "CFC_COUNTER"
	ChannelCitizen    Channel = "CITIZEN"
	ChannelDataEntry  Channel = "DATA_ENTRY"
	ChannelMigration  Channel = "MIGRATION"
)

type Relationship string

const (
	RelationshipFather  Relationship = "FATHER"
	RelationshipHusband Relationship = "HUSBAND"
)

