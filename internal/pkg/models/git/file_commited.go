package git

type FileCommitted struct {
	FileInfo FileInCommit `json:"file_info"`
	Content  string       `json:"content"`
}
