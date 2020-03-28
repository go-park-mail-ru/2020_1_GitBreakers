package models

type User struct {
	Id       int    `json:"id"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Login    string `json:"login"`
	Image    string `json:"image"`
	Email    string `json:"email"`
}
type SignInForm struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
