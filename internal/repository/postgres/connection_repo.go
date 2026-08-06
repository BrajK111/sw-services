package postgres

import (
	"errors"

	"gorm.io/gorm"

	"github.com/egovernments/sw-services-go/internal/domain/dto"
	"github.com/egovernments/sw-services-go/internal/domain/model"
)

type ConnectionRepository struct {
	db *gorm.DB
}

func NewConnectionRepository(db *gorm.DB) *ConnectionRepository {
	return &ConnectionRepository{db: db}
}

func (r *ConnectionRepository) preloaded() *gorm.DB {
	return r.db.
		Preload("ServiceDetail").
		Preload("Documents").
		Preload("PlumberInfo").
		Preload("RoadCuttingInfo").
		Preload("ConnectionHolders")
}

func (r *ConnectionRepository) Create(conn *model.SewerageConnection) error {
	return r.db.Session(&gorm.Session{FullSaveAssociations: true}).Create(conn).Error
}

func (r *ConnectionRepository) FindByID(id string) (*model.SewerageConnection, error) {
	var conn model.SewerageConnection
	err := r.preloaded().First(&conn, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &conn, nil
}

// Update replaces child associations wholesale for any relation the caller
// populated (see mapper.MergeIntoModel) and upserts the rest — the caller is
// expected to have already merged sparse fields onto a freshly-fetched row.
func (r *ConnectionRepository) Update(conn *model.SewerageConnection) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Session(&gorm.Session{FullSaveAssociations: true}).Save(conn).Error; err != nil {
			return err
		}
		// Only replace connection holders when the caller explicitly provided
		// a non-empty slice. An omitted/nil holders array means "leave existing
		// holders untouched" — mapper.MergeIntoModel only sets ConnectionHolders
		// when the incoming request included them, so this check is safe.
		if len(conn.ConnectionHolders) > 0 {
			if err := tx.Where("connectionid = ?", conn.ID).Delete(&model.ConnectionHolder{}).Error; err != nil {
				return err
			}
			if err := tx.Create(&conn.ConnectionHolders).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *ConnectionRepository) applyFilters(q *gorm.DB, c dto.SearchCriteria) *gorm.DB {
	if c.TenantID != "" {
		q = q.Where("tenantid = ?", c.TenantID)
	}
	if len(c.IDs) > 0 {
		q = q.Where("id IN ?", c.IDs)
	}
	if len(c.ApplicationNumber) > 0 {
		q = q.Where("applicationno IN ?", c.ApplicationNumber)
	}
	if len(c.ConnectionNumber) > 0 {
		q = q.Where("connectionno IN ?", c.ConnectionNumber)
	}
	if c.OldConnectionNumber != "" {
		q = q.Where("oldconnectionno = ?", c.OldConnectionNumber)
	}
	if c.Status != "" {
		q = q.Where("status = ?", c.Status)
	}
	if len(c.ApplicationStatus) > 0 {
		q = q.Where("applicationstatus IN ?", c.ApplicationStatus)
	}
	if c.PropertyID != "" {
		q = q.Where("property_id = ?", c.PropertyID)
	}
	if c.ApplicationType != "" {
		q = q.Where("applicationtype = ?", c.ApplicationType)
	}
	if c.Locality != "" {
		q = q.Where("locality = ?", c.Locality)
	}
	if c.FromDate > 0 {
		q = q.Where("createdtime >= ?", c.FromDate)
	}
	if c.ToDate > 0 {
		q = q.Where("createdtime <= ?", c.ToDate)
	}
	return q
}

func (r *ConnectionRepository) Search(c dto.SearchCriteria) ([]model.SewerageConnection, error) {
	var connections []model.SewerageConnection
	q := r.applyFilters(r.preloaded(), c)
	if c.Limit > 0 {
		q = q.Limit(c.Limit)
	}
	if c.Offset > 0 {
		q = q.Offset(c.Offset)
	}
	if err := q.Order("createdtime desc").Find(&connections).Error; err != nil {
		return nil, err
	}
	return connections, nil
}

func (r *ConnectionRepository) Count(c dto.SearchCriteria) (int64, error) {
	var count int64
	q := r.applyFilters(r.db.Model(&model.SewerageConnection{}), c)
	if err := q.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
