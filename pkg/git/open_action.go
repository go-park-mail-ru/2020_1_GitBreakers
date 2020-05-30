package git

import gogit "github.com/go-git/go-git/v5"

func OpenRepository(path string) (Repository, error) {
	if gogitRepo, err := gogit.PlainOpen(path); err != nil {
		return Repository{}, err
	} else {
		return Repository{gogitRepo}, nil
	}
}
