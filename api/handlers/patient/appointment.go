package patient

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/YahiaJouini/careflow/api/middleware"
	"github.com/YahiaJouini/careflow/internal/db/queries"
	"github.com/YahiaJouini/careflow/pkg/auth"
	"github.com/YahiaJouini/careflow/pkg/response"
	"github.com/gorilla/mux"
)

func CreateAppointment(w http.ResponseWriter, r *http.Request) {
	claims, _ := r.Context().Value(middleware.UserClaimsKey).(*auth.Claims)

	var req queries.AppointmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	data, err := queries.CreateAppointment(claims.UserID, req)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(w, data, "Appointment request sent successfully")
}

func GetAppointments(w http.ResponseWriter, r *http.Request) {
	claims, _ := r.Context().Value(middleware.UserClaimsKey).(*auth.Claims)

	data, err := queries.GetPatientAppointments(claims.UserID)
	if err != nil {
		response.ServerError(w, "Failed to retrieve appointments")
		return
	}

	response.Success(w, data, "Appointments retrieved successfully")
}

func UpdateAppointment(w http.ResponseWriter, r *http.Request) {
	claims, _ := r.Context().Value(middleware.UserClaimsKey).(*auth.Claims)
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var req queries.AppointmentUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	data, err := queries.UpdateAppointment(uint(id), claims.UserID, req)
	if err != nil {
		response.Error(w, http.StatusNotFound, err.Error())
		return
	}

	response.Success(w, data, "Appointment updated successfully")
}

func CancelAppointment(w http.ResponseWriter, r *http.Request) {
	claims, _ := r.Context().Value(middleware.UserClaimsKey).(*auth.Claims)
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	err := queries.CancelAppointment(uint(id), claims.UserID)
	if err != nil {
		response.Error(w, http.StatusNotFound, "Appointment not found or unauthorized")
		return
	}

	response.Success(w, nil, "Appointment cancelled successfully")
}

func GetMedicalHistory(w http.ResponseWriter, r *http.Request) {
	claims, _ := r.Context().Value(middleware.UserClaimsKey).(*auth.Claims)

	data, err := queries.GetMedicalHistory(claims.UserID)
	if err != nil {
		response.ServerError(w, "Failed to retrieve medical history")
		return
	}

	response.Success(w, data, "Medical history retrieved successfully")
}

func DeleteAppointment(w http.ResponseWriter, r *http.Request) {
	claims, _ := r.Context().Value(middleware.UserClaimsKey).(*auth.Claims)
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	err := queries.DeleteAppointment(uint(id), claims.UserID)
	if err != nil {
		response.Error(w, http.StatusNotFound, "Appointment not found or unauthorized")
		return
	}

	response.Success(w, nil, "Appointment deleted successfully")
}
