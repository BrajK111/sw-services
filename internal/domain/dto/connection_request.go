package dto

type SewerageConnectionRequest struct {
	RequestInfo        RequestInfo        `json:"RequestInfo"`
	SewerageConnection SewerageConnection `json:"sewerageConnection"`
	IsCreateCall       bool               `json:"isCreateCall,omitempty"`
	DisconnectRequest  bool               `json:"disconnectRequest,omitempty"`
	// CallerRoles is populated by the auth middleware from the Gin context,
	// NOT from the JSON body — clients cannot self-escalate privileges.
	CallerRoles []string `json:"-"`
}

type RequestInfoWrapper struct {
	RequestInfo RequestInfo `json:"RequestInfo"`
}

