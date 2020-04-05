package git

import "time"

type Repository struct {
	Id          int       `json:"-" valid:"-"`
	OwnerId     int       `json:"-" valid:"-"`
	Name        string    `json:"name" valid:"alphanum,required,stringlength(1|512)"`
	Description string    `json:"description" valid:"stringlength(0|2048)"`
	IsFork      bool      `json:"is_fork" valid:"-"`
	CreatedAt   time.Time `json:"created_at" valid:"-"`
	IsPublic    bool      `json:"is_public" valid:"-"`
}
