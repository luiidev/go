package models

import (
	"strings"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gorm.io/gorm"
)

type User struct {
	Model
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email" gorm:"unique"`
	Password  string `json:"-"`
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
func (u *User) CheckPasswordHash(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
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
