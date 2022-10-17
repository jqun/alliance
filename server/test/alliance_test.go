package test

import (
	"alliance/proto-message/pb"
	"alliance/server/module/server"
	"github.com/stretchr/testify/assert"
	"testing"
)

/*
* 仓库列表为空
* 1号，2号格子会依次放入5个1类1号物品
 */
func TestAllianceFindNull(t *testing.T) {
	info := &server.AllianceInfo{
		CapacityTimes: 0,
		ItemList:      [40]*server.AllianceItem{},
	}

	// 此测试数据会将数据分别放入index=1和index=2的格子，并且放满
	idType, id := int32(1), int32(1)
	number := int32(10)
	index := int32(1)
	maxIndex := int32(30)
	ret := make(map[int32]int32)
	err := server.Find(idType, id, number, index, maxIndex, info.ItemList, ret)
	if err == nil {
		info.UpdateItem(idType, id, ret)
	}
	assert.Equal(t, int32(5), info.ItemList[0].TotalNum)
	assert.Equal(t, int32(5), info.ItemList[1].TotalNum)
}

/*
* 1号格子仓库已放入4个1类物品
* 1号格子继续放入1个1类1号物品，2号格子放入5个1类1号物品，3号格子放入4个1类1号物品
 */
func TestAllianceFindOne(t *testing.T) {
	var items []*pb.TestItem
	items = append(items, &pb.TestItem{
		Id:       1,
		ItemType: 1,
		Number:   4,
	})
	info := &server.AllianceInfo{
		CapacityTimes: 0,
		ItemList: [40]*server.AllianceItem{
			0: {ItemType: 1, TotalNum: 4, ItemArray: pb.TestItem_Array{Items: items}},
		},
	}

	idType, id := int32(1), int32(1)
	number := int32(10)
	index := int32(1)
	maxIndex := int32(30)
	ret := make(map[int32]int32)
	err := server.Find(idType, id, number, index, maxIndex, info.ItemList, ret)
	if err == nil {
		info.UpdateItem(idType, id, ret)
	}
	assert.Equal(t, int32(5), info.ItemList[0].TotalNum)
	assert.Equal(t, int32(5), info.ItemList[1].TotalNum)
	assert.Equal(t, int32(4), info.ItemList[2].TotalNum)
}

/*
* 1号格子仓库已放入4个1类1号物品
* 1号格子不能放入2类1号物品，2号格子放入5个2类1号物品，3号格子放入5个2类1号物品
 */
func TestAllianceFindDiffType(t *testing.T) {
	var items []*pb.TestItem
	items = append(items, &pb.TestItem{
		Id:       1,
		ItemType: 1,
		Number:   4,
	})
	info := &server.AllianceInfo{
		CapacityTimes: 0,
		ItemList: [40]*server.AllianceItem{
			0: {ItemType: 1, TotalNum: 4, ItemArray: pb.TestItem_Array{Items: items}},
		},
	}

	idType, id := int32(2), int32(1)
	number := int32(10)
	index := int32(1)
	maxIndex := int32(30)
	ret := make(map[int32]int32)
	err := server.Find(idType, id, number, index, maxIndex, info.ItemList, ret)
	if err == nil {
		info.UpdateItem(idType, id, ret)
	}
	assert.Equal(t, int32(5), info.ItemList[1].TotalNum)
	assert.Equal(t, int32(5), info.ItemList[2].TotalNum)
}

/*
* 1号格子，放入1类物品4个，1类物品1个
* 整理后：
* 1号格子，放入1类物品5个
 */
func TestAllianceClearup(t *testing.T) {
	var itemsIndex1 []*pb.TestItem
	itemsIndex1 = append(itemsIndex1, &pb.TestItem{Id: 1, ItemType: 1, Number: 4})
	itemsIndex1 = append(itemsIndex1, &pb.TestItem{Id: 1, ItemType: 1, Number: 1})
	info := &server.AllianceInfo{
		CapacityTimes: 0,
		ItemList: [40]*server.AllianceItem{
			0: {ItemType: 1, TotalNum: 4, ItemArray: pb.TestItem_Array{Items: itemsIndex1}},
		},
	}

	info.Clearup()
	assert.Equal(t, int32(5), info.ItemList[0].TotalNum)
}

/*
* 1号格子，放入4个1类1号物品
* 2号格子，放入1个1类1号物品
* 3号格子，放入1个2类1号物品
* 整理后：
* 1号格子放入5个1号1类物品
* 2号格子放入1个2类1号物品
 */
func TestAllianceClearup2(t *testing.T) {
	var itemsIndex1 []*pb.TestItem
	itemsIndex1 = append(itemsIndex1, &pb.TestItem{Id: 1, ItemType: 1, Number: 4})

	var itemsIndex2 []*pb.TestItem
	itemsIndex2 = append(itemsIndex2, &pb.TestItem{Id: 1, ItemType: 1, Number: 1})

	var itemsIndex3 []*pb.TestItem
	itemsIndex3 = append(itemsIndex3, &pb.TestItem{Id: 1, ItemType: 2, Number: 1})

	info := &server.AllianceInfo{
		CapacityTimes: 0,
		ItemList: [40]*server.AllianceItem{
			0: {ItemType: 1, TotalNum: 4, ItemArray: pb.TestItem_Array{Items: itemsIndex1}},
			1: {ItemType: 1, TotalNum: 1, ItemArray: pb.TestItem_Array{Items: itemsIndex2}},
			2: {ItemType: 2, TotalNum: 1, ItemArray: pb.TestItem_Array{Items: itemsIndex3}},
		},
	}

	info.Clearup()
	assert.Equal(t, int32(5), info.ItemList[0].TotalNum)
	assert.Equal(t, int32(1), info.ItemList[1].TotalNum)
}

/*
* 1号格子，放入1个2类1号物品
* 2号格子，放入4个1类1号物品
* 3号格子，放入1个1类1号物品
* 整理后：
* 1号格子放入5个1号1类物品
* 2号格子放入1个2类1号物品
 */
func TestAllianceClearup3(t *testing.T) {
	var itemsIndex1 []*pb.TestItem
	itemsIndex1 = append(itemsIndex1, &pb.TestItem{Id: 1, ItemType: 2, Number: 1})

	var itemsIndex2 []*pb.TestItem
	itemsIndex2 = append(itemsIndex2, &pb.TestItem{Id: 1, ItemType: 1, Number: 4})

	var itemsIndex3 []*pb.TestItem
	itemsIndex3 = append(itemsIndex3, &pb.TestItem{Id: 1, ItemType: 1, Number: 1})

	info := &server.AllianceInfo{
		CapacityTimes: 0,
		ItemList: [40]*server.AllianceItem{
			0: {ItemType: 1, TotalNum: 4, ItemArray: pb.TestItem_Array{Items: itemsIndex1}},
			1: {ItemType: 1, TotalNum: 1, ItemArray: pb.TestItem_Array{Items: itemsIndex2}},
			2: {ItemType: 2, TotalNum: 1, ItemArray: pb.TestItem_Array{Items: itemsIndex3}},
		},
	}

	info.Clearup()
	assert.Equal(t, int32(5), info.ItemList[0].TotalNum)
	assert.Equal(t, int32(1), info.ItemList[1].TotalNum)
}