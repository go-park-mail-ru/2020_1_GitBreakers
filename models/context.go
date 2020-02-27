package models

// StoresContext
type StoresContext struct {
	UsersStore UsersStore
	SessionStore SessionsStore
	RepositoryStore []Repository
}

func CreateStoresContext() *StoresContext {
	ctx := new(StoresContext)
	ctx.UsersStore.users = map[string]User{}
	ctx.UsersStore.emails = map[string]bool{}
	ctx.UsersStore.logins = map[string]bool{}
	ctx.SessionStore.sessions = map[string]Session{}
	ctx.RepositoryStore = []Repository{}
	return ctx
}