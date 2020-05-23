package git

import (
	"fmt"
	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	gogitConfig "github.com/go-git/go-git/v5/config"
	gogitPlumbing "github.com/go-git/go-git/v5/plumbing"
	"github.com/pkg/errors"
	"path"
	"strings"
)

func (repo Repository) FetchBranchForce(remoteName, branchName string, depth int) error {
	refName := path.Clean(string(gogitPlumbing.NewRemoteReferenceName(remoteName, branchName)))

	branchRefPrefix := path.Clean(string(gogitPlumbing.NewBranchReferenceName("")))
	branchRefPrefix = strings.TrimSuffix(branchRefPrefix, "/")

	refSpec := config.RefSpec(fmt.Sprintf("+%s/*:%s", branchRefPrefix, refName))
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
	if err != nil {
		return errors.WithStack(err)
	}

	return err
}
