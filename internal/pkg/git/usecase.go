package git

type UseCase interface {
	Create()
	Update()
	GetRepo()
	GetRepoList()
}
