package queries

import (
	"errors"
	"time"

	"github.com/YahiaJouini/careflow/pkg/mails"
	"gorm.io/gorm"

	"github.com/YahiaJouini/careflow/internal/db"
	"github.com/YahiaJouini/careflow/internal/db/models"
)

func GetUserByID(userID uint) (*models.User, error) {
	var user models.User

	result := db.Db.Preload("Doctor").Preload("Doctor.Specialty").Take(&user, "id = ?", userID)
	if result.Error != nil {
		return nil, errors.New("user not found")
	}

	return &user, nil
}

func GetUserByEmail(email string) (*models.User, error) {
	var user models.User

	result := db.Db.Take(&user, "email = ?", email)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

type UpdateUserBody struct {
	// user Fields
	FirstName *string `json:"firstName" validate:"omitempty,min=3,max=30"`
	LastName  *string `json:"lastName" validate:"omitempty,min=3,max=30"`
	Image     *string `json:"image" validate:"omitempty,url"`

	// doctor fields
	Bio             *string  `json:"bio" validate:"omitempty,max=500"`
	ConsultationFee *float64 `json:"consultationFee" validate:"omitempty,gte=0"`
	IsAvailable     *bool    `json:"isAvailable" validate:"omitempty"`
}

func UpdateUser(userID uint, body UpdateUserBody) (*models.User, error) {
	var updatedUser models.User

	err := db.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Preload("Doctor").First(&updatedUser, userID).Error; err != nil {
			return err
		}

		// update basic fields
		if body.FirstName != nil {
			updatedUser.FirstName = *body.FirstName
		}
		if body.LastName != nil {
			updatedUser.LastName = *body.LastName
		}
		if body.Image != nil {
			updatedUser.Image = *body.Image
		}

		if err := tx.Save(&updatedUser).Error; err != nil {
			return err
		}

		if updatedUser.Role == "doctor" {
			if updatedUser.Doctor.UserID == 0 {
				return errors.New("doctor profile missing data integrity error")
			}

			// apply updates to the nested struct
			if body.Bio != nil {
				updatedUser.Doctor.Bio = *body.Bio
			}
			if body.ConsultationFee != nil {
				updatedUser.Doctor.ConsultationFee = *body.ConsultationFee
			}
			if body.IsAvailable != nil {
				updatedUser.Doctor.IsAvailable = *body.IsAvailable
			}

			if err := tx.Save(&updatedUser.Doctor).Error; err != nil {
				return err
			}
		}

		return nil // commit transaction
	})
	if err != nil {
		return nil, err
	}

	return GetUserByID(userID)
}

func DeleteUser(userID uint) error {
	return db.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", userID).Delete(&models.Doctor{}).Error; err != nil {
			return err
		}

		if err := tx.Delete(&models.User{}, userID).Error; err != nil {
			return err
		}

		return nil
	})
}

func CreateUser(db *gorm.DB, user *models.User) error {
	result := db.Create(user)
	return result.Error
}

// always use these in a transaction !!
func CreateDoctor(db *gorm.DB, doctor *models.Doctor) error {
	result := db.Create(doctor)
	return result.Error
}

func CreatePatient(tx *gorm.DB, patient *models.Patient) error {
	return tx.Create(patient).Error
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
