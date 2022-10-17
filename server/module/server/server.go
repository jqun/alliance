package server

import (
	network "alliance/common/net"
	"alliance/server/module/role"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type serverInfo struct {
	users       sync.Map
	AllianceMgr *allianceMgr
	NetSever    *network.Server
}

var SeverInfo *serverInfo

func init() {
	SeverInfo = &serverInfo{AllianceMgr: &allianceMgr{}}
}

func (m *serverInfo) GetAllianceMgr() *allianceMgr {
	return m.AllianceMgr
}

func (m *serverInfo) GetUser(roleName string) (*role.User, bool) {
	v, ok := m.users.Load(roleName)
	if !ok {
		return nil, ok
	}
	return v.(*role.User), ok
}

func (m *serverInfo) StoreUser(user *role.User) {
	m.users.Store(user.Role.Name, user)
}

func (m *serverInfo) UserLogout(roleName string) {
	m.users.Delete(roleName)
}

func (m *serverInfo) SetNetServer(s *network.Server) {
	m.NetSever = s
}

func (m *serverInfo) Run() {
	go m.signalListen()
	m.NetSever.Run()
}

func (m *serverInfo) stop() {
	m.NetSever.Stop()
	m.AllianceMgr.stop()
}

func (m *serverInfo) signalListen() {
	log.Printf("start up signal handle!")
	ticker := time.NewTicker(time.Minute * 2)
	sNotify := make(chan os.Signal)
	signal.Notify(sNotify, os.Interrupt, os.Kill, syscall.SIGTERM)
	for {
		select {
		case <-sNotify:
			m.stop()
			return
		case <-ticker.C:
			log.Printf("signal handle running!")
		}
	}
}
