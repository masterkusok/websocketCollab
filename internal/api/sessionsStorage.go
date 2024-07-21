package api

import "sync"

type SessionStorage struct {
	data    map[int]*Session
	rwMutex sync.RWMutex
}

func NewSessionStorage() *SessionStorage {
	return &SessionStorage{map[int]*Session{}, sync.RWMutex{}}
}

func (s *SessionStorage) Add(key int, session *Session) {
	s.rwMutex.Lock()
	s.data[key] = session
	s.rwMutex.Unlock()
}

func (s *SessionStorage) Get(key int) *Session {
	s.rwMutex.RLock()
	session := s.data[key]
	s.rwMutex.RUnlock()
	return session
}

func (s *SessionStorage) HasKey(key int) bool {
	s.rwMutex.RLock()
	_, ok := s.data[key]
	s.rwMutex.RUnlock()
	return ok
}

func (s *SessionStorage) Remove(key int) {
	s.rwMutex.Lock()
	delete(s.data, key)
	s.rwMutex.Unlock()
}
