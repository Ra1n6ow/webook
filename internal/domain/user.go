package domain

import "time"

type User struct {
	Id       int64     `json:"id"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	Birthday time.Time `json:"birthday"`
	Nickname string    `json:"nickname"`
	Intro    string    `json:"intro"`
}
