package consts

import "time"

const (
	MsgLength          = 1024
	ReadWriteMaxLength = 100
	ReadWriteMinLength = 8
	HeartBeatTime      = time.Second * 20
)

const (
	Dir = "_config"
)
