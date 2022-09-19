package model

import (
	"time"
	"github.com/dgrijalva/jwt-go"
)


type User struct {
	ID                   uint       `gorm:"primary_key" json:"id"`
	Name             	string     `gorm:"column:nama" json:"nama"`
	Password             string     `gorm:"column:password" json:"password"`
	Phone                string     `gorm:"column:phone" json:"phone"`
	PasswordConfirmation string     `json:"password_confirmation" gorm:"-"`
	AuthToken            string     `gorm:"column:auth_token" json:"auth_token"`
	CreatedAt            *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt            *time.Time `gorm:"column:updated_at" json:"updated_at"`
	LastLogin            time.Time  `gorm:"column:last_login" json:"last_login"`
}

type UserView struct {
	ID                   uint       `gorm:"primary_key" json:"id"`
	Name             	string     `gorm:"column:nama" json:"nama"`
	AuthToken            string     `gorm:"column:auth_token" json:"auth_token"`
}
type UserViews struct {
	ID                   uint       `gorm:"primary_key" json:"id"`
	Name             	string     `gorm:"column:nama" json:"nama"`
}

type Empty struct {

}

type Token struct {
	ID    uint   `json:"id"`
	Name    string   `json:"nama"`
	*jwt.StandardClaims
}