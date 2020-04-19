package models

import "time"

type Issue struct {
	ID        int       `json:"id" valid:"-" db:"id"`
	AuthorID  int       `json:"author" valid:"-" db:"author"`
	RepoID    int       `json:"repo" valid:"-" db:"repo"`
	Title     string    `json:"title" valid:"stringlength(1|256)" db:"title"`
	Message   string    `json:"message" valid:"stringlength(1|1024)" db:"message"`
	Label     string    `json:"label" valid:"stringlength(0|50)" db:"label"`
	IsClosed  bool      `json:"is_closed" valid:"-" db:"is_closed"`
	CreatedAt time.Time `json:"omitempty" valid:"-" db:"created_at"`
}

//пока нигде не использую
type Star struct {
	AuthorID int  `json:"-" valid:"-"`
	RepoID   int  `json:"repo" valid:"-"`
	Vote     bool `json:"vote" valid:"-"`
}
