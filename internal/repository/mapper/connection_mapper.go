// Package mapper converts between the flat wire-format dto.SewerageConnection
// and the multi-table model.SewerageConnection — the Go equivalent of Java's
// SewerageRowMapper, but built on GORM Preload instead of hand-rolled
// join-then-dedupe row scanning.
package mapper

import (
	"time"

	"github.com/google/uuid"

	"github.com/egovernments/sw-services-go/internal/domain/dto"
	"github.com/egovernments/sw-services-go/internal/domain/model"
	"github.com/egovernments/sw-services-go/internal/util"
)

func ToModel(d dto.SewerageConnection) *model.SewerageConnection {
	now := time.Now().UnixMilli()
	id := d.ID
	if id == "" {
		id = uuid.NewString()
	}

	m := &model.SewerageConnection{
		ID:                       id,
		TenantID:                 d.TenantID,
		PropertyID:               d.PropertyID,
		ApplicationNo:            d.ApplicationNo,
		ApplicationStatus:        d.ApplicationStatus,
		Status:                   d.Status,
		ConnectionNo:             d.ConnectionNo,
		OldConnectionNo:          d.OldConnectionNo,
		RoadType:                 d.RoadType,
		RoadCuttingArea:          d.RoadCuttingArea,
		ApplicationType:          d.ApplicationType,
		DateEffectiveFrom:        d.DateEffectiveFrom,
		Locality:                 d.Locality,
		OldApplication:           d.OldApplication,
		Channel:                  d.Channel,
		IsDisconnectionTemporary: d.IsDisconnectionTemporary,
		DisconnectionReason:      d.DisconnectionReason,
		AdditionalDetails:        model.JSONB(d.AdditionalDetails),
		ServiceDetail: &model.ConnectionServiceDetail{
			ConnectionID:               id,
			ConnectionExecutionDate:    d.ConnectionExecutionDate,
			DisconnectionExecutionDate: d.DisconnectionExecutionDate,
			ConnectionType:             d.ConnectionType,
			NoOfWaterClosets:           d.NoOfWaterClosets,
			NoOfToilets:                d.NoOfToilets,
			ProposedWaterClosets:       d.ProposedWaterClosets,
			ProposedToilets:            d.ProposedToilets,
			AppCreatedDate:             now,
		},
	}
	m.CreatedTime = now
	m.LastModifiedTime = now
	m.ServiceDetail.CreatedTime = now
	m.ServiceDetail.LastModifiedTime = now

	for _, doc := range d.Documents {
		docID := doc.ID
		if docID == "" {
			docID = uuid.NewString()
		}
		m.Documents = append(m.Documents, model.Document{
			ID:           docID,
			TenantID:     d.TenantID,
			DocumentType: doc.DocumentType,
			FileStoreID:  doc.FileStoreID,
			DocumentUID:  doc.DocumentUID,
			SWID:         id,
			Active:       doc.Status != "INACTIVE",
			AuditDetails: model.AuditDetails{CreatedTime: now, LastModifiedTime: now},
		})
	}

	for _, p := range d.PlumberInfo {
		pID := p.ID
		if pID == "" {
			pID = uuid.NewString()
		}
		m.PlumberInfo = append(m.PlumberInfo, model.PlumberInfo{
			ID:                    pID,
			TenantID:              d.TenantID,
			Name:                  p.Name,
			LicenseNo:             p.LicenseNo,
			MobileNumber:          p.MobileNumber,
			Gender:                p.Gender,
			FatherOrHusbandName:   p.FatherOrHusbandName,
			CorrespondenceAddress: p.CorrespondenceAddress,
			Relationship:          p.Relationship,
			SWID:                  id,
			AuditDetails:          model.AuditDetails{CreatedTime: now, LastModifiedTime: now},
		})
	}

	for _, rc := range d.RoadCuttingInfo {
		rcID := rc.ID
		if rcID == "" {
			rcID = uuid.NewString()
		}
		m.RoadCuttingInfo = append(m.RoadCuttingInfo, model.RoadCuttingInfo{
			ID:              rcID,
			TenantID:        d.TenantID,
			SWID:            id,
			Active:          rc.Status != "INACTIVE",
			RoadType:        rc.RoadType,
			RoadCuttingArea: rc.RoadCuttingArea,
			AuditDetails:    model.AuditDetails{CreatedTime: now, LastModifiedTime: now},
		})
	}

	for _, h := range d.ConnectionHolders {
		m.ConnectionHolders = append(m.ConnectionHolders, model.ConnectionHolder{
			TenantID:             d.TenantID,
			ConnectionID:         id,
			Status:               h.Status,
			UserID:               h.UUID,
			IsPrimaryHolder:      h.IsPrimaryOwner,
			ConnectionHolderType: h.OwnerType,
			HoldershipPercentage: h.OwnerShipPercentage,
			Relationship:         h.Relationship,
			AuditDetails:         model.AuditDetails{CreatedTime: now, LastModifiedTime: now},
		})
	}

	return m
}

