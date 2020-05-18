package models

import (
	gitmodels "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models/git"
	"time"
)

type Issue struct {
	ID          int64     `json:"id" valid:"-" db:"id"`
	AuthorID    int64     `json:"author_id" valid:"-" db:"author"`
	RepoID      int64     `json:"repo_id" valid:"-" db:"repo"`
	Title       string    `json:"title" valid:"stringlength(1|256)" db:"title"`
	Message     string    `json:"message" valid:"stringlength(1|1024)" db:"message"`
	Label       string    `json:"label" valid:"stringlength(0|50)" db:"label"`
	IsClosed    bool      `json:"is_closed" valid:"-" db:"is_closed"`
	CreatedAt   time.Time `json:"created_at,omitempty" valid:"-" db:"created_at"`
	AuthorLogin string    `json:"author_login" valid:"-" db:"user_login"`
	AuthorImage string    `json:"author_image" valid:"-" db:"user_avatar_path"`
}

//пока нигде не использую
type Star struct {
	AuthorID    int64  `json:"-" valid:"-"`
	RepoID      int64  `json:"repo" valid:"-"`
	Vote        bool   `json:"vote" valid:"-"`
	AuthorLogin string `json:"author_login"`
	RepoName    string `json:"repo_name" valid:"-"`
}

type News struct {
	ID          int64     `json:"id" db:"id"`
	AuthorID    int64     `json:"author_id" db:"author_id"`
	RepoID      int64     `json:"repo_id" db:"repository_id"`
	Mess        string    `json:"message" db:"message"`
	Label       string    `json:"label" db:"label"`
	Date        time.Time `json:"date" db:"created_at"`
	AuthorLogin string    `json:"author_login" db:"user_login"`
	AuthorImage string    `json:"author_image" db:"user_avatar_path"`
}

type PullRequest struct {
	ID         int64     `json:"id" db:"id"`
	AuthorId   int64     `json:"author_id" db:"author_id"`
	FromRepoID int64     `json:"from_repo_id" db:"from_repository_id"`
	ToRepoID   int64     `json:"to_repo_id" db:"to_repository_id"`
	Title      string    `json:"title" db:"title"`
	Message    string    `json:"message" db:"message"`
	Label      string    `json:"label" db:"label"`
	IsClosed   bool      `json:"is_closed" db:"is_closed"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	IsAccepted bool      `json:"is_accepted" db:"is_accepted"`
	BranchFrom string    `json:"branch_from" db:"from_repository_branch"`
	BranchTo   string    `json:"branch_to" db:"to_repository_branch"`
}
type PullReqSet []PullRequest

//easyjson:json
type NewsSet []News

//easyjson:json
type RepoSet []gitmodels.Repository

//easyjson:json
type IssuesSet []Issue

//easyjson -all path/to/file.go
