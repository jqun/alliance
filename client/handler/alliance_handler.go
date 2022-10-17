package handler

import (
	"alliance/common/net"
	"alliance/proto-message/pb"
	"log"
)

func ScAllianceInfo(data interface{}, c net.Conner) {
	resp, ok := data.(*pb.ScAllianceInfo)
	if !ok {
		return
	}

	result := resp.Result
	switch result {
	case 2:
		log.Printf("please login first")
	case 3:
		log.Printf("sorry, you have not join an alliance")
	case 4:
		log.Printf("lang sysytem error")
	case 1:
		log.Printf("AllianceName[%v], Members[%+v]", resp.AllianceName, resp.Members)
	}
}

func ScAllianceCreate(data interface{}, c net.Conner) {
	resp, ok := data.(*pb.ScAllianceCreate)
	if !ok {
		return
	}
	result := resp.Result
	switch result {
	case 2:
		log.Printf("sorry, you already have an alliance[%v]", resp.Name)
	case 3:
		log.Printf("sorry, the allianceName[%v] already exists", resp.Name)
	case 1:
		log.Printf("create alliance[%v] success", resp.Name)
	}
}

func ScAllianceList(data interface{}, _ net.Conner) {
	resp, ok := data.(*pb.ScAllianceList)
	if !ok {
		return
	}
	log.Printf("allianceList[%+v]", resp.AllianceList)
}

func ScAllianceJoin(data interface{}, _ net.Conner) {
	resp, ok := data.(*pb.ScAllianceJoin)
	if !ok {
		return
	}
	result := resp.Result
	switch result {
	case 2:
		log.Printf("sorry, you already have an alliance[%v]", resp.AllianceName)
	case 3:
		log.Printf("sorry, the allianceName[%v] not exists", resp.AllianceName)
	case 1:
		log.Printf("join alliance[%v] success", resp.AllianceName)
	}
}

func ScAllianceDismiss(data interface{}, _ net.Conner) {
	resp, ok := data.(*pb.ScAllianceDismiss)
	if !ok {
		return
	}
	result := resp.Result
	switch result {
	case 2:
		log.Printf("sorry, you have not an alliance")
	case 3:
		log.Printf("lang system error, alliance[%v]", resp.AllianceName)
	case 4:
		log.Printf("sorry, you're not the chair of alliance[%v]", resp.AllianceName)
	case 1:
		log.Printf("dismiss alliance[%v] success", resp.AllianceName)
	}
}

func ScAllianceIncreaseCapacity(data interface{}, _ net.Conner) {
	resp, ok := data.(*pb.ScAllianceIncreaseCapacity)
	if !ok {
		return
	}
	result := resp.Result
	switch result {
	case 2:
		log.Printf("sorry, you have not an alliance")
	case 3:
		log.Printf("lang system error")
	case 4:
		log.Printf("sorry, you're not the chair")
	case 5:
		log.Printf("sorry, the alliance already increase capactiy")
	case 1:
		log.Printf("increase capactiy success")
	}
}

func ScAllianceStoreItem(data interface{}, _ net.Conner) {
	resp, ok := data.(*pb.ScAllianceStoreItem)
	if !ok {
		return
	}
	result := resp.Result
	switch result {
	case 2:
		log.Printf("sorry, you have not an alliance")
	case 3:
		log.Printf("lang system error")
	case 4:
		log.Printf("sorry, the item index max limit")
	case 1:
		log.Printf("storeItem success")
	}
}

func ScAllianceDestroyItem(data interface{}, c net.Conner) {
	resp, ok := data.(*pb.ScAllianceDestroyItem)
	if !ok {
		return
	}
	result := resp.Result
	switch result {
	case 2:
		log.Printf("sorry, you have not an alliance")
	case 3:
		log.Printf("lang system error")
	case 4:
		log.Printf("sorry, you're not the chair")
	case 1:
		log.Printf("destroy item success")
	}
}

func ScAllianceClearup(data interface{}, c net.Conner) {
	resp, ok := data.(*pb.ScAllianceClearup)
	if !ok {
		return
	}
	result := resp.Result
	switch result {
	case 2:
		log.Printf("sorry, you have not an alliance")
	case 3:
		log.Printf("lang system error")
	case 1:
		log.Printf("clearup item success")
	}
}