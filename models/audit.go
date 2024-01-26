package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
)

// Audit is used by pop to map your audits database table to your go code.
type Audit struct {
	ID            int       `json:"id" db:"id"`
	GeneratedName string    `json:"generated_name" db:"generated_name"`
	AccessTime    time.Time `json:"access_time" db:"access_time"`
	LoginName     string    `json:"login_name" db:"login_name"`
	LoginType     string    `json:"login_type" db:"login_type"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// String is not required by pop and may be deleted
func (a Audit) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// Audits is not required by pop and may be deleted
type Audits []Audit

// String is not required by pop and may be deleted
func (a Audits) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (a *Audit) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (a *Audit) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (a *Audit) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
