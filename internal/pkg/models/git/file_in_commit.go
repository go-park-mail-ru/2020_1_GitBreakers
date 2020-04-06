package git

import "github.com/go-git/go-git/v5/plumbing/filemode"

type FileMode filemode.FileMode

func (fm FileMode) String() string {
	switch filemode.FileMode(fm) {
	case filemode.Empty:
		return "empty"
	case filemode.Dir:
		return "dir"
	case filemode.Regular:
		return "regular"
	case filemode.Deprecated:
		return "deprecated"
	case filemode.Executable:
		return "executable"
	case filemode.Symlink:
		return "symlink"
	case filemode.Submodule:
		return "submodule"
	default:
		return "unknown"
	}
}

type FileInCommit struct {
	Name        string `json:"name"`
	FileType    string `json:"file_type"`
	FileMode    string `json:"file_mode"`
	FileSize    int64  `json:"file_size"`
	ContentType string `json:"content_type"` // github.com/h2non/filetype, if FileType != 'blob' then this field will be empty
	EntryHash   string `json:"entry_hash"`
}
type FilesCommitRequest struct {
	UserName    string `json:"user_name"`
	Reponame    string `json:"reponame"`
	HashCommits string `json:"hash_commits"`
	Path        string `schema:"path"`
}
