package server

import (
	"alliance/proto-message/pb"
	"errors"
	"sort"
	"sync"
)

type allianceMgr struct {
	alliances sync.Map
}

const (
	DefaultItemIndex  = 30 // 默认仓库容量
	IncreaseItemIndex = 10 // 仓库扩容数量
	MaxItemIndex      = DefaultItemIndex + IncreaseItemIndex
	MaxItemNum        = 5 // 物品叠加数量
)

type allianceInfo struct {
	sync.RWMutex
	CapacityTimes int32                       // 扩容次数
	Members       map[string]struct{}         // 公会成员
	ChairName     string                      // 会长名字
	ItemList      [MaxItemIndex]*allianceItem // 设定最大长度为40的数组
	close         chan struct{}
}

type allianceItem struct {
	itemType  int32             // 指明此类格子物品类型
	totalNum  int32             // 此格子放置物品总数量
	ItemArray pb.TestItem_Array // itemId->item
}

func (m *allianceMgr) NewAllianceInfo(chairName, allianceName string) *allianceInfo {
	info := &allianceInfo{
		Members:   map[string]struct{}{chairName: {}},
		ChairName: chairName,
		close:     make(chan struct{}),
	}
	m.alliances.Store(allianceName, info)
	return info
}

func (m *allianceMgr) GetAlliance(allianceName string) *allianceInfo {
	v, ok := m.alliances.Load(allianceName)
	if !ok {
		return nil
	}
	return v.(*allianceInfo)
}

func (m *allianceMgr) GetAllAlliance() (allianceList []string) {
	m.alliances.Range(func(key, _ interface{}) bool {
		allianceList = append(allianceList, key.(string))
		return true
	})
	return
}

func (m *allianceMgr) DismissAlliance(allianceName string) {
	m.alliances.Delete(allianceName)
}

func (m *allianceMgr) stop() {
	m.alliances.Range(func(key, value interface{}) bool {
		r := value.(*allianceInfo)
		r.closeHandle()
		return true
	})
}

func (m *allianceInfo) closeHandle() {
	close(m.close)
}

func (m *allianceInfo) AddMember(roleName string) {
	m.Lock()
	defer m.Unlock()

	m.Members[roleName] = struct{}{}
}

func (m *allianceInfo) AddCapacity() {
	m.Lock()
	defer m.Unlock()
	m.CapacityTimes = 1
}

func (m *allianceInfo) CheckStoreItem(id, number, index int32) error {
	m.Lock()
	defer m.Unlock()

	indexNum := m.CapacityTimes*IncreaseItemIndex + DefaultItemIndex
	if index > indexNum || index < 1 {
		return errors.New("the index overflow")
	}

	idType, ok := itemInfo[id]
	if !ok {
		return errors.New("item id error")
	}

	var changeMap = make(map[int32]int32) // index->num
	err := find(idType, id, number, index, indexNum, m.ItemList, changeMap)
	if err != nil {
		return err
	}

	m.updateItem(idType, id, changeMap)
	return nil
}

func (m *allianceInfo) updateItem(idType, id int32, changeMap map[int32]int32) {
	for idx, num := range changeMap {
		itemM := m.ItemList[idx-1]
		if itemM == nil {
			itemM = &allianceItem{itemType: idType}
			m.ItemList[idx-1] = itemM
		}

		itemM.ItemArray.Items = append(itemM.ItemArray.Items, &pb.TestItem{
			Id:       id,
			Name:     "",
			ItemType: idType,
			Number:   num,
		})
	}
}

func (m *allianceInfo) DestroyItem(index int32) {
	m.Lock()
	defer m.Unlock()

	indexNum := m.CapacityTimes*IncreaseItemIndex + DefaultItemIndex
	if index > indexNum || index < 1 {
		return
	}
	m.ItemList[index-1] = nil
}

func (m *allianceInfo) Clearup() {
	m.Lock()
	defer m.Unlock()

	typeMap := make(map[int32]map[int32]*pb.TestItem) // type->id->number
	var typeList []int
	for _, item := range m.ItemList {
		if item == nil {
			continue
		}

		nums, ok := typeMap[item.itemType]
		if !ok {
			typeList = append(typeList, int(item.itemType))
			nums = make(map[int32]*pb.TestItem)
			typeMap[item.itemType] = nums
		}

		for _, testItem := range item.ItemArray.Items {
			if im, iOk := nums[testItem.Id]; iOk {
				im.Number += testItem.Number
			} else {
				nums[testItem.Id] = &pb.TestItem{
					Id:       testItem.Id,
					Name:     "",
					ItemType: item.itemType,
					Number:   testItem.Number,
				}
			}
		}
	}

	m.ItemList = [MaxItemIndex]*allianceItem{} // 清空仓库重新放置
	sort.Ints(typeList) // 按类型排序
	var index = int32(1)
	indexNum := m.CapacityTimes*IncreaseItemIndex + DefaultItemIndex
	for _, tt := range typeList {
		nums := typeMap[int32(tt)]

		var items []*pb.TestItem
		for _, item := range nums {
			items = append(items, item)
		}
		sort.Slice(items, func(i, j int) bool {
			return items[i].Number > items[j].Number
		})

		for _, item := range items {
			var changeMap = make(map[int32]int32) // index->num
			_ = find(item.ItemType, item.Id, item.Number, index, indexNum, m.ItemList, changeMap)
			m.updateItem(item.ItemType, item.Id, changeMap)
		}
	}
}

// 仓库物品放置主逻辑
func find(idType, id, number, index, maxIndex int32, itemList [MaxItemIndex]*allianceItem, ret map[int32]int32) error {
	if index > maxIndex || index < 1 { // 检查一轮后，仍然未放置完成，则放置失败
		return errors.New("the index overflow")
	}

	itemM := itemList[index-1]
	if itemM == nil {
		itemM = &allianceItem{itemType: idType}
	}

	indexNext, newMaxIndex := nextIndex(index, maxIndex)
	if itemM.itemType != idType { // 类型不同则一次放入下一个格子
		return find(idType, id, number, indexNext, newMaxIndex, itemList, ret)
	}

	totalNum := itemM.totalNum
	if totalNum >= MaxItemNum { // 此格子放置满
		// 寻找下一个可以放的格子，如果都不能放，则存储失败
		return find(idType, id, number, indexNext, newMaxIndex, itemList, ret)
	} else if totalNum+number > MaxItemNum { // 加起来超过限制
		n := MaxItemNum - totalNum
		ret[index] = n
		return find(idType, id, number-n, indexNext, newMaxIndex, itemList, ret)
	} else { // 放置完成
		ret[index] = number
		return nil
	}
}

func nextIndex(index, maxIndex int32) (int32, int32) {
	if index+1 <= maxIndex {
		return index + 1, maxIndex
	}
	return 1, index - 1
}

var itemInfo = make(map[int32]int32) // id->idType
