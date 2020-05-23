package git

import gogit "github.com/go-git/go-git/v5"

type Repository struct {
	*gogit.Repository
}
