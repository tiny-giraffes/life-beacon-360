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
