package git

import gogit "github.com/go-git/go-git/v5"

type Remote struct {
	*gogit.Remote
}
