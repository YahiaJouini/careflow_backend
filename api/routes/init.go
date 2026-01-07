package routes

import (
	"net/http"

	"github.com/YahiaJouini/careflow/api/middleware"
	"github.com/gorilla/mux"
)

func InitializeRoutes() *mux.Router {
	router := mux.NewRouter()
	authRouter := router.PathPrefix("/auth").Subrouter()
	meRouter := router.PathPrefix("/me").Subrouter()
	meRouter.Use(middleware.AuthMiddleware(middleware.All))
	InitUserRoutes(meRouter)

	adminRouter := router.PathPrefix("/admin").Subrouter()
	adminRouter.Use(middleware.AuthMiddleware(middleware.Admin))

	patientRouter := router.PathPrefix("/patient").Subrouter()
	patientRouter.Use(middleware.AuthMiddleware(middleware.Patient))

	doctorRouter := router.PathPrefix("/doctor").Subrouter()
	doctorRouter.Use(middleware.AuthMiddleware(middleware.Doctor))

	publicRouter := router.PathPrefix("/public").Subrouter()

	InitAuthRoutes(authRouter)

	InitAdminRoutes(adminRouter)
	InitPatientRoutes(patientRouter)
	InitPublicRoutes(publicRouter)
	InitDoctorRoutes(doctorRouter)
	return router
}

func UseSecurityMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		next.ServeHTTP(w, r)
	})
}
