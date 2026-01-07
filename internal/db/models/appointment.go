package models

import (
	"time"

	"gorm.io/gorm"
)

const (
	StatusPending   = "pending"
	StatusConfirmed = "confirmed"
	StatusCancelled = "cancelled"
	StatusCompleted = "completed"
)

type Appointment struct {
	ID uint `gorm:"primaryKey" json:"id"`

	PatientID uint `gorm:"not null" json:"patientId"`
	Patient   User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"patient,omitempty"`

	DoctorID uint   `gorm:"not null" json:"doctorId"`
	Doctor   Doctor `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"doctor,omitempty"`

	AppointmentDate time.Time `gorm:"not null" json:"appointmentDate"`
	Reason          string    `gorm:"type:text" json:"reason"`
	Status          string    `gorm:"type:varchar(20);default:'pending';check(status IN ('pending', 'confirmed', 'cancelled', 'completed'))" json:"status"`

	DoctorNotes string   `gorm:"type:text" json:"doctorNotes"`
	Medications []string `gorm:"type:jsonb;serializer:json" json:"medications"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
