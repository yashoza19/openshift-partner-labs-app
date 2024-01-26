package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
)

// Note is used by pop to map your notes database table to your go code.
type Note struct {
	ID        int       `json:"id" db:"id"`
	LabID     int       `json:"lab_id" db:"lab_id" form:"lab_id" validate:"required"`
	UserID    string    `json:"user_id" db:"user_id" form:"user_id" validate:"required"`
	Note      string    `json:"note" db:"note" form:"note" validate:"required"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// String is not required by pop and may be deleted
func (n Note) String() string {
	jn, _ := json.Marshal(n)
	return string(jn)
}

// Notes is not required by pop and may be deleted
type Notes []Note

// String is not required by pop and may be deleted
func (n Notes) String() string {
	jn, _ := json.Marshal(n)
	return string(jn)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (n *Note) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.IntIsPresent{Field: n.LabID, Name: "LabID"},
		&validators.StringIsPresent{Field: n.UserID, Name: "UserID"},
		&validators.StringIsPresent{Field: n.Note, Name: "Note"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (n *Note) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (n *Note) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
