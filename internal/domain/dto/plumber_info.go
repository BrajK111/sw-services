package dto

type PlumberInfo struct {
	ID                    string `json:"id,omitempty"`
	Name                  string `json:"name,omitempty"`
	LicenseNo             string `json:"licenseNo,omitempty"`
	MobileNumber          string `json:"mobileNumber,omitempty"`
	Gender                string `json:"gender,omitempty"`
	FatherOrHusbandName   string `json:"fatherOrHusbandName,omitempty"`
	CorrespondenceAddress string `json:"correspondenceAddress,omitempty"`
	Relationship          string `json:"relationship,omitempty"`
}
