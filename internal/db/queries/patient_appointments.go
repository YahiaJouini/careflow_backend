package queries

import (
	"errors"

	"github.com/YahiaJouini/careflow/internal/db"
	"github.com/YahiaJouini/careflow/internal/db/models"
)

func CreateAppointment(patientID uint, req AppointmentRequest) (*models.Appointment, error) {
	var doctor models.Doctor

	if err := db.Db.First(&doctor, req.DoctorID).Error; err != nil {
		return nil, errors.New("Doctor not found")
	}
	if !doctor.IsAvailable {
		return nil, errors.New("Doctor is currently unavailable")
	}

	appointment := models.Appointment{
		PatientID:       patientID,
		DoctorID:        req.DoctorID,
		AppointmentDate: req.AppointmentDate,
		Reason:          req.Reason,
		Status:          models.StatusPending,
	}

	if err := db.Db.Create(&appointment).Error; err != nil {
		return nil, err
	}
	return &appointment, nil
}

func GetPatientAppointments(patientID uint) ([]models.Appointment, error) {
	var appointments []models.Appointment

	err := db.Db.Preload("Doctor").Preload("Doctor.User").
		Where("patient_id = ? AND status != ?", patientID, models.StatusCompleted).
		Order("appointment_date desc").
		Find(&appointments).Error

	return appointments, err
}

func GetMedicalHistory(patientID uint) ([]models.Appointment, error) {
	var appointments []models.Appointment

	err := db.Db.Preload("Doctor").Preload("Doctor.User").
		Where("patient_id = ? AND status = ?", patientID, models.StatusCompleted).
		Order("appointment_date desc").
		Find(&appointments).Error

	return appointments, err
}

func UpdateAppointment(appointmentID uint, patientID uint, req AppointmentUpdateRequest) (*models.Appointment, error) {
	var appointment models.Appointment

	if err := db.Db.Where("id = ? AND patient_id = ?", appointmentID, patientID).First(&appointment).Error; err != nil {
		return nil, errors.New("Appointment not found")
	}

	if !appointment.AppointmentDate.Equal(req.AppointmentDate) {
		if appointment.Status == models.StatusConfirmed {
			appointment.Status = models.StatusPending
		}
	}

	appointment.AppointmentDate = req.AppointmentDate
	appointment.Reason = req.Reason

	if err := db.Db.Save(&appointment).Error; err != nil {
		return nil, err
	}

	return &appointment, nil
}

func CancelAppointment(appointmentID uint, patientID uint) error {
	var appointment models.Appointment

	if err := db.Db.Where("id = ? AND patient_id = ?", appointmentID, patientID).First(&appointment).Error; err != nil {
		return errors.New("Appointment not found")
	}

	appointment.Status = models.StatusCancelled
	return db.Db.Save(&appointment).Error
}

func DeleteAppointment(appointmentID uint, patientID uint) error {
	var appointment models.Appointment

	if err := db.Db.Where("id = ? AND patient_id = ?", appointmentID, patientID).First(&appointment).Error; err != nil {
		return errors.New("Appointment not found")
	}

	return db.Db.Delete(&appointment).Error
}
