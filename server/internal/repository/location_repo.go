// Life Beacon 360
// Copyright (C) 2025 Tim Yashin/tiny-giraffes
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

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
