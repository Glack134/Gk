package model

import "time"

type User struct {
	Id         int    `json:"-" db:"id"`
	Сountry    string `json:"country" binding:"required"`
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	Email      string `json:"email" binding:"required"`
	ResetToken string `json:"reset_token" db:"reset_token"`
}

type User_update struct {
	Id       int     `json:"id"`
	Country  *string `json:"country,omitempty"`
	Username *string `json:"username,omitempty"`
	Password *string `json:"password,omitempty"`
	Email    *string `json:"email,omitempty"`
}

type TokenRecord struct {
	UserID int       `json:"-" db:"id"`
	Token  string    `json:"token" db:"token"`
	Expiry time.Time `json:"expiry" db:"expiry"`
}
