package model

type User struct {
	Id       int    `json:"-" db:"id"`
	Сountry  string `json:"country" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

type User_update struct {
	Id       int     `json:"id" db:"-"`
	Country  *string `json:"country,omitempty"`
	Username *string `json:"username,omitempty"`
	Password *string `json:"password,omitempty"`
	Email    *string `json:"email,omitempty"`
}
