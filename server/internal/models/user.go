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

package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	GroupID      uuid.UUID `gorm:"type:uuid;not null" json:"group_id"`
	Username     string    `gorm:"type:varchar(100);unique;not null" json:"username"`
	Email        string    `gorm:"type:varchar(255)" json:"email"`
	PasswordHash string    `gorm:"type:varchar(255);not null" json:"-"`
	Role         string    `gorm:"type:varchar(50);default:'member';not null" json:"role"` // 'admin' or 'member'
	CreatedAt    time.Time `gorm:"type:timestamptz;default:current_timestamp;not null" json:"created_at"`
	UpdatedAt    time.Time `gorm:"type:timestamptz;default:current_timestamp;not null" json:"updated_at"`
	
	// Relations
	Group     Group      `gorm:"foreignKey:GroupID" json:"group,omitempty"`
	Locations []Location `gorm:"foreignKey:UserID" json:"locations,omitempty"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

// IsAdmin checks if the user has admin role
func (u *User) IsAdmin() bool {
	return u.Role == "admin"
}