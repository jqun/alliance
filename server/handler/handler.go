package handler

import (
	"alliance/common/net"
	"alliance/proto-message/pb"
	pbId "alliance/proto-message/pb-id"
)
func InitHandler() {
	net.RegisterHandler(pbId.CsAccountLoginId, &pb.CsAccountLogin{}, CsAccountLogin)
	net.RegisterHandler(pbId.CsAccountLogoutId, &pb.CsAccountLogout{}, CsAccountLogout)

	net.RegisterHandler(pbId.CsAllianceInfoId, &pb.CsAllianceInfo{}, CsAllianceInfo)
	net.RegisterHandler(pbId.CsAllianceCreateId, &pb.CsAllianceCreate{}, CsAllianceCreate)
	net.RegisterHandler(pbId.CsAllianceListId, &pb.CsAllianceList{}, CsAllianceList)
	net.RegisterHandler(pbId.CsAllianceJoinId, &pb.CsAllianceJoin{}, CsAllianceJoin)
	net.RegisterHandler(pbId.CsAllianceDismissId, &pb.CsAllianceDismiss{}, CsAllianceDismiss)

	net.RegisterHandler(pbId.CsAllianceIncreaseCapacityId, &pb.CsAllianceIncreaseCapacity{}, CsAllianceIncreaseCapacity)
	net.RegisterHandler(pbId.CsAllianceStoreItemId, &pb.CsAllianceStoreItem{}, CsAllianceStoreItem)
	net.RegisterHandler(pbId.CsAllianceDestroyItemId, &pb.CsAllianceDestroyItem{}, CsAllianceDestroyItem)
	net.RegisterHandler(pbId.CsAllianceClearupId, &pb.CsAllianceClearup{}, CsAllianceClearup)
}