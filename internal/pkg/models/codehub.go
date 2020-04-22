package models

import (
	gitmodels "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models/git"
	"time"
)

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

//easyjson:json
type RepoSet []gitmodels.Repository

//easyjson:json
type IssuesSet []Issue

//easyjson -all path/to/file.go
