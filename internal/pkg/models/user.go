package models

type User struct {
	Id       string `json:"id"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Login    string `json:"login"`
	Image    string `json:"image"`
	Email    string `json:"email"`
}
