package models

import (
	"encoding/json"
	"errors"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
)

// Lab is used by pop to map your labs database table to your go code.
type Lab struct {
	ID               int       `json:"id" db:"id"`
	ClusterID        uuid.UUID `json:"cluster_id" db:"cluster_id" form:"cluster_id"`
	GeneratedName    string    `json:"generated_name" db:"generated_name" form:"generated_name"`
	ClusterName      string    `json:"cluster_name" db:"cluster_name" form:"cluster_name"`
	OpenShiftVersion string    `json:"openshift_version" db:"openshift_version" form:"openshift_version"`
	ClusterSize      string    `json:"cluster_size" db:"cluster_size" form:"cluster_size"`
	CompanyName      string    `json:"company_name" db:"company_name" form:"company_name"`
	RequestType      string    `json:"request_type" db:"request_type" form:"request_type"`
	Partner          bool      `json:"partner" db:"partner" form:"partner"`
	Sponsor          string    `json:"sponsor" db:"sponsor" form:"sponsor"`
	CloudProvider    string    `json:"cloud_provider" db:"cloud_provider" form:"cloud_provider"`
	PrimaryFirst     string    `json:"primary_first" db:"primary_first" form:"primary_first"`
	PrimaryLast      string    `json:"primary_last" db:"primary_last" form:"primary_last"`
	PrimaryEmail     string    `json:"primary_email" db:"primary_email" form:"primary_email"`
	SecondaryFirst   string    `json:"secondary_first" db:"secondary_first" form:"secondary_first"`
	SecondaryLast    string    `json:"secondary_last" db:"secondary_last" form:"secondary_last"`
	SecondaryEmail   string    `json:"secondary_email" db:"secondary_email" form:"secondary_email"`
	Region           string    `json:"region" db:"region" form:"region"`
	AlwaysOn         bool      `json:"always_on" db:"always_on" form:"always_on"`
	ProjectName      string    `json:"project_name" db:"project_name" form:"project_name"`
	LeaseTime        string    `json:"lease_time" db:"lease_time" form:"lease_time"`
	Description      string    `json:"description" db:"description" form:"description"`
	Notes            string    `json:"notes" db:"notes" form:"notes"`
	StartDate        time.Time `json:"start_date" db:"start_date" form:"start_date"`
	EndDate          time.Time `json:"end_date" db:"end_date" form:"end_date"`
	State            string    `json:"state" db:"state" form:"state"`
	Hold             bool      `json:"hold" db:"hold" form:"hold"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// String is not required by pop and may be deleted
func (l Lab) String() string {
	jl, _ := json.Marshal(l)
	return string(jl)
}

// Labs is not required by pop and may be deleted
type Labs []Lab

// String is not required by pop and may be deleted
func (l Labs) String() string {
	jl, _ := json.Marshal(l)
	return string(jl)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (l *Lab) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (l *Lab) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	verrs := validate.NewErrors()
	l.ValidateCreation(verrs)
	return verrs, nil
}

func (l *Lab) ValidateCreation(errors *validate.Errors) {
	var err error

	l.ClusterID, err = uuid.NewV4()
	if err != nil {
		errors.Add("cluster_id", "Unable to generate cluster id")
	}

	is_alphanumeric := regexp.MustCompile(`^[a-zA-Z0-9]*$`)
	if !is_alphanumeric.MatchString(l.ClusterName) {
		errors.Add("cluster_name", "Cluster name can only contain alphanumeric characters")
	}

	if len(l.ClusterName) > 12 || len(l.ClusterName) < 4 {
		errors.Add("cluster_name", "Cluster name must be 4-12 characters")
	}

	l.GeneratedName = strings.Split(l.ClusterID.String(), "-")[0] + "-" + strings.ToLower(l.ClusterName)
	l.EndDate, err = setEndDate(l.StartDate.Format("2006-01-02"), l.LeaseTime)
	if err != nil {
		errors.Add("end_date", "Unable to set end date: "+err.Error())
	}
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (l *Lab) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	verrs := validate.NewErrors()
	l.UpdateValidation(verrs)
	return verrs, nil
}

func (l *Lab) UpdateValidation(errors *validate.Errors) {
	var err error
	re := regexp.MustCompile(`[A-Z][^A-Z]*`)

	l.ClusterID, err = uuid.NewV4()
	if err != nil {
		errors.Add("cluster_id", "Unable to generate cluster id")
	}

	labval := *l

	structType := reflect.TypeOf(labval)
	if structType.Kind() != reflect.Struct {
		errors.Add("struct_type", "Invalid type passed to IsValid")
	}

	structVal := reflect.ValueOf(labval)
	fieldNum := structVal.NumField()

	for i := 0; i < fieldNum; i++ {
		field := structVal.Field(i)
		fieldName := structType.Field(i).Name

		if fieldName == "ID" ||
			fieldName == "RequestType" ||
			fieldName == "GeneratedName" ||
			fieldName == "Partner" ||
			fieldName == "Sponsor" ||
			fieldName == "SecondaryFirst" ||
			fieldName == "SecondaryLast" ||
			fieldName == "SecondaryEmail" ||
			fieldName == "Notes" ||
			fieldName == "EndDate" ||
			fieldName == "State" ||
			fieldName == "Hold" ||
			fieldName == "AlwaysOn" ||
			fieldName == "CreatedAt" ||
			fieldName == "UpdatedAt" {
			continue
		}

		isSet := field.IsValid() && !field.IsZero()
		if !isSet {
			submatchall := re.FindAllString(fieldName, -1)
			formattedFieldName := ""
			for _, v := range submatchall {
				if v == "OpenShiftVersion" {
					formattedFieldName += "OpenShift Version"
				} else {
					strings.ToLower(v)
					formattedFieldName += v + " "
				}
			}
			errors.Add(fieldName, formattedFieldName+" is required")
		}
	}

	l.OpenShiftVersion = "ocp-" + l.OpenShiftVersion
	l.GeneratedName = strings.Split(l.ClusterID.String(), "-")[0] + "-" + strings.ToLower(l.ClusterName)
	l.EndDate, err = setEndDate(l.StartDate.Format("2006-01-02"), l.LeaseTime)
	if err != nil {
		errors.Add("end_date", "Unable to set end date: "+err.Error())
	}
}

func setEndDate(startdate string, leasetime string) (endtime time.Time, err error) {
	layout := "2006-01-02"
	var enddate time.Time
	start, err := time.Parse(layout, startdate)
	if err != nil {
		return time.Time{}, errors.New("invalid start_date")
	}

	switch leasetime {
	case "1d":
		enddate = start.AddDate(0, 0, 1)
	case "1w":
		enddate = start.AddDate(0, 0, 7)
	case "2w":
		enddate = start.AddDate(0, 0, 14)
	case "1m":
		enddate = start.AddDate(0, 1, 0)
	case "":
		enddate = start.AddDate(1, 0, 0)
	default:
		return time.Time{}, errors.New("invalid lease_time")
	}

	endDate := enddate.Format("2006-01-02 15:04:05")
	endtime, err = time.Parse("2006-01-02 15:04:05", endDate)

	return endtime, nil
}
