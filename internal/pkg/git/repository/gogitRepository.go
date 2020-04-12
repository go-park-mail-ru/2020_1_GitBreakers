package repository

import (
	"fmt"
	gogit "github.com/go-git/go-git/v5"
	gogitPlumbing "github.com/go-git/go-git/v5/plumbing"
	gogitPlumbingObj "github.com/go-git/go-git/v5/plumbing/object"
	"github.com/pkg/errors"
	"strings"
)

type gogitRepository struct {
	*gogit.Repository
}

func plainOpenGoGitRepository(path string) (gogitRepo gogitRepository, err error) {
	gogitPlumbRepo, err := gogit.PlainOpen(path)
	if err != nil {
		return gogitRepo, err
	}
	gogitRepo = gogitRepository{gogitPlumbRepo}

	return gogitRepo, nil
}

// getLastChangesOfFileByPath is very slow
func (gogitRepo gogitRepository) getLastChangesOfFileByPath(filePath string,
	startCommitHash gogitPlumbing.Hash) (*gogitPlumbingObj.Commit, error) {

	if filePath == "" {
		return nil, fmt.Errorf("getLastChangesOfFileByPath called with empty path, gogitRepo=%+v",
			gogitRepo)
	}

	commitIter, err := gogitRepo.Log(&gogit.LogOptions{
		From:  startCommitHash,
		Order: gogit.LogOrderCommitterTime,
		PathFilter: func(s string) bool {
			return strings.HasPrefix(s, filePath)
		},
	})
	if err != nil {
		return nil, errors.Wrapf(err, "error while getting getLastChangesOfFileByPath in "+
			"repo=%+v, filePath=%+v", gogitRepo, filePath)
	}
	defer commitIter.Close()

	commit, err := commitIter.Next()
	if err != nil {
		return nil, errors.Wrapf(err, "error while getting getLastChangesOfFileByPath in "+
			"repo=%+v, filePath=%+v", gogitRepo, filePath)
	}
	return commit, nil
}
