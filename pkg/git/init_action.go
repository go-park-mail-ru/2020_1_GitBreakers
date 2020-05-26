package git

import (
	gogit "github.com/go-git/go-git/v5"
)

func Init(path string, isBare bool) (Repository, error) {
	if gogitRepo, err := gogit.PlainInit(path, isBare); err != nil {
		return Repository{}, err
	} else {
		return Repository{gogitRepo}, nil
	}
}
