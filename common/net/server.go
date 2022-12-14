package net

import (
	"log"
	"net"
	"sync"
)

type Server struct {
	sync.RWMutex
	Addr          string
	Listener      net.Listener
	CSocketConn   map[string]Conner // roleName -> Conner
	ListenerClose bool
}

func NewServer(addr string) *Server {
	m := &Server{}
	m.Addr = addr
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Panicf("listen addr[%v] err,error is [%v]", addr, err.Error())
	}
	m.Listener = l
	m.CSocketConn = make(map[string]Conner)
	return m
}

func (m *Server) Run() {
	for {
		conn, err := m.Listener.Accept()
		if err != nil {
			if m.ListenerClose {
				break
			} else {
				continue
			}
		}
		log.Printf("new a addr[%v] conn", conn.RemoteAddr().String())
		conner := NewConner(conn)
		conner.bySeverCreate = true
		conner.setReadDeadLine()
		go conner.Start()
	}
}

func (m *Server) StoreConner(roleName string, c Conner) {
	if roleName == "" {
		return
	}
	m.Lock()
	defer m.Unlock()
	m.CSocketConn[roleName] = c
}

func (m *Server) GetConner(roleName string) Conner {
	m.RLock()
	defer m.RUnlock()
	return m.CSocketConn[roleName]
}

func (m *Server) Stop() {
	err := m.Listener.Close()
	if err != nil {
		log.Printf("server close listen socket err,error is [%v]", err.Error())
	}
	m.ListenerClose = true
	var wg sync.WaitGroup
	m.RLock()
	for _, c := range m.CSocketConn {
		wg.Add(1)
		go func(c Conner) {
			defer wg.Done()
			c.Stop()
		}(c)
	}
	m.RUnlock()
	wg.Wait()
}
