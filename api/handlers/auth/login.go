package auth

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/YahiaJouini/chat-app-backend/internal/db/queries"
	"github.com/YahiaJouini/chat-app-backend/pkg/auth"
	"github.com/YahiaJouini/chat-app-backend/pkg/response"
	"github.com/go-playground/validator/v10"
)

type LoginBody struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password"`
}

var Validate = validator.New()

func Login(w http.ResponseWriter, r *http.Request) {
	var body LoginBody

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.ServerError(w, err.Error())
		return
	}
	//  validate req body
	if err := Validate.Struct(body); err != nil {
		response.Error(w, 0, err.Error())
		return
	}

	user, err := queries.GetUserByEmail(body.Email)
	if err != nil {
		response.Error(w, http.StatusNotFound, err.Error())
		return
	}

	if !user.Verified {
		response.Error(w, http.StatusForbidden, "User not verified yet")
		return
	}

	ok := auth.VerifyPassword(body.Password, user.Password)
	if !ok {
		response.Error(w, http.StatusNotFound, "Invalid credentials")
		return
	}

	refreshToken := auth.GenerateToken(user, auth.RefreshToken)
	accessToken := auth.GenerateToken(user, auth.AccessToken)

	userAgent := r.Header.Get("User-Agent")

	if strings.Contains(userAgent, "Android") || strings.Contains(userAgent, "CareFlow") {

		mobileData := struct {
			AccessToken  string      `json:"accessToken"`
			RefreshToken string      `json:"refreshToken"`
			User         interface{} `json:"user"`
		}{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			User:         user,
		}

		response.Success(w, mobileData, "Login successful")
		return
	}

	// assign cookies
	auth.SetAuthCookie(w, refreshToken, auth.Add)
	response.Success(w, accessToken, "Login successful")
}

func Logout(w http.ResponseWriter, r *http.Request) {
	auth.SetAuthCookie(w, "", auth.Remove)
	response.Success(w, nil, "Logged out successfully")
}
