package routes

import (
	"github.com/YahiaJouini/careflow/api/handlers/doctor"
	"github.com/gorilla/mux"
)

func InitDoctorRoutes(router *mux.Router) {
	router.HandleFunc("/stats", doctor.GetDashboardOverview).Methods("GET")

	// appointments routes
	router.HandleFunc("/appointments", doctor.GetAppointments).Methods("GET")
	router.HandleFunc("/appointments/{id}/validate", doctor.ValidateAppointment).Methods("PUT")
	router.HandleFunc("/appointments/{id}", doctor.UpdateAppointment).Methods("PUT")
	router.HandleFunc("/appointments/{id}", doctor.CancelAppointment).Methods("DELETE")

	// patients routes
	router.HandleFunc("/patients", doctor.GetPatients).Methods("GET")
}
