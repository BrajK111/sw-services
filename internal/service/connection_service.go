// Package service mirrors SewerageServiceImpl at the CRUD-orchestration
// level. Workflow transitions, IDGen calls, peer-service validation
// (MDMS/property), encryption, and Kafka side-effects are Phase 2+ — this
// phase wires the shape (validate -> enrich -> persist -> respond) with
// TODOs marking exactly where those calls plug in later.
package service

import (
	"log"

	"github.com/egovernments/sw-services-go/internal/config"
	"github.com/egovernments/sw-services-go/internal/domain/dto"
	swerrors "github.com/egovernments/sw-services-go/internal/domain/errors"
	"github.com/egovernments/sw-services-go/internal/domain/model"
	"github.com/egovernments/sw-services-go/internal/repository/mapper"
	"github.com/egovernments/sw-services-go/internal/repository/postgres"
	"github.com/egovernments/sw-services-go/internal/util"
	"github.com/egovernments/sw-services-go/internal/validator"
)

// nextStatus computes the next applicationStatus for a given role.
// It mirrors the DIGIT sewerage workflow state machine:
//
//	CITIZEN / SW_CEMP  →  PENDING_FOR_FIELD_INSPECTION
//	SW_FIELD_INSPECTOR →  PENDING_FOR_APPROVAL
//	SW_APPROVER        →  APPROVED  (also activates the connection)
//	SUPERUSER          →  accepts any status in the incoming request unchanged
func nextStatus(currentStatus string, callerRoles []string) (string, bool) {
	roleSet := make(map[string]bool, len(callerRoles))
	for _, r := range callerRoles {
		roleSet[r] = true
	}

	// SUPERUSER may set any status directly — do not override.
	if roleSet["SUPERUSER"] {
		return "", false // caller controls the value
	}
	if roleSet["SW_APPROVER"] {
		return string(model.AppStatusApproved), true
	}
	if roleSet["SW_FIELD_INSPECTOR"] {
		return string(model.AppStatusPendingForApproval), true
	}
	if roleSet["SW_CEMP"] || roleSet["EMPLOYEE"] {
		return string(model.AppStatusPendingForFieldInspection), true
	}
	// Default: move to the next logical step.
	switch currentStatus {
	case string(model.AppStatusInitiated):
		return string(model.AppStatusPendingForFieldInspection), true
	default:
		return "", false
	}
}

type ConnectionService struct {
	repo       *postgres.ConnectionRepository
	pagination config.PaginationConfig

	propertyValidator *validator.PropertyValidator
	mdmsValidator     *validator.MDMSValidator
	idgen             *util.IdGenClient
	userClient        *util.UserClient
}

func NewConnectionService(
	repo *postgres.ConnectionRepository,
	pagination config.PaginationConfig,
	propertyValidator *validator.PropertyValidator,
	mdmsValidator *validator.MDMSValidator,
	idgen *util.IdGenClient,
	userClient *util.UserClient,
) *ConnectionService {
	return &ConnectionService{
		repo:              repo,
		pagination:        pagination,
		propertyValidator: propertyValidator,
		mdmsValidator:     mdmsValidator,
		idgen:             idgen,
		userClient:        userClient,
	}
}

func (s *ConnectionService) validateCreate(c dto.SewerageConnection) error {
	if c.TenantID == "" {
		return swerrors.New("EG_SW_TENANTID_MANDATORY", "tenantId is mandatory")
	}
	if c.PropertyID == "" {
		return swerrors.New("EG_SW_PROPERTYID_MANDATORY", "propertyId is mandatory")
	}
	return nil
}

// resolveConnectionHolders mirrors Java's UserService.createUser: for each
// connection holder with a mobile number but no resolved UUID, search
// egov-user, and create a CITIZEN account if none exists. A peer-service
// outage is logged, not fatal — the connection still gets created without a
// resolved UUID for that holder, same non-blocking resilience pattern used
// elsewhere in this service.
func (s *ConnectionService) resolveConnectionHolders(requestInfo dto.RequestInfo, tenantID string, holders []dto.OwnerInfo) {
	for i := range holders {
		h := &holders[i]
		if h.UUID != "" || h.MobileNumber == "" {
			continue
		}
		existing, err := s.userClient.SearchByMobileNumber(requestInfo, tenantID, h.MobileNumber)
		if err != nil {
			log.Printf("WARN: user search for mobileNumber=%s failed (continuing without UUID): %v", h.MobileNumber, err)
			continue
		}
		if existing != nil {
			h.UUID = existing.UUID
			continue
		}
		created, err := s.userClient.CreateNoValidate(requestInfo, tenantID, h.Name, h.MobileNumber)
		if err != nil {
			log.Printf("WARN: user create for mobileNumber=%s failed (continuing without UUID): %v", h.MobileNumber, err)
			continue
		}
		h.UUID = created.UUID
	}
}

