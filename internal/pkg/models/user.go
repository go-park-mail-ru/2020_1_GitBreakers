package models

type User struct {
	Id       int    `json:"id" valid:"-"`
	Password string `json:"password,omitempty" valid:"stringlength(5|128)"`
	Name     string `json:"name" valid:"optional,alphanum"`
	Login    string `json:"login" valid:"stringlength(3|128)"`
	Image    string `json:"image" valid:"-"`
	Email    string `json:"email" valid:"email,stringlength(5|128)"`
}

type SignInForm struct {
	Login    string `json:"login" valid:"alphanum,stringlength(3|128)"`
	Password string `json:"password" valid:"stringlength(5|128)"`
}
