package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/YahiaJouini/careflow/internal/db"
	"github.com/YahiaJouini/careflow/internal/db/models"
	"github.com/YahiaJouini/careflow/internal/db/queries"
	"github.com/YahiaJouini/careflow/pkg/auth"
	"github.com/YahiaJouini/careflow/pkg/mails"
	"github.com/YahiaJouini/careflow/pkg/response"
	"github.com/YahiaJouini/careflow/pkg/utils"
	"gorm.io/gorm"
)

type RegisterBody struct {
	FirstName string `json:"firstName" validate:"required,min=3,max=30"`
	LastName  string `json:"lastName" validate:"required,min=3,max=30"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6"`
	Role      string `json:"role" validate:"omitempty,oneof=patient doctor"`

	// Doctor-specific fields
	SpecialtyID   *uint   `json:"specialtyId"`
	LicenseNumber *string `json:"licenseNumber"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	var body RegisterBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// validate req body
	if err := utils.Validate.Struct(body); err != nil {
		response.Error(w, 0, err.Error())
		return
	}

	// default to patient if empty or invalid
	if body.Role == "" {
		body.Role = "patient"
	}

	if body.Role == "doctor" {
		if body.SpecialtyID == nil || body.LicenseNumber == nil {
			response.Error(w, http.StatusBadRequest, "Doctors must provide a SpecialtyID and LicenseNumber")
			return
		}
	}

	var existingUser models.User
	if result := db.Db.Select("id").Take(&existingUser, "email = ?", body.Email); result.Error == nil {
		response.Error(w, http.StatusConflict, "Account with this email already exists.")
		return
	}

	hashedPassword, err := auth.HashPassword(body.Password)
	if err != nil {
		response.ServerError(w)
		return
	}

	verificationCode, err := mails.GenerateVerificationCode()
	if err != nil {
		response.ServerError(w)
		return
	}

	expiresAt := time.Now().Add(15 * time.Minute)

	user := models.User{
		FirstName:          body.FirstName,
		LastName:           body.LastName,
		Email:              body.Email,
		Password:           hashedPassword,
		Role:               body.Role,
		VerificationCode:   verificationCode,
		CodeExpirationTime: expiresAt,
		Verified:           false,
	}

	err = db.Db.Transaction(func(tx *gorm.DB) error {
		// create the User
		if err := queries.CreateUser(tx, &user); err != nil {
			return err
		}

		if body.Role == "doctor" {
			doctor := models.Doctor{
				UserID:        user.ID,
				SpecialtyID:   *body.SpecialtyID,
				LicenseNumber: *body.LicenseNumber,
			}
			if err := queries.CreateDoctor(tx, &doctor); err != nil {
				return err
			}
		} else if body.Role == "patient" {
			patient := models.Patient{
				UserID: user.ID,
			}
			if err := queries.CreatePatient(tx, &patient); err != nil {
				return err
			}
		}

		return nil // commit
	})
	if err != nil {
		response.ServerError(w, "Database error: "+err.Error())
		return
	}

	if result := mails.SendMail(body.Email, verificationCode); result.Err != nil {
		response.ServerError(w, "An error occured sending your verification code")
		return
	}

	response.Success(w, nil, "Please validate your email")
}
