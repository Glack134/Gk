package model

import "time"

type User struct {
	Id             int    `json:"-" db:"id"`
	Сountry        string `json:"country" binding:"required"`
	Username       string `json:"username" binding:"required"`
	Password       string `json:"password" binding:"required"`
	Email          string `json:"email" binding:"required"`
	ResetToken     string `json:"reset_token" db:"reset_token"`
	TwoFASecret    string `json:"two_fa_secret" db:"two_fa_secret"`
	IsTwoFAEnabled bool   `json:"is_two_fa_enabled" db:"is_two_fa_enabled"`
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

type Subscription struct {
	ID        int       `json:"id"`
	Plan      string    `json:"plan"`
	Status    string    `json:"status"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}
