package models

import (
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primary_key"`
	FirstName string         `json:"first_name" validate:"required,max=50"`
	LastName  string         `json:"last_name" validate:"required,max=50"`
	Email     string         `json:"email" validate:"required,email"`
	Password  string         `json:"-" validate:"required,min=8,max=24"`
	CreatedAt time.Time      `json:"created_at,omitempty"`
	UpdatedAt time.Time      `json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// Hashear la contraseña antes de almacenarla
func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// Verificar si la contraseña ingresada coincide con el hash almacenado
func (u *User) CheckPasswordHash(hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(u.Password))
	return err == nil
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.FirstName = cases.Title(language.Spanish).String(u.FirstName)
	u.LastName = cases.Title(language.Spanish).String(u.LastName)
	u.Email = strings.ToLower(u.Email)
	passwordHashed, err := hashPassword(u.Password)
	if err != nil {
		return err
	}
	u.Password = passwordHashed

	return
}
