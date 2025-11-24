package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/YahiaJouini/chat-app-backend/internal/db"
	"github.com/YahiaJouini/chat-app-backend/internal/db/models"
	"github.com/YahiaJouini/chat-app-backend/pkg/auth"
	"github.com/YahiaJouini/chat-app-backend/pkg/mails"
	"github.com/YahiaJouini/chat-app-backend/pkg/response"
)

type RegisterBody struct {
	FirstName string `json:"firstName" validate:"required,min=3,max=30"`
	LastName  string `json:"lastName" validate:"required,min=3,max=30"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	var body RegisterBody

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	// validate req body
	if err := Validate.Struct(body); err != nil {
		response.Error(w, 0, err.Error())
		return
	}
	var user models.User
	if result := db.Db.Take(&user, "email = ?", body.Email); result.Error == nil {
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
	user = models.User{
		FirstName:          body.FirstName,
		LastName:           body.LastName,
		Email:              body.Email,
		Password:           hashedPassword,
		VerificationCode:   verificationCode,
		CodeExpirationTime: expiresAt,
	}

	if result := mails.SendMail(body.Email, verificationCode); result.Err != nil {
		response.ServerError(w, "An error occured sending your verification code")
		return
	}

	result := db.Db.Create(&user)
	if result.Error != nil {
		response.ServerError(w)
		return
	}
	response.Success(w, nil, "Please validate your email")
}
