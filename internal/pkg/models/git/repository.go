package git

import "time"

type Repository struct {
	ID          int64     `json:"id" valid:"-"`
	OwnerID     int64     `json:"owner_id" valid:"-"`
	Name        string    `json:"name" valid:"alphanum,required,stringlength(1|512)"`
	Description string    `json:"description" valid:"stringlength(0|2048)"`
	IsFork      bool      `json:"is_fork" valid:"-"`
	CreatedAt   time.Time `json:"created_at" valid:"-"`
	IsPublic    bool      `json:"is_public" valid:"-"`
	Stars       int64     `json:"stars" valid:"-"`
}

//easyjson:json
type RepositorySet []Repository
