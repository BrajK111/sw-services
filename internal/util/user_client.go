package util

import (
	"fmt"

	"github.com/egovernments/sw-services-go/internal/domain/dto"
)

type User struct {
	UUID         string `json:"uuid,omitempty"`
	UserName     string `json:"userName,omitempty"`
	Name         string `json:"name,omitempty"`
	MobileNumber string `json:"mobileNumber,omitempty"`
	Type         string `json:"type,omitempty"`
	TenantID     string `json:"tenantId,omitempty"`
	Roles        []Role `json:"roles,omitempty"`
	Active       bool   `json:"active"`
}

type Role struct {
	Code     string `json:"code"`
	TenantID string `json:"tenantId,omitempty"`
}

type userSearchRequest struct {
	RequestInfo  dto.RequestInfo `json:"RequestInfo"`
	TenantID     string          `json:"tenantId"`
	MobileNumber string          `json:"mobileNumber,omitempty"`
	UserType     string          `json:"userType,omitempty"`
}

type userSearchResponse struct {
	ResponseInfo dto.ResponseInfo `json:"ResponseInfo"`
	User         []User           `json:"user"`
}

type createUserRequest struct {
	RequestInfo dto.RequestInfo `json:"RequestInfo"`
	User        User            `json:"user"`
}

type createUserResponse struct {
	ResponseInfo dto.ResponseInfo `json:"ResponseInfo"`
	User         User             `json:"user"`
}

type UserClient struct {
	host             string
	searchPath       string
	createNoValidate string
}

func NewUserClient(host, searchPath, createNoValidatePath string) *UserClient {
	return &UserClient{host: host, searchPath: searchPath, createNoValidate: createNoValidatePath}
}

// SearchByMobileNumber returns (nil, nil) when no matching citizen exists —
// only a transport-level failure is a non-nil error, so callers can cleanly
// decide "not found -> create" vs "peer unreachable -> propagate".
func (c *UserClient) SearchByMobileNumber(requestInfo dto.RequestInfo, tenantID, mobileNumber string) (*User, error) {
	req := userSearchRequest{RequestInfo: requestInfo, TenantID: tenantID, MobileNumber: mobileNumber, UserType: "CITIZEN"}
	var res userSearchResponse
	if err := PostJSON(JoinURL(c.host, c.searchPath), nil, req, &res); err != nil {
		return nil, fmt.Errorf("user search: %w", err)
	}
	if len(res.User) == 0 {
		return nil, nil
	}
	return &res.User[0], nil
}

// CreateNoValidate mirrors Java's UserService.createUser: creates a CITIZEN
// account for a connection holder without OTP/duplicate validation, matching
// the internal-microservice trust boundary Java's _createnovalidate uses.
func (c *UserClient) CreateNoValidate(requestInfo dto.RequestInfo, tenantID, name, mobileNumber string) (*User, error) {
	user := User{
		Name:         name,
		UserName:     mobileNumber,
		MobileNumber: mobileNumber,
		Type:         "CITIZEN",
		TenantID:     tenantID,
		Active:       true,
		Roles:        []Role{{Code: "CITIZEN", TenantID: tenantID}},
	}

	req := createUserRequest{RequestInfo: requestInfo, User: user}
	var res createUserResponse
	if err := PostJSON(JoinURL(c.host, c.createNoValidate), nil, req, &res); err != nil {
		return nil, fmt.Errorf("user create: %w", err)
	}
	return &res.User, nil
}
