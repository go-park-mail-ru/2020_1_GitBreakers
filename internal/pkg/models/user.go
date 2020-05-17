package models

import (
	"time"
)

type User struct {
	ID        int64     `json:"id" valid:"-" db:"id"`
	Password  string    `json:"password,omitempty" valid:"stringlength(5|128)" db:"password"`
	Name      string    `json:"name" valid:"optional,stringlength(0|128)" db:"name"`
	Login     string    `json:"login" valid:"stringlength(3|128)" db:"login"`
	Image     string    `json:"image" valid:"-" db:"avatar_path"`
	Email     string    `json:"email" valid:"email,stringlength(4|128)" db:"email"`
	CreatedAt time.Time `json:"created_at" valid:"-" db:"created_at"`
}

type SignInForm struct {
	Login    string `json:"login" valid:"alphanum,stringlength(3|128)"`
	Password string `json:"password" valid:"stringlength(5|128)"`
}

//easyjson:json
type UserSet []User

//easyjson -all path/to/file.go
