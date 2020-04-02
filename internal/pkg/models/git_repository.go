package models

import "time"

type GitRepository struct {
	Id          int       `json:"id" valid:"-"`
	OwnerId     int       `json:"owner_id" valid:"-"`
	Name        string    `json:"name" valid:"stringlength(1|512)"`
	Description string    `json:"description" valid:"stringlength(1|2048)"`
	IsFork      bool      `json:"is_fork" valid:"-"`
	CreatedAt   time.Time `json:"created_at" valid:"-"`
}
