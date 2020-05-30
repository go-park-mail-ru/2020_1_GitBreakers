package git

import (
	"fmt"
	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	gogitConfig "github.com/go-git/go-git/v5/config"
	"github.com/pkg/errors"
)

func (repo Repository) FetchBranchForce(remoteName, branchName string, depth int) error {
	branchRef := CreateBranchRefName(branchName)
	refName := CreateRemoteRefName(remoteName, branchName)

	refSpec := config.RefSpec(fmt.Sprintf("+%s:%s", branchRef, refName))
	if err := refSpec.Validate(); err != nil {
		return errors.WithStack(err)
	}

	err := repo.Fetch(&gogit.FetchOptions{
		RemoteName: remoteName,
		RefSpecs:   []gogitConfig.RefSpec{refSpec},
		Depth:      depth,
		Tags:       gogit.NoTags,
		Force:      true,
	})

	switch {
	case err == gogit.NoErrAlreadyUpToDate:
		return nil
	case err != nil:
		return errors.WithStack(err)
	}

	return nil
}
