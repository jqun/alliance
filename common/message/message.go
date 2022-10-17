package message

import (
	"alliance/common/consts"
	"bytes"
	"encoding/binary"
	"github.com/golang/protobuf/proto"
	"log"
)

type Message struct {
	MsgId int64  //消息id
	Data  []byte //数据
}

func BytesToMsg(data []byte) *Message {
	return &Message{
		MsgId: decodeMsgId(data[:consts.ReadWriteMinLength]),
		Data:  data[consts.ReadWriteMinLength:],
	}
}

func MsgToBytes(msgId uint64, msgData proto.Message) []byte {
	dataBytes, err := proto.Marshal(msgData)
	if err != nil {
		log.Printf("proto marshal err,error is [%v]", err.Error())
	}
	idBytes := encodeMsgId(msgId)
	b := bytes.NewBuffer([]byte{})
	b.Write(idBytes)
	b.Write(dataBytes)
	return b.Bytes()
}

func decodeMsgId(msgBytes []byte) int64 {
	var id int64
	b := bytes.NewBuffer(msgBytes)
	err := binary.Read(b, binary.BigEndian, &id)
	if err != nil {
		log.Printf("decode msg id error[%v]", err.Error())
	}
	return id
}

func encodeMsgId(msgId uint64) []byte {
	b := bytes.NewBuffer([]byte{})
	err := binary.Write(b, binary.BigEndian, msgId)
	if err != nil {
		log.Printf("encode msg id error[%v]", err.Error())
	}
	return b.Bytes()
}