func ToDTO(m *model.SewerageConnection) dto.SewerageConnection {
	d := dto.SewerageConnection{
		ID:                       m.ID,
		TenantID:                 m.TenantID,
		PropertyID:               m.PropertyID,
		ApplicationNo:            m.ApplicationNo,
		ApplicationStatus:        m.ApplicationStatus,
		Status:                   m.Status,
		ConnectionNo:             m.ConnectionNo,
		OldConnectionNo:          m.OldConnectionNo,
		RoadType:                 m.RoadType,
		RoadCuttingArea:          m.RoadCuttingArea,
		ApplicationType:          m.ApplicationType,
		DateEffectiveFrom:        m.DateEffectiveFrom,
		Locality:                 m.Locality,
		OldApplication:           m.OldApplication,
		Channel:                  m.Channel,
		IsDisconnectionTemporary: m.IsDisconnectionTemporary,
		DisconnectionReason:      m.DisconnectionReason,
		AdditionalDetails:        map[string]interface{}(m.AdditionalDetails),
		AuditDetails: &dto.AuditDetails{
			CreatedBy:        m.CreatedBy,
			LastModifiedBy:   m.LastModifiedBy,
			CreatedTime:      m.CreatedTime,
			LastModifiedTime: m.LastModifiedTime,
		},
	}

	if m.ServiceDetail != nil {
		d.ConnectionExecutionDate = m.ServiceDetail.ConnectionExecutionDate
		d.DisconnectionExecutionDate = m.ServiceDetail.DisconnectionExecutionDate
		d.ConnectionType = m.ServiceDetail.ConnectionType
		d.NoOfWaterClosets = m.ServiceDetail.NoOfWaterClosets
		d.NoOfToilets = m.ServiceDetail.NoOfToilets
		d.ProposedWaterClosets = m.ServiceDetail.ProposedWaterClosets
		d.ProposedToilets = m.ServiceDetail.ProposedToilets
	}

	for _, doc := range m.Documents {
		status := "ACTIVE"
		if !doc.Active {
			status = "INACTIVE"
		}
		d.Documents = append(d.Documents, dto.Document{
			ID:           doc.ID,
			TenantID:     doc.TenantID,
			DocumentType: doc.DocumentType,
			FileStoreID:  doc.FileStoreID,
			DocumentUID:  doc.DocumentUID,
			Status:       status,
		})
	}

	for _, p := range m.PlumberInfo {
		d.PlumberInfo = append(d.PlumberInfo, dto.PlumberInfo{
			ID:                    p.ID,
			Name:                  p.Name,
			LicenseNo:             p.LicenseNo,
			MobileNumber:          p.MobileNumber,
			Gender:                p.Gender,
			FatherOrHusbandName:   p.FatherOrHusbandName,
			CorrespondenceAddress: p.CorrespondenceAddress,
			Relationship:          p.Relationship,
		})
	}

	for _, rc := range m.RoadCuttingInfo {
		status := "ACTIVE"
		if !rc.Active {
			status = "INACTIVE"
		}
		d.RoadCuttingInfo = append(d.RoadCuttingInfo, dto.RoadCuttingInfo{
			ID:              rc.ID,
			RoadType:        rc.RoadType,
			RoadCuttingArea: rc.RoadCuttingArea,
			Status:          status,
		})
	}

	for _, h := range m.ConnectionHolders {
		d.ConnectionHolders = append(d.ConnectionHolders, dto.OwnerInfo{
			UUID:                h.UserID,
			IsPrimaryOwner:      h.IsPrimaryHolder,
			OwnerShipPercentage: h.HoldershipPercentage,
			OwnerType:           h.ConnectionHolderType,
			Relationship:        h.Relationship,
			Status:              h.Status,
		})
	}

	return d
}

