package dto

type Document struct {
	ID           string `json:"id,omitempty"`
	TenantID     string `json:"tenantId,omitempty"`
	DocumentType string `json:"documentType"`
	FileStoreID  string `json:"fileStoreId"`
	DocumentUID  string `json:"documentUid,omitempty"`
	Status       string `json:"status,omitempty"`
}
