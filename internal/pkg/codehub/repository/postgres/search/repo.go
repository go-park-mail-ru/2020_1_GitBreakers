package search

import "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"

type RepoSearch struct {
}

func (repo RepoSearch) GetFromUsers(query string, limit int64, offset int64) (models.UserSet, error) {
	return models.UserSet{}, nil
}

func (repo RepoSearch) GetFromStarredRepos(query string, limit int64, offset int64) (models.RepoSet, error) {
	return models.RepoSet{}, nil
}

func (repo RepoSearch) GetFromAllRepos(query string, limit int64, offset int64) (models.RepoSet, error) {
	return models.RepoSet{}, nil
}

func (repo RepoSearch) GetFromOwnRepos(query string, limit int64, offset int64) (models.RepoSet, error) {
	return models.RepoSet{}, nil
}
