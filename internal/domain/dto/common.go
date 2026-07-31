package dto

// RequestInfo/ResponseInfo mirror the minimal DIGIT common-contract fields
// every service accepts/returns. Full parity (plainAccessRequest/ABAC fields)
// is not needed until the Phase 2+ encryption/validator work lands.
type RequestInfo struct {
	APIID       string    `json:"apiId,omitempty"`
	Ver         string    `json:"ver,omitempty"`
	Ts          int64     `json:"ts,omitempty"`
	Action      string    `json:"action,omitempty"`
	DID         string    `json:"did,omitempty"`
	Key         string    `json:"key,omitempty"`
	MsgID       string    `json:"msgId,omitempty"`
	RequesterID string    `json:"requesterId,omitempty"`
	UserInfo    *UserInfo `json:"userInfo,omitempty"`
	AuthToken   string    `json:"authToken,omitempty"`
}

// Role mirrors the DIGIT role object inside userInfo.roles[].
type Role struct {
	ID       *int64 `json:"id"`
	Code     string `json:"code"`
	Name     string `json:"name,omitempty"`
	TenantID string `json:"tenantId,omitempty"`
}

type UserInfo struct {
	ID           *int64 `json:"id"`
	UUID         string `json:"uuid,omitempty"`
	UserName     string `json:"userName,omitempty"`
	Name         string `json:"name"`
	Type         string `json:"type,omitempty"`
	MobileNumber string `json:"mobileNumber,omitempty"`
	EmailID      string `json:"emailId"`
	TenantID     string `json:"tenantId,omitempty"`
	Roles        []Role `json:"roles,omitempty"`
}


type ResponseInfo struct {
	APIID    string `json:"apiId,omitempty"`
	Ver      string `json:"ver,omitempty"`
	Ts       int64  `json:"ts,omitempty"`
	ResMsgID string `json:"resMsgId,omitempty"`
	MsgID    string `json:"msgId,omitempty"`
	Status   string `json:"status"`
}

func ResponseInfoFromRequestInfo(req *RequestInfo, statusCode string) ResponseInfo {
	ri := ResponseInfo{Status: statusCode}
	if req != nil {
		ri.APIID = req.APIID
		ri.Ver = req.Ver
		ri.Ts = req.Ts
		ri.MsgID = req.MsgID
	}
	return ri
}

type AuditDetails struct {
	CreatedBy        string `json:"createdBy,omitempty"`
	LastModifiedBy   string `json:"lastModifiedBy,omitempty"`
	CreatedTime      int64  `json:"createdTime,omitempty"`
	LastModifiedTime int64  `json:"lastModifiedTime,omitempty"`
}
