package entity

import "github.com/google/uuid"

type User struct {
	ID uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Username string `json:"username" gorm:"unique;not null"`
	Email string `json:"email" gorm:"unique;not null"`
	Password string `json:"-" gorm:"not null"`
	Role string `json:"role" gorm:"default:user"`
	Verified bool `json:"verified" gorm:"default:false"`
}