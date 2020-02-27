package models

import (
	"sync"
)

//var sessionsStore SessionsStore

type Session struct {
	ID        string `json:"session_id"`
	UserLogin string `json:"-"`
}

type SessionsStore struct {
	sessions map[string]Session
	mu       *sync.Mutex
}

const (
	_ = iota
	StatusSessionAlreadyExist
	StatusSessionNotExist
)

func (store *SessionsStore) StoreSession(session Session) CommonError {
	defer store.mu.Unlock()
	store.mu.Lock()

	if _, isUnique := store.sessions[session.ID]; isUnique {
		return NewModelError(`session with same id already in SessionsStore`, StatusSessionAlreadyExist)
	}
	store.sessions[session.ID] = session
	return nil
}

func (store *SessionsStore) HaveSession(sessionId string) bool {
	_, in := store.sessions[sessionId]
	return in
}

func (store *SessionsStore) GetSession(sessionId string) (Session, CommonError) {
	session, inStore := store.sessions[sessionId]
	if !inStore {
		return session, NewModelError(`session not exits`, StatusSessionNotExist)
	}
	return session, nil
}

func (store *SessionsStore) DeleteSession(sessionId string) {
	defer store.mu.Unlock()
	store.mu.Lock()

	if _, inStore := store.sessions[sessionId]; inStore {
		delete(store.sessions, sessionId)
	}
}
