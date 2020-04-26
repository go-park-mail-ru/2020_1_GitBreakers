package gitserver

type UseCase interface {
	HaveAccess (currentUserId *int64, userLogin string, repoName string) (bool, error)
}
