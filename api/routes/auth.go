package routes

import (
	"net/http"

	"github.com/YahiaJouini/chat-app-backend/api/handlers/auth"
	"github.com/YahiaJouini/chat-app-backend/api/middleware"
	"github.com/gorilla/mux"
)

func InitAuthRoutes(router *mux.Router) {
	router.HandleFunc("/login", auth.Login).Methods("POST")
	router.HandleFunc("/register", auth.Register).Methods("POST")
	router.HandleFunc("/verify-email", auth.ValidateCode).Methods("POST")
	router.Handle("/verify-email/token", middleware.SetReturnTokenMiddleware(http.HandlerFunc(auth.ValidateCode))).Methods("POST")
	router.HandleFunc("/resend-verification", auth.ResendCode).Methods("POST")

	// after login
	router.HandleFunc("/logout", auth.Logout).Methods("POST")
	router.HandleFunc("/refresh-token", auth.RefreshToken).Methods("POST")
	// check if user is authenticated
	router.Handle("/verify", middleware.AuthMiddleware(middleware.All)(http.HandlerFunc(auth.Authenticated))).Methods("GET")
}
