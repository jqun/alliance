package handler

import (
	"alliance/common/net"
	"alliance/proto-message/pb"
	pbId "alliance/proto-message/pb-id"
)

func InitHandler() {
	net.RegisterHandler(pbId.ScAccountLoginId, &pb.ScAccountLogin{}, ScAccountLogin)
	net.RegisterHandler(pbId.ScAccountLogoutId, &pb.ScAccountLogout{}, ScAccountLogout)

	net.RegisterHandler(pbId.ScAllianceInfoId, &pb.ScAllianceInfo{}, ScAllianceInfo)
	net.RegisterHandler(pbId.ScAllianceCreateId, &pb.ScAllianceCreate{}, ScAllianceCreate)
	net.RegisterHandler(pbId.ScAllianceListId, &pb.ScAllianceList{}, ScAllianceList)
	net.RegisterHandler(pbId.ScAllianceJoinId, &pb.ScAllianceJoin{}, ScAllianceJoin)
	net.RegisterHandler(pbId.ScAllianceDismissId, &pb.ScAllianceDismiss{}, ScAllianceDismiss)

	net.RegisterHandler(pbId.ScAllianceIncreaseCapacityId, &pb.ScAllianceIncreaseCapacity{}, ScAllianceIncreaseCapacity)
	net.RegisterHandler(pbId.ScAllianceStoreItemId, &pb.ScAllianceStoreItem{}, ScAllianceStoreItem)
	net.RegisterHandler(pbId.ScAllianceDestroyItemId, &pb.ScAllianceDestroyItem{}, ScAllianceDestroyItem)
	net.RegisterHandler(pbId.ScAllianceClearupId, &pb.ScAllianceClearup{}, ScAllianceClearup)
}
