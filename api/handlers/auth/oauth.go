package auth

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/YahiaJouini/careflow/internal/db"
	"github.com/YahiaJouini/careflow/internal/db/models"
	"github.com/YahiaJouini/careflow/internal/db/queries"
	"github.com/YahiaJouini/careflow/pkg/auth"
	"github.com/YahiaJouini/careflow/pkg/response"
	"gorm.io/gorm"
)

type GoogleLoginBody struct {
	AccessToken string `json:"accessToken" validate:"required"`
}

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
}

func fetchGoogleUserInfo(accessToken string) (*GoogleUserInfo, error) {
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfo GoogleUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}

func GoogleLogin(w http.ResponseWriter, r *http.Request) {
	var body GoogleLoginBody

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.ServerError(w, "Invalid request body")
		return
	}

	userInfo, err := fetchGoogleUserInfo(body.AccessToken)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "Failed to fetch Google user info: "+err.Error())
		return
	}

	if userInfo.Email == "" {
		response.Error(w, http.StatusUnauthorized, "No email received from Google")
		return
	}

	if !userInfo.VerifiedEmail {
		response.Error(w, http.StatusForbidden, "Google email not verified")
		return
	}

	user, err := queries.GetUserByEmail(userInfo.Email)
	if err != nil {
		newUser := models.User{
			FirstName: userInfo.GivenName,
			LastName:  userInfo.FamilyName,
			Email:     userInfo.Email,
			Image:     userInfo.Picture,
			Verified:  true,
			Role:      "patient",
		}

		txErr := db.Db.Transaction(func(tx *gorm.DB) error {
			if err := queries.CreateUser(tx, &newUser); err != nil {
				return err
			}
			patient := models.Patient{
				UserID: newUser.ID,
			}
			if err := queries.CreatePatient(tx, &patient); err != nil {
				return err
			}
			return nil
		})
		if txErr != nil {
			response.ServerError(w, "Could not create user: "+txErr.Error())
			return
		}
		user = &newUser
	} else {
		updateUserBody := queries.UpdateUserBody{
			Image: &userInfo.Picture,
		}
		queries.UpdateUser(user.ID, updateUserBody)
	}


	refreshToken := auth.GenerateToken(user, auth.RefreshToken)
	accessToken := auth.GenerateToken(user, auth.AccessToken)

	userAgent := r.Header.Get("User-Agent")
	if strings.Contains(userAgent, "Android") {
		mobileData := auth.MobileAuthResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			User:         user,
		}
		response.Success(w, mobileData, "Google Login successful")
		return
	}

	auth.SetAuthCookie(w, refreshToken, auth.Add)
	response.Success(w, accessToken, "Google Login successful")
}
