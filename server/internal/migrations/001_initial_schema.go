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

package migrations

import (
	"github.com/tiny-giraffes/life-beacon-360/server/internal/models"
	"gorm.io/gorm"
)

// InitialSchemaMigration creates the initial database schema
type InitialSchemaMigration struct{}

// ID returns the migration identifier
func (m *InitialSchemaMigration) ID() string {
	return "001_initial_schema"
}

// Up creates the initial database tables
func (m *InitialSchemaMigration) Up(db *gorm.DB) error {
	// Create locations table
	if err := db.AutoMigrate(&models.Location{}); err != nil {
		return err
	}

	return nil
}

// Down removes the initial database tables
func (m *InitialSchemaMigration) Down(db *gorm.DB) error {
	// Drop tables in reverse order
	if err := db.Migrator().DropTable(&models.Location{}); err != nil {
		return err
	}

	return nil
}
