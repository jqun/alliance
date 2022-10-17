package net

import (
	"alliance/common/consts"
	"alliance/common/message"
	"alliance/common/util"
	"github.com/golang/protobuf/proto"
	"log"
	"net"
	"reflect"
	"time"
)

type Conner interface {
	Start()
	Stop()
	ReadMsg()
	WriteMsg()
	SetRoleName(string)
	GetRoleName() string
}

type socketConn struct {
	roleId        int64
	roleName      string
	conn          net.Conn
	In            chan *message.Message
	Out           chan []byte
	bySeverCreate bool //是否为服务器监听创建的socket,如果是将设置dead time
	close         chan struct{}
}

func NewConner(c net.Conn) *socketConn {
	m := &socketConn{
		conn:  c,
		In:    make(chan *message.Message, consts.MsgLength),
		Out:   make(chan []byte, consts.MsgLength),
		close: make(chan struct{}),
	}
	return m
}

func (m *socketConn) Start() {
	defer util.RunPaniced()
	go m.msgHandle()
	go m.WriteMsg()
	m.ReadMsg()
}

func (m *socketConn) SetRoleId(roleId int64) {
	m.roleId = roleId
}

func (m *socketConn) GetRoleId() int64 {
	return m.roleId
}

func (m *socketConn) SetRoleName(roleName string) {
	m.roleName = roleName
}

func (m *socketConn) GetRoleName() string {
	return m.roleName
}

func (m *socketConn) Stop() {
	m.waitMsgHandle()
	_ = m.conn.Close()
	close(m.close)
}

func (m *socketConn) ReadMsg() {
	for {
		buf := make([]byte, consts.ReadWriteMaxLength)
		n, err := m.conn.Read(buf)
		if err != nil {
			m.Stop()
			break
		}
		if n < consts.ReadWriteMinLength {
			log.Printf("msg not match condition!")
			continue
		}
		msg := message.BytesToMsg(buf[:n])
		m.In <- msg
		m.setReadDeadLine()
	}
}

func (m *socketConn) WriteMsg() {
	defer util.RunPaniced()
	for {
		select {
		case bytes := <-m.Out:
			_, err := m.conn.Write(bytes)
			if err != nil {
				log.Printf("write msg error[%v]", err.Error())
			}
		case <-m.close:
			break
		}

	}
}

func (m *socketConn) setReadDeadLine() {
	if !m.bySeverCreate {
		return
	}
	_ = m.conn.SetReadDeadline(time.Now().Add(10 * time.Minute))
}

func (m *socketConn) SendMsg(msgId uint64, msg proto.Message) {
	bytes := message.MsgToBytes(msgId, msg)
	m.Out <- bytes
}

func (m *socketConn) msgHandle() {
	defer util.RunPaniced()
	for {
		select {
		case msg, ok := <-m.In:
			if !ok {
				continue
			}
			m.Handle(msg)
		case <-m.close:
			break
		}
	}
}

func (m *socketConn) Handle(msg *message.Message) {
	hm, ok := GetHandlerModel(msg.MsgId)
	if !ok {
		return
	}
	data := reflect.New(hm.T.Elem()).Interface().(proto.Message)
	if err := proto.Unmarshal(msg.Data, data); err != nil {
		log.Printf("Handle err,error is [%v]", err.Error())
		return
	}
	hm.H(data, m)
}

func (m *socketConn) waitMsgHandle() {
	close(m.In)
	for {
		msg, ok := <-m.In
		if !ok {
			break
		}
		m.Handle(msg)
	}
}
