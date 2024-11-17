// internal/repository/location_repo.go

package repository

import (
	"github.com/tiny-giraffes/life-beacon-360/server/internal/models"
	"gorm.io/gorm"
)

// SaveCoordinate saves a Coordinate record to the database
func SaveCoordinate(db *gorm.DB, coord *models.Location) error {
	return db.Create(coord).Error
}
