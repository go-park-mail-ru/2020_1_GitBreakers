package git

const (
	TypeDir  = "dir"
	TypeFile = "file"
	TypeExec = "exec"
)

type FileInCommit struct {
	Name        string
	Type        string
	ContentType string // github.com/h2non/filetype, if Type == 'dir' then this field will be empty
	Sha1Hash    string
}
