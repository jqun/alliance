package util

import (
	"log"
	"runtime/debug"
	"strconv"
	"strings"
)

func RunPaniced(){
	if err := recover(); err != nil {
		log.Printf("Err is [%v]\n Stack:[%v]", err, string(debug.Stack()))
	}
}

func SplitSpace(str string) []string {
	return strings.Fields(strings.TrimSpace(str))
}

func StrToInt32(str string) int32 {
	v, _ := strconv.ParseInt(str, 10, 64)
	return int32(v)
}