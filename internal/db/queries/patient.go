package queries

import (
	"errors"
	"time"

	"github.com/YahiaJouini/careflow/internal/db"
	"github.com/YahiaJouini/careflow/internal/db/models"
)

type AppointmentRequest struct {
	DoctorID        uint      `json:"doctorId"`
	AppointmentDate time.Time `json:"appointmentDate"`
	Reason          string    `json:"reason"`
}

type AppointmentUpdateRequest struct {
	AppointmentDate time.Time `json:"appointmentDate"`
	Reason          string    `json:"reason"`
}


func GetPatientByUserID(userID uint) (*models.Patient, error) {
	var patient models.Patient
	if err := db.Db.Preload("User").Where("user_id = ?", userID).First(&patient).Error; err != nil {
		return nil, err
	}
	return &patient, nil
}

type UpdatePatientBody struct {
	Height            *float64  `json:"height" validate:"omitempty,gte=0"`
	Weight            *float64  `json:"weight" validate:"omitempty,gte=0"`
	BloodType         *string   `json:"bloodType" validate:"omitempty,oneof=A+ A- B+ B- AB+ AB- O+ O-"`
	ChronicConditions *[]string `json:"chronicConditions" validate:"omitempty"`
	Allergies         *[]string `json:"allergies" validate:"omitempty"`
	Medications       *[]string `json:"medications" validate:"omitempty"`
}

func UpdatePatient(userID uint, body UpdatePatientBody) (*models.Patient, error) {
	var patient models.Patient

	if err := db.Db.Where("user_id = ?", userID).First(&patient).Error; err != nil {
		return nil, errors.New("patient not found")
	}

	if body.Height != nil {
		patient.Height = *body.Height
	}
	if body.Weight != nil {
		patient.Weight = *body.Weight
	}
	if body.BloodType != nil {
		patient.BloodType = *body.BloodType
	}
	if body.ChronicConditions != nil {
		patient.ChronicConditions = *body.ChronicConditions
	}
	if body.Allergies != nil {
		patient.Allergies = *body.Allergies
	}
	if body.Medications != nil {
		patient.Medications = *body.Medications
	}

	if err := db.Db.Save(&patient).Error; err != nil {
		return nil, err
	}

	return &patient, nil
}