func (s *ConnectionService) CreateConnection(req dto.SewerageConnectionRequest) (dto.SewerageConnection, error) {
	if err := s.validateCreate(req.SewerageConnection); err != nil {
		return dto.SewerageConnection{}, err
	}

	c := req.SewerageConnection

	if err := s.propertyValidator.ValidatePropertyForConnection(req.RequestInfo, c.TenantID, c.PropertyID); err != nil {
		return dto.SewerageConnection{}, err
	}
	if err := s.mdmsValidator.ValidateConnectionType(req.RequestInfo, c.TenantID, c.ConnectionType); err != nil {
		return dto.SewerageConnection{}, err
	}

	if c.ApplicationNo == "" {
		generated, err := s.idgen.GenerateApplicationNumber(req.RequestInfo, c.TenantID)
		if err != nil {
			log.Printf("WARN: idgen unreachable, falling back to local application number: %v", err)
			generated = util.LocalApplicationNumber()
		}
		c.ApplicationNo = generated
	}
	if c.ApplicationStatus == "" {
		c.ApplicationStatus = "INITIATED"
	}
	if c.Status == "" {
		c.Status = "INACTIVE"
	}

	s.resolveConnectionHolders(req.RequestInfo, c.TenantID, c.ConnectionHolders)

	entity := mapper.ToModel(c)
	if err := s.repo.Create(entity); err != nil {
		return dto.SewerageConnection{}, err
	}
	return mapper.ToDTO(entity), nil
}

func (s *ConnectionService) UpdateConnection(req dto.SewerageConnectionRequest) (dto.SewerageConnection, error) {
	incoming := req.SewerageConnection
	if incoming.ID == "" {
		return dto.SewerageConnection{}, swerrors.New("EG_SW_ID_MANDATORY", "connection id is mandatory for update")
	}

	existing, err := s.repo.FindByID(incoming.ID)
	if err != nil {
		return dto.SewerageConnection{}, err
	}
	if existing == nil {
		return dto.SewerageConnection{}, swerrors.New("EG_SW_CONNECTION_NOT_FOUND", "no sewerage connection exists for the given id")
	}

	mapper.MergeIntoModel(existing, incoming)

	// Prioritize explicitly requested status (e.g. REJECTED, CANCELLED) over automatic role progression.
	if incoming.ApplicationStatus != "" && (incoming.ApplicationStatus == string(model.AppStatusRejected) || incoming.ApplicationStatus == "CANCELLED") {
		existing.ApplicationStatus = incoming.ApplicationStatus
		existing.Status = string(model.StatusInactive)
	} else if next, ok := nextStatus(existing.ApplicationStatus, req.CallerRoles); ok {
		existing.ApplicationStatus = next
	} else if incoming.ApplicationStatus != "" {
		existing.ApplicationStatus = incoming.ApplicationStatus
	}

	// On approval: activate the connection.
	if existing.ApplicationStatus == string(model.AppStatusApproved) {
		existing.Status = string(model.StatusActive)
	}

	if req.DisconnectRequest {
		existing.Status = string(model.StatusInactive)
		existing.ApplicationStatus = string(model.AppStatusPendingApprovalForDisconnect)
	}


	if err := s.repo.Update(existing); err != nil {
		return dto.SewerageConnection{}, err
	}

	refreshed, err := s.repo.FindByID(existing.ID)
	if err != nil {
		return dto.SewerageConnection{}, err
	}
	return mapper.ToDTO(refreshed), nil
}


func (s *ConnectionService) normalizePagination(c *dto.SearchCriteria) {
	if c.Limit <= 0 {
		c.Limit = s.pagination.DefaultLimit
	}
	if c.Limit > s.pagination.MaxLimit {
		c.Limit = s.pagination.MaxLimit
	}
	if c.Offset < 0 {
		c.Offset = s.pagination.DefaultOffset
	}
}

func (s *ConnectionService) Search(criteria dto.SearchCriteria) ([]dto.SewerageConnection, int64, error) {
	s.normalizePagination(&criteria)

	entities, err := s.repo.Search(criteria)
	if err != nil {
		return nil, 0, err
	}
	total, err := s.repo.Count(criteria)
	if err != nil {
		return nil, 0, err
	}

	results := make([]dto.SewerageConnection, 0, len(entities))
	for i := range entities {
		results = append(results, mapper.ToDTO(&entities[i]))
	}
	return results, total, nil
}

// PlainSearch matches Java's /_plainsearch route shape today. Java's real
// skip-level-search optimization (bypassing privacy/decryption) requires the
// ABAC context this Go service doesn't have until encryption work lands, so
// this intentionally behaves the same as Search in Phase 1 (documented, not
// silently wrong the way the audited prior attempt was).
func (s *ConnectionService) PlainSearch(criteria dto.SearchCriteria) ([]dto.SewerageConnection, int64, error) {
	return s.Search(criteria)
}
