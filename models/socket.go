package models

import (
	"sync"

	socketio "github.com/googollee/go-socket.io"
)

type SocketConn struct {
	sync.RWMutex
	Conns map[string]socketio.Conn
}

func NewSocketConn() *SocketConn {
	return &SocketConn{
		Conns: make(map[string]socketio.Conn),
	}
}

func (s *SocketConn) Get(id string) socketio.Conn {
	s.RLock()
	defer s.RUnlock()
	return s.Conns[id]
}

func (s *SocketConn) Delete(id string) {
	s.Lock()
	defer s.Unlock()
	delete(s.Conns, id)
}

func (s *SocketConn) GetAll() map[string]socketio.Conn {
	s.RLock()
	defer s.RUnlock()
	return s.Conns
}

func (s *SocketConn) Set(id string, conn socketio.Conn) {
	s.Lock()
	defer s.Unlock()
	s.Conns[id] = conn
}
