package git

import (
	gogitPlumbing "github.com/go-git/go-git/v5/plumbing"
	"github.com/pkg/errors"
)

func (repo Repository) MergeBase(firstCommitHash, secondCommitHash string) ([]string, error) {
	gogitFirstCommitHash := gogitPlumbing.NewHash(firstCommitHash)
	gogitSecondCommitHash := gogitPlumbing.NewHash(secondCommitHash)

	firstCommit, err := repo.CommitObject(gogitFirstCommitHash)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	secondCommit, err := repo.CommitObject(gogitSecondCommitHash)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	mergeBases, err := firstCommit.MergeBase(secondCommit)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	mergeBasesHashes := make([]string, 0, len(mergeBases))

	for i := range mergeBases {
		mergeBasesHashes = append(mergeBasesHashes, mergeBases[i].Hash.String())
	}

	return mergeBasesHashes, nil
}
