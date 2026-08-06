package dto

type SewerageConnectionResponse struct {
	ResponseInfo        ResponseInfo         `json:"ResponseInfo"`
	SewerageConnections []SewerageConnection `json:"SewerageConnections"`
	TotalCount          int64                `json:"totalCount,omitempty"`
}
