package models

import (
	gitmodels "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models/git"
	"time"
)

type Issue struct {
	ID        int64     `json:"id" valid:"-" db:"id"`
	AuthorID  int64     `json:"author_id" valid:"-" db:"author"`
	RepoID    int64     `json:"repo_id" valid:"-" db:"repo"`
	Title     string    `json:"title" valid:"stringlength(1|256)" db:"title"`
	Message   string    `json:"message" valid:"stringlength(1|1024)" db:"message"`
	Label     string    `json:"label" valid:"stringlength(0|50)" db:"label"`
	IsClosed  bool      `json:"is_closed" valid:"-" db:"is_closed"`
	CreatedAt time.Time `json:"created_at,omitempty" valid:"-" db:"created_at"`
}

//пока нигде не использую
type Star struct {
	AuthorID int64 `json:"-" valid:"-"`
	RepoID   int64 `json:"repo" valid:"-"`
	Vote     bool  `json:"vote" valid:"-"`
}

type News struct {
	ID       int64     `json:"id"`
	AuthorID int64     `json:"author_id"`
	Mess     string    `json:"message"`
	Date     time.Time `json:"date"`
}

//easyjson:json
type NewsSet []News

//easyjson:json
type RepoSet []gitmodels.Repository

//easyjson:json
type IssuesSet []Issue

//easyjson -all path/to/file.go
