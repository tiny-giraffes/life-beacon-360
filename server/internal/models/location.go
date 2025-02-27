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
)

type Location struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Latitude  float64   `gorm:"type:float8;not null" json:"latitude"`
	Longitude float64   `gorm:"type:float8;not null" json:"longitude"`
	CreatedAt time.Time `gorm:"type:timestamptz;default:current_timestamp;not null" json:"createdAt"`
}

type LocationRequest struct {
	Latitude  float64 `json:"latitude" validate:"required"`
	Longitude float64 `json:"longitude" validate:"required"`
}
