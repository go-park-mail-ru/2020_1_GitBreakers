package git

import "time"

type Commit struct {
	CommitHash        string    `json:"commit_hash" valid:"-"`
	CommitAuthorName  string    `json:"commit_author_name" valid:"-"`
	CommitAuthorEmail string    `json:"commit_author_email" valid:"-"`
	CommitAuthorWhen  time.Time `json:"commit_author_when" valid:"-"`
	CommitterName     string    `json:"committer_name" valid:"-"`
	CommitterEmail    string    `json:"committer_email" valid:"-"`
	CommitterWhen     time.Time `json:"committer_when" valid:"-"`
	TreeHash          string    `json:"tree_hash" valid:"-"`
	Message           string    `json:"message" valid:"-"`
	CommitParents     []string  `json:"commit_parents" valid:"-"`
}
type CommitRequest struct {
	UserLogin  string `json:"user_login"`
	RepoName   string `json:"repo_name"`
	CommitHash string `json:"commit_hash"`
	Offset     int64  `json:"offset" schema:"offset"`
	Limit      int64  `json:"limit" schema:"limit"`
}

//easyjson:json
type CommitSet []Commit
