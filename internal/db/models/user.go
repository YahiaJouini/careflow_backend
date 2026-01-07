package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	FirstName string `gorm:"type:varchar(100); not null" json:"firstName"`
	LastName  string `gorm:"type:varchar(100); not null" json:"lastName"`
	Email     string `gorm:"type:varchar(255); not null; unique" json:"email"`
	Image     string `gorm:"type:varchar(255); default:https://avatar.iran.liara.run/public" json:"image"`
	Password  string `gorm:"type:varchar(255)" json:"-"`
	Verified  bool   `gorm:"type:boolean; default:false"`
	Role      string `gorm:"type:varchar(255); not null; default:'patient'; check(role IN ('admin', 'doctor','patient'))" json:"role"`

	Doctor             *Doctor        `json:"doctor,omitempty" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Patient            *Patient       `json:"patient,omitempty" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	VerificationCode   string         `gorm:"type:varchar(6)" json:"-"`
	CodeExpirationTime time.Time      `gorm:"type:timestamp; not null" json:"-"`
	CreatedAt          time.Time      `json:"createdAt"`
	UpdatedAt          time.Time      `json:"-"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`
}
