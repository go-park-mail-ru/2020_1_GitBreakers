package search

import "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"

type RepoSearch struct {
}

func (repo RepoSearch) GetFromUsers(query string, limit int64, offset int64) (models.UserSet, error) {
	panic("implement me")
}

func (repo RepoSearch) GetFromStarredRepos(query string, limit int64, offset int64) (models.RepoSet, error) {
	panic("implement me")
}

func (repo RepoSearch) GetFromAllRepos(query string, limit int64, offset int64) (models.RepoSet, error) {
	panic("implement me")
}

func (repo RepoSearch) GetFromOwnRepos(query string, limit int64, offset int64) (models.RepoSet, error) {
	panic("implement me")
}
