package models

type GitBranch struct {
	Name   string    `json:"name" valid:"-"`
	Commit GitCommit `json:"commit" valid:"-"`
}

type GitBranchCommits struct {
	Commits []GitCommit `json:"commits" valid:"-"`
}
