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
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

// Migration represents a database migration
type Migration interface {
	ID() string
	Up(db *gorm.DB) error
	Down(db *gorm.DB) error
}

// SchemaMigration tracks applied migrations in the database
type SchemaMigration struct {
	ID        string    `gorm:"primaryKey"`
	AppliedAt time.Time `gorm:"autoCreateTime"`
}

// MigrationRunner handles the execution of database migrations
type MigrationRunner struct {
	db         *gorm.DB
	migrations []Migration
}

// NewMigrationRunner creates a new migration runner
func NewMigrationRunner(db *gorm.DB) *MigrationRunner {
	return &MigrationRunner{
		db:         db,
		migrations: getAllMigrations(),
	}
}

// RunMigrations executes all pending migrations
func (mr *MigrationRunner) RunMigrations() error {
	// Create schema_migrations table if it doesn't exist
	if err := mr.db.AutoMigrate(&SchemaMigration{}); err != nil {
		return fmt.Errorf("failed to create schema_migrations table: %w", err)
	}

	// Get applied migrations
	var appliedMigrations []SchemaMigration
	if err := mr.db.Find(&appliedMigrations).Error; err != nil {
		return fmt.Errorf("failed to query applied migrations: %w", err)
	}

	appliedMap := make(map[string]bool)
	for _, migration := range appliedMigrations {
		appliedMap[migration.ID] = true
	}

	// Run pending migrations
	for _, migration := range mr.migrations {
		if !appliedMap[migration.ID()] {
			log.Printf("Running migration: %s", migration.ID())

			if err := migration.Up(mr.db); err != nil {
				return fmt.Errorf("failed to run migration %s: %w", migration.ID(), err)
			}

			// Record the migration as applied
			schemaMigration := SchemaMigration{
				ID:        migration.ID(),
				AppliedAt: time.Now(),
			}

			if err := mr.db.Create(&schemaMigration).Error; err != nil {
				return fmt.Errorf("failed to record migration %s: %w", migration.ID(), err)
			}

			log.Printf("Successfully applied migration: %s", migration.ID())
		}
	}

	return nil
}

// getAllMigrations returns all available migrations in order
func getAllMigrations() []Migration {
	return []Migration{
		&InitialSchemaMigration{},
		&AddUserGroupTablesMigration{},
	}
}
