package models

type GitRepository struct {
	Id          int
	OwnerId     int
	Name        string
	Description string
	IsFork      bool
}
