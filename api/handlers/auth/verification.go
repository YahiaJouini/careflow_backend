package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/YahiaJouini/chat-app-backend/internal/db/queries"
	"github.com/YahiaJouini/chat-app-backend/pkg/auth"
	"github.com/YahiaJouini/chat-app-backend/pkg/mails"
	"github.com/YahiaJouini/chat-app-backend/pkg/response"
)

type ValidateCodeReq struct {
	Code  string `json:"code" validate:"required,len=6"`
	Email string `json:"email" validate:"required,email"`
}

type ResendCodeReq struct {
	Email string `json:"email" validate:"required,email"`
}

func ValidateCode(w http.ResponseWriter, r *http.Request) {
	var body ValidateCodeReq

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, 0, err.Error())
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
	if user.Verified {
		response.Error(w, 0, "Email already verified")
		return
	}

	if time.Now().After(user.CodeExpirationTime) {
		response.Error(w, 0, "Verification code has expired")
		return
	}

	if user.VerificationCode != body.Code {
		response.Error(w, 0, "Invalid verification code")
		return
	}

	err = queries.MarkAsVerified(body.Email)
	if err != nil {
		fmt.Println(err)
		response.ServerError(w)
		return
	}

	userAgent := r.Header.Get("User-Agent")

	refreshToken := auth.GenerateToken(user, auth.RefreshToken)
	accessToken := auth.GenerateToken(user, auth.AccessToken)

	if strings.Contains(userAgent, "Android") || strings.Contains(userAgent, "CareFlow") {

		mobileData := auth.MobileAuthResponse{
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

func ResendCode(w http.ResponseWriter, r *http.Request) {
	var body ResendCodeReq

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := Validate.Struct(body); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := queries.GetUserByEmail(body.Email)
	if err != nil {
		response.Error(w, http.StatusNotFound, err.Error())
		return
	}
	if user.Verified {
		response.Error(w, 0, "Email already verified")
		return
	}
	code, status, err := queries.UpdateVerificationCode(body.Email)
	if err != nil {
		response.Error(w, status, err.Error())
		return
	}

	if result := mails.SendMail(body.Email, code); result.Err != nil {
		response.ServerError(w, "An error occured sending your verification code")
		return
	}

	response.Success(w, nil, "New verification code sent")
}
