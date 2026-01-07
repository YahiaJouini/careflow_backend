package db

import (
	"fmt"
	"log"

	"github.com/YahiaJouini/careflow/internal/db/models"
)

func Migrate() {
	err := Db.AutoMigrate(
		&models.User{},
		&models.Specialty{},
		&models.Doctor{},
		&models.Appointment{},
		&models.Patient{},
	)
	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	seedSpecialties()
	fmt.Println("migrations and seeding applied successfully")
}

func seedSpecialties() {
	var count int64
	Db.Model(&models.Specialty{}).Count(&count)

	if count == 0 {
		specialties := []models.Specialty{
			{
				Name:        "Generalist",
				Description: "Primary care and general health",
				Icon:        "https://cdn-icons-png.flaticon.com/512/2966/2966356.png",
			},
			{
				Name:        "Dentist",
				Description: "Teeth and gum health",
				Icon:        "https://cdn-icons-png.flaticon.com/512/2966/2966486.png",
			},
			{
				Name:        "Cardiologist",
				Description: "Heart and cardiovascular system",
				Icon:        "https://cdn-icons-png.flaticon.com/512/2966/2966334.png",
			},
			{
				Name:        "Neurologist",
				Description: "Disorders of the nervous system",
				Icon:        "https://cdn-icons-png.flaticon.com/512/2966/2966467.png",
			},
			{
				Name:        "Ophthalmologist",
				Description: "Eye and vision care",
				Icon:        "https://cdn-icons-png.flaticon.com/512/2966/2966367.png",
			},
			{
				Name:        "Pediatrician",
				Description: "Medical care for infants, children, and adolescents",
				Icon:        "https://cdn-icons-png.flaticon.com/512/2966/2966356.png",
			},
			{
				Name:        "Psychiatrist",
				Description: "Mental health and behavioral disorders",
				Icon:        "https://cdn-icons-png.flaticon.com/512/2966/2966356.png",
			},
		}

		if err := Db.Create(&specialties).Error; err != nil {
			log.Println("Failed to seed specialties:", err)
		} else {
			fmt.Println("Default specialties seeded.")
		}
	}
}
