package models

type Patient struct {
	ID uint `gorm:"primaryKey" json:"id"`

	UserID uint `gorm:"unique;not null" json:"userId"`
	User   User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user,omitempty"`

	Height    float64 `gorm:"type:decimal(5,2)" json:"height"`
	Weight    float64 `gorm:"type:decimal(5,2)" json:"weight"`
	BloodType string  `gorm:"type:varchar(3)" json:"bloodType"`

	ChronicConditions []string `gorm:"type:jsonb;serializer:json" json:"chronicConditions"`
	Allergies         []string `gorm:"type:jsonb;serializer:json" json:"allergies"`
	Medications       []string `gorm:"type:jsonb;serializer:json" json:"medications"`
}
