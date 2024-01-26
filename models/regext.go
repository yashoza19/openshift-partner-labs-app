package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
)

// Regext is used by pop to map your regexts database table to your go code.
type Regext struct {
	ID          int       `json:"id" db:"id"`
	LabID       int       `json:"lab_id" db:"lab_id"`
	CurrentUser string    `json:"current_user" db:"current_user"`
	Extension   string    `json:"extension" db:"extension"`
	Date        time.Time `json:"date" db:"date"`
	Status      string    `json:"status" db:"status"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// String is not required by pop and may be deleted
func (r Regext) String() string {
	jr, _ := json.Marshal(r)
	return string(jr)
}

// Regexts is not required by pop and may be deleted
type Regexts []Regext

// String is not required by pop and may be deleted
func (r Regexts) String() string {
	jr, _ := json.Marshal(r)
	return string(jr)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (r *Regext) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (r *Regext) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (r *Regext) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
