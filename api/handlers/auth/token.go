package auth

import (
	"net/http"

	"github.com/YahiaJouini/chat-app-backend/internal/db/queries"
	"github.com/YahiaJouini/chat-app-backend/pkg/auth"
	"github.com/YahiaJouini/chat-app-backend/pkg/response"
)

func RefreshToken(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetRefreshToken(r)
	if err != nil {
		response.Unauthorized(w, err.Error())
		return
	}

	claims, err := auth.VerifyToken(token, auth.RefreshToken)
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
		auth.Logout(w)
		response.Unauthorized(w, "User role has changed. Please log in again.")
		return
	}

	accessToken := auth.GenerateToken(user, auth.AccessToken)
	response.Success(w, accessToken, "Access token refreshed successfully")
}
