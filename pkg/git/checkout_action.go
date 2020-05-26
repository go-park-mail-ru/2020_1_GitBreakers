package git

import (
	gogit "github.com/go-git/go-git/v5"
	gogitPlumbing "github.com/go-git/go-git/v5/plumbing"
)

func (repo Repository) Checkout(refName gogitPlumbing.ReferenceName, isForce bool) error {
	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	err = worktree.Checkout(&gogit.CheckoutOptions{
		Branch: refName,
		Force:  isForce,
	})

	return err
}
