package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/YahiaJouini/chat-app-backend/internal/db/queries"
	"github.com/YahiaJouini/chat-app-backend/pkg/auth"
	"github.com/YahiaJouini/chat-app-backend/pkg/response"
)

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

func RefreshToken(w http.ResponseWriter, r *http.Request) {
	var currentRefreshToken string
	var err error

	userAgent := r.Header.Get("User-Agent")
	isMobile := strings.Contains(userAgent, "Android")

	if isMobile {
		// mobile: extract from request body
		var body RefreshTokenRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		fmt.Println("Token : ", body.RefreshToken)
		currentRefreshToken = body.RefreshToken
	} else {
		currentRefreshToken, err = auth.GetRefreshToken(r)
		if err != nil {
			response.Unauthorized(w, err.Error())
			return
		}
	}

	claims, err := auth.VerifyToken(currentRefreshToken, auth.RefreshToken)
	if err != nil {
		response.Unauthorized(w, err.Error())
		return
	}
	user, err := queries.GetUserByID(claims.UserID)
	if err != nil {
		response.Unauthorized(w, err.Error())
		return
	}

	if claims.Role != user.Role {
		if !isMobile {
			auth.Logout(w)
		}
		response.Unauthorized(w, "User role has changed. Please log in again.")
		return
	}

	newAccessToken := auth.GenerateToken(user, auth.AccessToken)
	newRefreshToken := auth.GenerateToken(user, auth.RefreshToken)

	if isMobile {
		data := auth.MobileAuthResponse{
			AccessToken:  newAccessToken,
			RefreshToken: newRefreshToken,
			User:         user,
		}
		response.Success(w, data, "Access token refreshed successfully")
	} else {
		auth.SetAuthCookie(w, newRefreshToken, auth.Add)
		response.Success(w, newAccessToken, "Access token refreshed successfully")
	}
}
