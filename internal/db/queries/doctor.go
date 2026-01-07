package queries

import (
	"errors"
	"time"

	"github.com/YahiaJouini/careflow/internal/db"
	"github.com/YahiaJouini/careflow/internal/db/models"
)

type DoctorUpdateAppointmentRequest struct {
	AppointmentDate time.Time `json:"appointmentDate"`
	DoctorNotes     string    `json:"doctorNotes"`
}

type ValidateAppointmentRequest struct {
	Status string `json:"status"` // "confirmed" or "completed"
}

type PatientDetailsResponse struct {
	FirstName         string   `json:"firstName"`
	LastName          string   `json:"lastName"`
	Email             string   `json:"email"`
	Image             string   `json:"image"`
	Height            float64  `json:"height"`
	Weight            float64  `json:"weight"`
	BloodType         string   `json:"bloodType"`
	ChronicConditions []string `json:"chronicConditions"`
	Allergies         []string `json:"allergies"`
	Medications       []string `json:"medications"`
}

func getDoctorID(userID uint) (uint, error) {
	var doctor models.Doctor
	if err := db.Db.Where("user_id = ?", userID).First(&doctor).Error; err != nil {
		return 0, errors.New("Doctor profile not found")
	}
	return doctor.ID, nil
}

func GetDoctorAppointments(userID uint) ([]models.Appointment, error) {
	doctorID, err := getDoctorID(userID)
	if err != nil {
		return nil, err
	}

	var appointments []models.Appointment
	err = db.Db.Preload("Patient").
		Where("doctor_id = ?", doctorID).
		Find(&appointments).Error

	return appointments, err
}

func ValidateAppointment(userID uint, appointmentID uint, status string) (*models.Appointment, error) {
	doctorID, err := getDoctorID(userID)
	if err != nil {
		return nil, err
	}

	var appt models.Appointment
	if err := db.Db.Where("id = ? AND doctor_id = ?", appointmentID, doctorID).First(&appt).Error; err != nil {
		return nil, errors.New("Appointment not found")
	}

	if status != models.StatusConfirmed && status != models.StatusCompleted {
		return nil, errors.New("Invalid status. Use 'confirmed' or 'completed'")
	}

	appt.Status = status
	if err := db.Db.Save(&appt).Error; err != nil {
		return nil, err
	}

	return &appt, nil
}

func UpdateAppointmentDoctor(userID uint, appointmentID uint, req DoctorUpdateAppointmentRequest) (*models.Appointment, error) {
	doctorID, err := getDoctorID(userID)
	if err != nil {
		return nil, err
	}

	var appt models.Appointment
	if err := db.Db.Where("id = ? AND doctor_id = ?", appointmentID, doctorID).First(&appt).Error; err != nil {
		return nil, errors.New("Appointment not found")
	}

	if !req.AppointmentDate.IsZero() {
		appt.AppointmentDate = req.AppointmentDate
	}
	appt.DoctorNotes = req.DoctorNotes

	if err := db.Db.Save(&appt).Error; err != nil {
		return nil, err
	}

	return &appt, nil
}

func CancelAppointmentDoctor(userID uint, appointmentID uint) error {
	doctorID, err := getDoctorID(userID)
	if err != nil {
		return err
	}

	var appt models.Appointment
	if err := db.Db.Where("id = ? AND doctor_id = ?", appointmentID, doctorID).First(&appt).Error; err != nil {
		return errors.New("Appointment not found")
	}

	appt.Status = models.StatusCancelled
	return db.Db.Save(&appt).Error
}

func GetDoctorPatients(userID uint) ([]models.User, error) {
	doctorID, err := getDoctorID(userID)
	if err != nil {
		return nil, err
	}

	var patients []models.User

	err = db.Db.Model(&models.User{}).
		Joins("JOIN appointments ON appointments.patient_id = users.id").
		Where("appointments.doctor_id = ?", doctorID).
		Group("users.id").
		Find(&patients).Error

	return patients, err
}

func GetDoctorPatientDetails(doctorUserID, patientUserID uint) (*models.Patient, error) {
	doctorID, err := getDoctorID(doctorUserID)
	if err != nil {
		return nil, err
	}

	var patient models.Patient
	err = db.Db.Preload("User").
		Joins("JOIN appointments ON appointments.patient_id = patients.user_id").
		Where("appointments.doctor_id = ? AND patients.user_id = ?", doctorID, patientUserID).
		First(&patient).Error

	if err != nil {
		return nil, errors.New("patient not found or not associated with this doctor")
	}

	return &patient, nil
}
