package models

import "time"

type GitCommit struct {
	CommitHash        string    `json:"commit_hash" valid:"-"`
	CommitAuthorName  string    `json:"commit_author_name" valid:"-"`
	CommitAuthorEmail string    `json:"commit_author_email" valid:"-"`
	CommitAuthorWhen  time.Time `json:"commit_author_when" valid:"-"`
	CommitterName     string    `json:"committer_name" valid:"-"`
	CommitterEmail    string    `json:"committer_email" valid:"-"`
	CommitterWhen     time.Time `json:"committer_when" valid:"-"`
	TreeHash          string    `json:"tree_hash" valid:"-"`
	CommitParents     []string  `json:"commit_parents" valid:"-"`
}
