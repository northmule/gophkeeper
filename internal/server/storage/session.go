package storage

import (
	"sync"
	"time"
)

// Session данные по авторизованным пользователям
type Session struct {
	values map[string]time.Time
	mx     sync.RWMutex
}

type SessionManager interface {
	Add(token string, expire time.Time)
	IsValid(token string) bool
}

func NewSession() SessionManager {
	return &Session{
		values: make(map[string]time.Time),
	}
}

func (s *Session) Add(token string, expire time.Time) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.values[token] = expire
}

func (s *Session) IsValid(token string) bool {
	s.mx.RLock()
	defer s.mx.RUnlock()
	expireTime, ok := s.values[token]
	if !ok {
		return false
	}

	return expireTime.After(time.Now())
}
