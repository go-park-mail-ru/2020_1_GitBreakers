package models

type User struct {
	ID       int    `json:"id" valid:"-" db:"id"`
	Password string `json:"password,omitempty" valid:"stringlength(5|128)" db:"password"`
	Name     string `json:"name" valid:"optional,alphanum" db:"name"`
	Login    string `json:"login" valid:"stringlength(3|128)" db:"login"`
	Image    string `json:"image" valid:"-" db:"avatar_path"`
	Email    string `json:"email" valid:"email,stringlength(4|128)" db:"email"`
}

type SignInForm struct {
	Login    string `json:"login" valid:"alphanum,stringlength(3|128)"`
	Password string `json:"password" valid:"stringlength(5|128)"`
}
