package user

import (
	"encoding/json"
	"net/http"

	"github.com/YahiaJouini/chat-app-backend/api/middleware"
	"github.com/YahiaJouini/chat-app-backend/internal/db/queries"
	"github.com/YahiaJouini/chat-app-backend/pkg/auth"
	"github.com/YahiaJouini/chat-app-backend/pkg/response"
	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

func GetUser(w http.ResponseWriter, r *http.Request) {
	claims, _ := r.Context().Value(middleware.UserClaimsKey).(*auth.Claims)

	// fetch user again because jwt claims aren't the source of truth
	data, err := queries.GetUserByID(claims.UserID)
	if err != nil {
		auth.SetAuthCookie(w, "", auth.Remove)
		response.Unauthorized(w, "User not found. Logged out.")
		return
	}

	response.Success(w, data, "User retrieved successfully")
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(middleware.UserClaimsKey).(*auth.Claims)

	var body queries.UpdateUserBody

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.ServerError(w, err.Error())
		return
	}

	// validate req body
	if err := Validate.Struct(body); err != nil {
		response.Error(w, 0, err.Error())
		return
	}

	updatedUser, err := queries.UpdateUser(claims.UserID, body)
	if err != nil {
		response.Error(w, 0, err.Error())
		return
	}

	response.Success(w, updatedUser, "User updated successfully")
}
