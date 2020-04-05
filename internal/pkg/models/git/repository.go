package git

import "time"

type Repository struct {
	Id          int       `json:"-" valid:"-"`
	OwnerId     int       `json:"-" valid:"-"`
	Name        string    `json:"name" valid:"stringlength(1|512)"`
	Description string    `json:"description" valid:"stringlength(1|2048)"`
	IsPublic    bool      `json:"is_public" valid:"-"`
	IsFork      bool      `json:"is_fork" valid:"-"`
	CreatedAt   time.Time `json:"created_at" valid:"-"`
}