// MergeIntoModel applies the "fetch-before-merge" sparse-update strategy: only
// fields actually present in the incoming request overwrite the existing row,
// so a partial update payload can never silently wipe columns the client
// didn't send — the exact bug class Resource/sw-services issues.pdf flagged.
func MergeIntoModel(existing *model.SewerageConnection, incoming dto.SewerageConnection) {
	now := time.Now().UnixMilli()

	if incoming.ApplicationStatus != "" {
		existing.ApplicationStatus = incoming.ApplicationStatus
	}
	if incoming.Status != "" {
		existing.Status = incoming.Status
	}
	if incoming.ConnectionNo != "" {
		existing.ConnectionNo = incoming.ConnectionNo
	}
	if incoming.OldConnectionNo != "" {
		existing.OldConnectionNo = incoming.OldConnectionNo
	}
	if incoming.RoadType != "" {
		existing.RoadType = incoming.RoadType
	}
	if incoming.RoadCuttingArea != 0 {
		existing.RoadCuttingArea = incoming.RoadCuttingArea
	}
	if incoming.ApplicationType != "" {
		existing.ApplicationType = incoming.ApplicationType
	}
	if incoming.DateEffectiveFrom != 0 {
		existing.DateEffectiveFrom = incoming.DateEffectiveFrom
	}
	if incoming.Locality != "" {
		existing.Locality = incoming.Locality
	}
	if incoming.Channel != "" {
		existing.Channel = incoming.Channel
	}
	existing.IsDisconnectionTemporary = incoming.IsDisconnectionTemporary
	if incoming.DisconnectionReason != "" {
		existing.DisconnectionReason = incoming.DisconnectionReason
	}
	if incoming.AdditionalDetails != nil {
		existing.AdditionalDetails = model.JSONB(incoming.AdditionalDetails)
	}
	existing.LastModifiedTime = now

	if existing.ServiceDetail == nil {
		existing.ServiceDetail = &model.ConnectionServiceDetail{
			ConnectionID: existing.ID,
			AuditDetails: model.AuditDetails{CreatedTime: now},
		}
	}
	if incoming.ConnectionExecutionDate != 0 {
		existing.ServiceDetail.ConnectionExecutionDate = incoming.ConnectionExecutionDate
	}
	if incoming.DisconnectionExecutionDate != 0 {
		existing.ServiceDetail.DisconnectionExecutionDate = incoming.DisconnectionExecutionDate
	}
	if incoming.ConnectionType != "" {
		existing.ServiceDetail.ConnectionType = incoming.ConnectionType
	}
	if incoming.NoOfWaterClosets != 0 {
		existing.ServiceDetail.NoOfWaterClosets = incoming.NoOfWaterClosets
	}
	if incoming.NoOfToilets != 0 {
		existing.ServiceDetail.NoOfToilets = incoming.NoOfToilets
	}
	if incoming.ProposedWaterClosets != 0 {
		existing.ServiceDetail.ProposedWaterClosets = incoming.ProposedWaterClosets
	}
	if incoming.ProposedToilets != 0 {
		existing.ServiceDetail.ProposedToilets = incoming.ProposedToilets
	}
	existing.ServiceDetail.LastModifiedTime = now

	// UnmaskingUtil mechanism (issues doc §5/§7): if the incoming plumber
	// mobileNumber is masked (contains '*' — the client is echoing back a
	// previously-masked value it never actually saw), substitute the real
	// value already on file instead of persisting asterisks over real data.
	if incoming.PlumberInfo != nil {
		incoming.PlumberInfo = resolvePlumberMasking(existing.PlumberInfo, incoming.PlumberInfo)
	}

	// Child collections use replace-if-provided semantics: an omitted/nil
	// array in the request leaves the existing rows untouched; a provided
	// array (even if some entries are pre-existing IDs) replaces the full set.
	needsFresh := incoming.Documents != nil || incoming.PlumberInfo != nil ||
		incoming.RoadCuttingInfo != nil || incoming.ConnectionHolders != nil
	if needsFresh {
		fresh := ToModel(incoming)
		if incoming.Documents != nil {
			existing.Documents = fresh.Documents
		}
		if incoming.PlumberInfo != nil {
			existing.PlumberInfo = fresh.PlumberInfo
		}
		if incoming.RoadCuttingInfo != nil {
			existing.RoadCuttingInfo = fresh.RoadCuttingInfo
		}
		if incoming.ConnectionHolders != nil {
			existing.ConnectionHolders = fresh.ConnectionHolders
		}
	}
}

// resolvePlumberMasking matches incoming plumber entries to existing ones by
// ID and swaps in the real mobileNumber wherever the incoming value looks
// masked — Java's UnmaskingUtil.getUnmaskedPlumberInfo does the same
// id-matched swap for update/modify/disconnect flows.
func resolvePlumberMasking(existing []model.PlumberInfo, incoming []dto.PlumberInfo) []dto.PlumberInfo {
	existingByID := make(map[string]model.PlumberInfo, len(existing))
	for _, p := range existing {
		existingByID[p.ID] = p
	}

	resolved := make([]dto.PlumberInfo, len(incoming))
	for i, p := range incoming {
		if trusted, ok := existingByID[p.ID]; ok {
			p.MobileNumber = util.ResolveIfMasked(p.MobileNumber, trusted.MobileNumber)
		}
		resolved[i] = p
	}
	return resolved
}
