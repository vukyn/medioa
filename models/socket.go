package models

import (
	"sync"

	"github.com/olahol/melody"
)

type SocketConn struct {
	sync.RWMutex
	conns map[string]*melody.Session
}

func NewSocketConn() *SocketConn {
	return &SocketConn{
		conns: make(map[string]*melody.Session),
	}
}

func (s *SocketConn) Get(id string) *melody.Session {
	s.RLock()
	defer s.RUnlock()
	return s.conns[id]
}

func (s *SocketConn) Delete(id string) {
	s.Lock()
	defer s.Unlock()
	delete(s.conns, id)
}

func (s *SocketConn) GetAll() map[string]*melody.Session {
	s.RLock()
	defer s.RUnlock()
	return s.conns
}

func (s *SocketConn) Set(id string, conn *melody.Session) {
	s.Lock()
	defer s.Unlock()
	s.conns[id] = conn
}
