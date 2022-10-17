package handler

import (
	"alliance/common/net"
	"alliance/proto-message/pb"
	pbId "alliance/proto-message/pb-id"
	"alliance/server/module/server"
)

func CsAllianceInfo(data interface{}, c net.Conner) {
	_, ok := data.(*pb.CsAllianceInfo)
	if !ok {
		return
	}
	roleName := c.GetRoleName()
	user, ok := server.SeverInfo.GetUser(roleName)
	ret := &pb.ScAllianceInfo{Result: 1}
	if !ok {
		ret.Result = 2
		net.SendMsg(pbId.ScAllianceInfoId, ret, c)
		return
	}

	allianceName := user.Role.AllianceName
	if allianceName == "" {
		ret.Result = 3
		net.SendMsg(pbId.ScAllianceInfoId, ret, c)
		return
	}

	allianceInfo := server.SeverInfo.GetAllianceMgr().GetAlliance(allianceName)
	if allianceInfo == nil {
		ret.Result = 4
		net.SendMsg(pbId.ScAllianceInfoId, ret, c)
		return
	}
	ret.AllianceName = allianceName
	ret.Members = make([]string, len(allianceInfo.Members))
	for m := range allianceInfo.Members {
		ret.Members = append(ret.Members, m)
	}
	net.SendMsg(pbId.ScAllianceInfoId, ret, c)
}

func CsAllianceCreate(data interface{}, c net.Conner) {
	req, ok := data.(*pb.CsAllianceCreate)
	if !ok {
		return
	}
	allianceName := req.Name
	roleName := c.GetRoleName()
	user, ok := server.SeverInfo.GetUser(roleName)
	ret := &pb.ScAllianceCreate{Result: 1}
	if !ok || user.Role.AllianceName != "" {
		ret.Result = 2
		ret.Name = user.Role.AllianceName
		net.SendMsg(pbId.ScAllianceCreateId, ret, c)
		return
	}

	allianceInfo := server.SeverInfo.GetAllianceMgr().GetAlliance(allianceName)
	if allianceInfo != nil {
		ret.Result = 3
		ret.Name = allianceName
		net.SendMsg(pbId.ScAllianceCreateId, ret, c)
		return
	}

	allianceInfo = server.SeverInfo.GetAllianceMgr().NewAllianceInfo(roleName, allianceName)
	user.SetAllianceName(allianceName)

	ret.Name = allianceName
	net.SendMsg(pbId.ScAllianceCreateId, ret, c)
	return
}

func CsAllianceList(data interface{}, c net.Conner) {
	_, ok := data.(*pb.CsAllianceList)
	if !ok {
		return
	}

	allianceList := server.SeverInfo.GetAllianceMgr().GetAllAlliance()
	ret := &pb.ScAllianceList{AllianceList: allianceList}
	net.SendMsg(pbId.ScAllianceListId, ret, c)
	return
}

func CsAllianceJoin(data interface{}, c net.Conner) {
	req, ok := data.(*pb.CsAllianceJoin)
	if !ok {
		return
	}
	allianceName := req.AllianceName
	roleName := c.GetRoleName()
	user, ok := server.SeverInfo.GetUser(roleName)
	ret := &pb.ScAllianceJoin{Result: 1}
	if !ok || user.Role.AllianceName != "" {
		ret.Result = 2
		ret.AllianceName = user.Role.AllianceName
		net.SendMsg(pbId.ScAllianceJoinId, ret, c)
		return
	}

	allianceInfo := server.SeverInfo.GetAllianceMgr().GetAlliance(allianceName)
	if allianceInfo == nil {
		ret.Result = 3
		ret.AllianceName = allianceName
		net.SendMsg(pbId.ScAllianceJoinId, ret, c)
		return
	}

	allianceInfo.AddMember(roleName)
	user.SetAllianceName(allianceName)

	ret.AllianceName = allianceName
	net.SendMsg(pbId.ScAllianceJoinId, ret, c)
	return
}

func CsAllianceDismiss(data interface{}, c net.Conner) {
	_, ok := data.(*pb.CsAllianceDismiss)
	if !ok {
		return
	}

	roleName := c.GetRoleName()
	user, ok := server.SeverInfo.GetUser(roleName)
	ret := &pb.ScAllianceDismiss{Result: 1}
	if !ok || user.Role.AllianceName == "" {
		ret.Result = 2
		net.SendMsg(pbId.ScAllianceDismissId, ret, c)
		return
	}

	allianceName := user.Role.AllianceName
	allianceInfo := server.SeverInfo.GetAllianceMgr().GetAlliance(allianceName)
	if allianceInfo == nil {
		ret.Result = 3
		ret.AllianceName = allianceName
		net.SendMsg(pbId.ScAllianceDismissId, ret, c)
		return
	}

	if allianceInfo.ChairName != c.GetRoleName() {
		ret.Result = 4
		ret.AllianceName = allianceName
		net.SendMsg(pbId.ScAllianceDismissId, ret, c)
		return
	}

	server.SeverInfo.GetAllianceMgr().DismissAlliance(allianceName)
	for memberName := range allianceInfo.Members {
		if memberUser, ok := server.SeverInfo.GetUser(memberName); ok {
			memberUser.Role.AllianceName = ""
			server.SeverInfo.StoreUser(memberUser)
		}
	}
	ret.AllianceName = allianceName
	net.SendMsg(pbId.ScAllianceDismissId, ret, c)
	return
}

