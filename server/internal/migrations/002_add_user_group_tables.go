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

// AddUserGroupTablesMigration adds groups and users tables and updates locations table
type AddUserGroupTablesMigration struct{}

// ID returns the migration identifier
func (m *AddUserGroupTablesMigration) ID() string {
	return "002_add_user_group_tables"
}

// Up creates the groups and users tables and adds user_id to locations
func (m *AddUserGroupTablesMigration) Up(db *gorm.DB) error {
	// Create groups table
	if err := db.AutoMigrate(&models.Group{}); err != nil {
		return err
	}

	// Create users table
	if err := db.AutoMigrate(&models.User{}); err != nil {
		return err
	}

	// Add user_id column to locations table if it doesn't exist
	if db.Migrator().HasColumn(&models.Location{}, "user_id") {
		return nil
	}

	// Add the user_id column
	if err := db.Migrator().AddColumn(&models.Location{}, "user_id"); err != nil {
		return err
	}

	// Add foreign key constraint
	if err := db.Migrator().CreateConstraint(&models.Location{}, "User"); err != nil {
		return err
	}

	return nil
}

// Down removes the groups and users tables and user_id from locations
func (m *AddUserGroupTablesMigration) Down(db *gorm.DB) error {
	// Drop foreign key constraint first
	if err := db.Migrator().DropConstraint(&models.Location{}, "User"); err != nil {
		return err
	}

	// Remove user_id column from locations
	if err := db.Migrator().DropColumn(&models.Location{}, "user_id"); err != nil {
		return err
	}

	// Drop tables in reverse order (to handle foreign key constraints)
	if err := db.Migrator().DropTable(&models.User{}); err != nil {
		return err
	}

	if err := db.Migrator().DropTable(&models.Group{}); err != nil {
		return err
	}

	return nil
}
