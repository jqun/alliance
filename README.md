# alliance
简易仓库管理系统

## 目录树 ##
```
alliance
├─client
│  ├─config   
│  ├─c_net    
│  ├─gm  
│  ├─handler   
│  ├─_client    
│  └─_config       
├─common
│  ├─consts   
│  ├─message 
│  ├─net    
│  └─util        
├─proto
│  └─script        
├─proto-message
│  ├─pb   
│  └─pb-id    
└─server
    ├─config 
    ├─handler   
    ├─module
    │  ├─role 
    │  └─server    
    ├─test
    └─_config

```
## client ##
使用命令行模式的简易客户端
提供如下命令：
> 登录：login RoleName
> 登出：logout
> 查看公会信息：whichAlliance
> 创建公会：createAlliance allianceName
> 会长解散公会：dismissAlliance
> 公会列表：allianceList
> 加入公会：joinAlliance allianceName
> 会长扩容仓库：increaseCapacity
> 公会成员提交物品：storeItem itemId itemNum index
> 会长销毁仓库：destroyItem index
> 仓库整理：clearup

### server ###


1. 将CS协议注册到handlerMap
2. 启动服务并监听客户端连接请求
3. 监听退出信号
	```
	func (m *serverInfo) Run() {
		go m.signalListen()
		m.NetSever.Run()
	}

	func (m *Server) Run() {
		for {
			conn, err := m.Listener.Accept()
			if err != nil {
				if m.ListenerClose {
					break
				} else {
					log.Printf("accept socket err,error is [%v]", err.Error())
					continue
				}
			}
			log.Printf("new a addr[%v] conn", conn.RemoteAddr().String())
			conner := NewConner(conn)
			conner.bySeverCreate = true
			conner.setReadDeadLine()
			go conner.Start()
		}
	}
	```
4. 用户关键结构设计
	```
	type Role struct {
		Name         string
		AllianceName string // 公会名字
	}

	type User struct {
		Role *Role
	}
	```
5. 公会关键结构设计
	```
	type AllianceInfo struct {
		sync.RWMutex
		CapacityTimes int32                       // 扩容次数
		Members       map[string]struct{}         // 公会成员
		ChairName     string                      // 会长名字
		ItemList      [MaxItemIndex]*AllianceItem // 设定最大长度为40的数组
		close         chan struct{}
	}
	```
6. 仓库放置主逻辑
    ```
    func Find(idType, id, number, index, maxIndex int32, itemList [MaxItemIndex]*AllianceItem, ret map[int32]int32) error {
    	if index > maxIndex || index < 1 { // 检查一轮后，仍然未放置完成，则放置失败
    		return errors.New("the index overflow")
    	}
    
    	itemM := itemList[index-1]
    	if itemM == nil {
    		itemM = &AllianceItem{ItemType: idType}
    	}
    
    	indexNext, newMaxIndex := nextIndex(index, maxIndex)
    	if itemM.ItemType != idType { // 类型不同则一次放入下一个格子
    		return Find(idType, id, number, indexNext, newMaxIndex, itemList, ret)
    	}
    
    	totalNum := itemM.TotalNum
    	if totalNum >= MaxItemNum { // 此格子放置满
    		// 寻找下一个可以放的格子，如果都不能放，则存储失败
    		return Find(idType, id, number, indexNext, newMaxIndex, itemList, ret)
    	} else if totalNum+number > MaxItemNum { // 加起来超过限制,需往下一组格子存储
    		n := MaxItemNum - totalNum
    		ret[index] = n
    		return Find(idType, id, number-n, indexNext, newMaxIndex, itemList, ret)
    	} else { // 放置完成
    		ret[index] = number
    		return nil
    	}
    }
    ```

## 部署 ##
