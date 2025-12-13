package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Username    string    `gorm:"type:varchar(100);not null;unique" json:"username"`
	FullName    string    `gorm:"type:varchar(100)" json:"full_name"`
	Email       string    `gorm:"type:varchar(255);unique;not null" json:"email"`
	Password    string    `gorm:"type:varchar(255);not null" json:"-"`
	PhoneNumber string    `gorm:"type:varchar(20)" json:"phone_number"`
	Gender      string    `gorm:"type:varchar(15)" json:"gender"`
	Address     string    `gorm:"type:varchar(100)" json:"address"`
	SystemRole  string    `gorm:"type:varchar(50);not null;default:'USER'" json:"system_role"`
	IsActive    *bool     `gorm:"not null;default:true" json:"is_active"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (u *User) TableName() string {
	return "users"
}
