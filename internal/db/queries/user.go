package queries

import (
	"errors"
	"time"

	"github.com/YahiaJouini/chat-app-backend/pkg/mails"

	"github.com/YahiaJouini/chat-app-backend/internal/db"
	"github.com/YahiaJouini/chat-app-backend/internal/db/models"
)

func GetUserByID(userID uint) (*models.User, error) {
	var user models.User

	result := db.Db.Take(&user, "id = ?", userID)
	if result.Error != nil {
		return nil, errors.New("user not found")
	}

	return &user, nil
}

func GetUserByEmail(email string) (*models.User, error) {
	var user models.User

	result := db.Db.Take(&user, "email = ?", email)
	if result.Error != nil {
		return nil, errors.New("user not found")
	}

	return &user, nil
}

func MarkAsVerified(email string) error {
	result := db.Db.Model(models.User{}).Where("email = ?", email).Update("verified", true)
	return result.Error
}

func UpdateVerificationCode(email string) (string, int, error) {
	user, err := GetUserByEmail(email)
	if err != nil {
		return "", 404, err
	}
	newCode, err := mails.GenerateVerificationCode()
	if err != nil {
		return "", 500, err
	}
	expiresAt := time.Now().Add(15 * time.Minute)
	user.VerificationCode = newCode
	user.CodeExpirationTime = expiresAt
	db.Db.Save(&user)
	return newCode, 200, nil
}

func UpdateUser(userID uint, body UpdateUserBody) (*models.User, error) {
	// fetch user first (same principle as GetUserByID)
	user, err := GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	// apply only non-nil updates
	if body.FirstName != nil {
		user.FirstName = *body.FirstName
	}
	if body.LastName != nil {
		user.LastName = *body.LastName
	}
	if body.Image != nil {
		user.Image = *body.Image
	}

	// save changes
	if err := db.Db.Save(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

type UpdateUserBody struct {
	FirstName *string `json:"firstName" validate:"omitempty,min=3,max=30"`
	LastName  *string `json:"lastName" validate:"omitempty,min=3,max=30"`
	Image     *string `json:"image" validate:"omitempty,url"`
}
