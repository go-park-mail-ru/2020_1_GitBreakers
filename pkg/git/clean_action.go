package git

import gogit "github.com/go-git/go-git/v5"

func (repo Repository) Clean() error {
	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	err = worktree.Clean(&gogit.CleanOptions{
		Dir: true,
	})

	return err
}