func CsAllianceIncreaseCapacity(data interface{}, c net.Conner) {
	_, ok := data.(*pb.CsAllianceIncreaseCapacity)
	if !ok {
		return
	}

	roleName := c.GetRoleName()
	user, ok := server.SeverInfo.GetUser(roleName)
	ret := &pb.ScAllianceIncreaseCapacity{Result: 1}
	if !ok || user.Role.AllianceName == "" {
		ret.Result = 2
		net.SendMsg(pbId.ScAllianceIncreaseCapacityId, ret, c)
		return
	}

	allianceName := user.Role.AllianceName
	allianceInfo := server.SeverInfo.GetAllianceMgr().GetAlliance(allianceName)
	if allianceInfo == nil {
		ret.Result = 3
		net.SendMsg(pbId.ScAllianceIncreaseCapacityId, ret, c)
		return
	}

	if allianceInfo.ChairName != c.GetRoleName() {
		ret.Result = 4
		net.SendMsg(pbId.ScAllianceIncreaseCapacityId, ret, c)
		return
	}

	if allianceInfo.CapacityTimes >= 1 {
		ret.Result = 5
		net.SendMsg(pbId.ScAllianceIncreaseCapacityId, ret, c)
		return
	}

	allianceInfo.AddCapacity()
	net.SendMsg(pbId.ScAllianceIncreaseCapacityId, ret, c)
	return
}

func CsAllianceStoreItem(data interface{}, c net.Conner) {
	req, ok := data.(*pb.CsAllianceStoreItem)
	if !ok {
		return
	}

	roleName := c.GetRoleName()
	user, ok := server.SeverInfo.GetUser(roleName)
	ret := &pb.ScAllianceStoreItem{Result: 1}
	if !ok || user.Role.AllianceName == "" {
		ret.Result = 2
		net.SendMsg(pbId.ScAllianceStoreItemId, ret, c)
		return
	}

	allianceName := user.Role.AllianceName
	allianceInfo := server.SeverInfo.GetAllianceMgr().GetAlliance(allianceName)
	if allianceInfo == nil {
		ret.Result = 3
		net.SendMsg(pbId.ScAllianceStoreItemId, ret, c)
		return
	}

	itemId, itemNumber, index := req.Id, req.Number, req.Index
	err := allianceInfo.CheckStoreItem(itemId, itemNumber, index)
	if err != nil {
		ret.Result = 4
		net.SendMsg(pbId.ScAllianceStoreItemId, ret, c)
		return
	}

	net.SendMsg(pbId.ScAllianceStoreItemId, ret, c)
	return
}

func CsAllianceDestroyItem(data interface{}, c net.Conner) {
	req, ok := data.(*pb.CsAllianceDestroyItem)
	if !ok {
		return
	}

	roleName := c.GetRoleName()
	user, ok := server.SeverInfo.GetUser(roleName)
	ret := &pb.ScAllianceDestroyItem{Result: 1}
	if !ok || user.Role.AllianceName == "" {
		ret.Result = 2
		net.SendMsg(pbId.ScAllianceDestroyItemId, ret, c)
		return
	}

	allianceName := user.Role.AllianceName
	allianceInfo := server.SeverInfo.GetAllianceMgr().GetAlliance(allianceName)
	if allianceInfo == nil {
		ret.Result = 3
		net.SendMsg(pbId.ScAllianceDestroyItemId, ret, c)
		return
	}

	if allianceInfo.ChairName != c.GetRoleName() {
		ret.Result = 4
		net.SendMsg(pbId.ScAllianceDestroyItemId, ret, c)
		return
	}

	allianceInfo.DestroyItem(req.Index)
	net.SendMsg(pbId.ScAllianceDestroyItemId, ret, c)
	return
}

func CsAllianceClearup(data interface{}, c net.Conner) {
	_, ok := data.(*pb.CsAllianceClearup)
	if !ok {
		return
	}

	roleName := c.GetRoleName()
	user, ok := server.SeverInfo.GetUser(roleName)
	ret := &pb.ScAllianceClearup{Result: 1}
	if !ok || user.Role.AllianceName == "" {
		ret.Result = 2
		net.SendMsg(pbId.ScAllianceClearupId, ret, c)
		return
	}

	allianceName := user.Role.AllianceName
	allianceInfo := server.SeverInfo.GetAllianceMgr().GetAlliance(allianceName)
	if allianceInfo == nil {
		ret.Result = 3
		net.SendMsg(pbId.ScAllianceClearupId, ret, c)
		return
	}

	allianceInfo.Clearup()
	net.SendMsg(pbId.ScAllianceClearupId, ret, c)
	return
}
