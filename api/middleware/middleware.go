package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/YahiaJouini/chat-app-backend/pkg/auth"
	"github.com/YahiaJouini/chat-app-backend/pkg/response"
)

type contextKey string

const UserClaimsKey contextKey = "userClaims"

type role string

const (
	Patient role = "patient"
	Doctor  role = "doctor"
	Admin   role = "admin"
	All     role = "all"
)

func AuthMiddleware(requiredRole role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.Unauthorized(w, "Missing Authorization header.")
				return
			}

			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				response.Unauthorized(w, "Authorization header must be in 'Bearer <token>' format.")
				return
			}

			claims, err := auth.VerifyToken(string(tokenParts[1]), auth.AccessToken)
			if err != nil {
				auth.Logout(w)
				response.Unauthorized(w, "Invalid or expired token")
				return
			}

			if claims.Role != string(requiredRole) && requiredRole != All {
				auth.Logout(w)
				response.Unauthorized(w, "Insufficient permissions")
				return
			}

			// Add claims to the request context
			// kinda like in express you do req.auth.userID
			ctx := context.WithValue(r.Context(), UserClaimsKey, claims)

			// Pass the request with the new context to the next handler
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func SetReturnTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if val := r.URL.Query().Get("return-token"); val != "" {
			response.Error(w, http.StatusBadRequest, "Cannot set return-token manually")
			return
		}
		ctx := context.WithValue(r.Context(), "return-token", true)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
