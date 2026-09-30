package repo

import (
	"server-alarm/api/internal/app/entity"
	"server-alarm/api/internal/pkg/db"
)

// ClimateRepoDB represents a DB repo for entity.Climate.
type ClimateRepoDB struct {
	dbInst *db.DB
}

// NewClimateRepoDB returns a new instance of ClimateRepoDB.
func NewClimateRepoDB(dbInst *db.DB) *ClimateRepoDB {
	return &ClimateRepoDB{
		dbInst: dbInst,
	}
}

// GetLastRecord returns the last climate db record.
func (r *ClimateRepoDB) GetLastRecord() (*entity.Climate, error) {
	record := &entity.Climate{}

	err := r.dbInst.QueryRow(`
		SELECT id, datetime, temperature1
		FROM climate
		ORDER BY id DESC
		LIMIT 1`,
	).Scan(&record.ID, &record.Datetime, &record.Temperature)

	return record, err // err OR nil
}
