package entity

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type Auditable struct {
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

// Implement the Valuer and Scanner interfaces for Auditable
func (a Auditable) Value() (driver.Value, error) {
    return json.Marshal(a)
}

func (a *Auditable) Scan(value interface{}) error {
    return json.Unmarshal(value.([]byte), a)
}

func NewAuditable() Auditable {
	return Auditable{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: gorm.DeletedAt{},
	}
}

func UpdateAuditable() Auditable {
	return Auditable{
		UpdatedAt: time.Now(),
	}
}