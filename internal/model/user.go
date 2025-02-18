package model

type User struct {
	Id       int    `json:"-" db:"id"`
	Сountry  string `json:"country" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
}
