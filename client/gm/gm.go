package gm

import "fmt"

const (
	// 玩家登录登出命令
	GmLogin  = "login"  // 登录
	GmLogout = "logout" // 登出

	// 公会操作命令
	GmWhichAlliance    = "whichAlliance"    // 查询所属公会
	GmCreateAlliance   = "createAlliance"   // 创建公会
	GmAllianceList     = "allianceList"     // 查询公会列表
	GmJoinAlliance     = "joinAlliance"     // 没有公会加入公会
	GmDismissAlliance  = "dismissAlliance"  // 解散公会
	GmIncreaseCapacity = "increaseAlliance" // 公会扩容

	// 仓库操作
	GmStoreItem    = "storeItem"   // 成员存储物品
	GmDestroyItem  = "destroyItem" // 会长销毁某一个仓库格子
	GmCleanup      = "clearup"     // 整理仓库
	GmGetGmExplain = "help"        // 获取gm说明
)

var GmExplain = []string{
	fmt.Sprintf("GM[%v] usage:login server,e.g.:%v [roleName]", GmLogin, GmLogin),
	fmt.Sprintf("GM[%v] usage:logout server,e.g.:%v", GmLogout, GmLogout),
	fmt.Sprintf("GM[%v] usage:check which alliance,e.g.:%v", GmWhichAlliance, GmWhichAlliance),
	fmt.Sprintf("GM[%v] usage:create alliance,e.g.:%v [allianceName]", GmCreateAlliance, GmCreateAlliance),
	fmt.Sprintf("GM[%v] usage:list all alliance,e.g.:%v", GmAllianceList, GmAllianceList),
	fmt.Sprintf("GM[%v] usage:join alliance,e.g.:%v [allianceName]", GmJoinAlliance, GmJoinAlliance),
	fmt.Sprintf("GM[%v] usage:dismiss the alliance,e.g.:%v", GmDismissAlliance, GmDismissAlliance),
	fmt.Sprintf("GM[%v] usage:chair increate allicance capacity,e.g.:%v", GmIncreaseCapacity, GmIncreaseCapacity),
	fmt.Sprintf("GM[%v] usage:store item,e.g.:%v [itemId] [itemNum] [index]", GmStoreItem, GmStoreItem),
	fmt.Sprintf("GM[%v] usage:destory alliance,e.g.:%v [index]", GmDestroyItem, GmDestroyItem),
	fmt.Sprintf("GM[%v] usage:clean up the alliance,e.g.:%v", GmCleanup, GmCleanup),
	fmt.Sprintf("GM[%v] usage:print all gm,e.g.:%v", GmGetGmExplain, GmGetGmExplain),
}
