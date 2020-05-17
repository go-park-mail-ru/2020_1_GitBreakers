package git

import "time"

type ParentRepositoryInfo struct {
	ID          *int64  `json:"id,omitempty" valid:"-"`
	OwnerID     *int64  `json:"owner_id,omitempty" valid:"-"`
	Name        *string `json:"name,omitempty" valid:"-"`
	AuthorLogin *string `json:"author_login,omitempty" valid:"-"`
}

type Repository struct {
	ID                   int64                `json:"id" valid:"-"`
	OwnerID              int64                `json:"owner_id" valid:"-"`
	Name                 string               `json:"name" valid:"alphanum,required,stringlength(1|512)"`
	Description          string               `json:"description" valid:"stringlength(0|2048)"`
	IsFork               bool                 `json:"is_fork" valid:"-"`
	CreatedAt            time.Time            `json:"created_at" valid:"-"`
	IsPublic             bool                 `json:"is_public" valid:"-"`
	Stars                int64                `json:"stars" valid:"-"`
	Forks                int64                `json:"forks" valid:"-"`
	AuthorLogin          string               `json:"author_login,omitempty" valid:"-"`
	ParentRepositoryInfo ParentRepositoryInfo `json:"parent_repository_info,omitempty" valid:"-"`
}
type RepoFork struct {
	FromRepoID     int64  `json:"from_repo_id" valid:"-"`
	FromAuthorName string `json:"from_author_name" valid:"-"`
	FromRepoName   string `json:"from_repo_name" valid:"-"`
	NewName        string `json:"new_name" valid:"stringlength(1|512)"`
}

//easyjson:json
type RepositorySet []Repository
