package routes

import (
	"github.com/YahiaJouini/careflow/api/handlers/patient"
	"github.com/gorilla/mux"
)

func InitPatientRoutes(router *mux.Router) {
	router.HandleFunc("/me", patient.GetPatient).Methods("GET")
	router.HandleFunc("/me", patient.UpdatePatient).Methods("PUT")

	router.HandleFunc("/appointments", patient.GetAppointments).Methods("GET")
	router.HandleFunc("/appointments", patient.CreateAppointment).Methods("POST")
	router.HandleFunc("/appointments/history", patient.GetMedicalHistory).Methods("GET")
	router.HandleFunc("/appointments/{id}", patient.UpdateAppointment).Methods("PUT")
	router.HandleFunc("/appointments/{id}", patient.CancelAppointment).Methods("PUT")
	router.HandleFunc("/appointments/{id}", patient.DeleteAppointment).Methods("DELETE")
}
