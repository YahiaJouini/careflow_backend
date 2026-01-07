package doctor

import (
	"net/http"
	"strconv"

	"github.com/YahiaJouini/careflow/api/middleware"
	"github.com/YahiaJouini/careflow/internal/db/queries"
	"github.com/YahiaJouini/careflow/pkg/auth"
	"github.com/YahiaJouini/careflow/pkg/response"
	"github.com/gorilla/mux"
)

func GetPatients(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserClaimsKey).(*auth.Claims)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	patients, err := queries.GetDoctorPatients(claims.UserID)
	if err != nil {
		response.ServerError(w, err.Error())
		return
	}

	response.Success(w, patients, "Patients retrieved successfully")
}



func GetPatientDetails(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserClaimsKey).(*auth.Claims)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	patientIDStr := vars["id"]
	patientID, err := strconv.ParseUint(patientIDStr, 10, 32)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid patient ID")
		return
	}

	patient, err := queries.GetDoctorPatientDetails(claims.UserID, uint(patientID))
	if err != nil {
		response.Error(w, http.StatusNotFound, err.Error())
		return
	}

	resp := queries.PatientDetailsResponse{
		FirstName:         patient.User.FirstName,
		LastName:          patient.User.LastName,
		Email:             patient.User.Email,
		Image:             patient.User.Image,
		Height:            patient.Height,
		Weight:            patient.Weight,
		BloodType:         patient.BloodType,
		ChronicConditions: patient.ChronicConditions,
		Allergies:         patient.Allergies,
		Medications:       patient.Medications,
	}

	response.Success(w, resp, "Patient details retrieved successfully")
}
