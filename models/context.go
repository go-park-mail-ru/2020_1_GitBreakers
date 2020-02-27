package models

// StoresContext
type StoresContext struct {
	UsersStore UsersStore
	SessionStore SessionsStore
	RepositoryStore []Repository
}
