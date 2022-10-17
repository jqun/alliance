package _client

import (
	"alliance/client/c_net"
	"alliance/client/gm"
	"alliance/common/consts"
	network "alliance/common/net"
	"alliance/common/util"
	"alliance/proto-message/pb"
	pbId "alliance/proto-message/pb-id"
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type clientInfo struct {
	RoleName string
	Conner   network.Conner
	Gm       chan string
	Close    chan struct{}
}

var Client *clientInfo

func NewClient(addr string) *clientInfo {
	m := &clientInfo{}
	m.Gm = make(chan string, 10)
	m.Close = make(chan struct{})
	conner := c_net.NewClientConner(addr)
	m.Conner = conner
	m.flagsParse()
	return m
}

func (m *clientInfo) flagsParse() {
	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	fs.StringVar(&m.RoleName, "roleName", "", "set _client role name")
	_ = fs.Parse(os.Args[1:])
}

func (m *clientInfo) Run() {
	log.Printf("_client start running")
	go m.gmStart()
	go m.ScanStdin()
	go m.signalHandle()
	m.Conner.Start()
}

func (m *clientInfo) gmStart() {
	heartTicker := time.NewTicker(consts.HeartBeatTime)
	for {
		select {
		case gmStr := <-m.Gm:
			m.handleGmId(gmStr)
		case <-heartTicker.C:
			m.sendClientSocketHeart()
		case <-m.Close:
			break
		}
	}
}

func (m *clientInfo) signalHandle() {
	ticker := time.NewTicker(time.Minute * 2)
	sNotify := make(chan os.Signal)
	signal.Notify(sNotify, os.Interrupt, os.Kill, syscall.SIGTERM)
	for {
		select {
		case <-sNotify:
			m.logout()
			m.Stop()
			return
		case <-ticker.C:
			log.Printf("signal listen running")
		}
	}
}

func (m *clientInfo) Stop() {
	m.Conner.Stop()
	close(m.Close)
}

func (m *clientInfo) handleGmId(str string) {
	gmArr := util.SplitSpace(str)
	if len(gmArr) == 0 {
		return
	}
	switch gmArr[0] {
	case gm.GmLogin:
		m.login(gmArr)
	case gm.GmLogout:
		m.logout()
	case gm.GmWhichAlliance:
		m.whichAlliance()
	case gm.GmCreateAlliance:
		m.createAlliance(gmArr)
	case gm.GmAllianceList:
		m.allianceList()
	case gm.GmJoinAlliance:
		m.joinAlliance(gmArr)
	case gm.GmDismissAlliance:
		m.dismissAlliance()
	case gm.GmIncreaseCapacity:
		m.increaseCapacity()
	case gm.GmStoreItem:
		m.storeItem(gmArr)
	case gm.GmDestroyItem:
		m.destroyItem(gmArr)
	case gm.GmCleanup:
		m.clearup()
	case gm.GmGetGmExplain:
		m.printGmExplain()
	default:
		log.Printf("gm[%v] don't handle", gmArr[0])
	}
}

func (m *clientInfo) receiveGmStr(str string) {
	m.Gm <- str
}

func (m *clientInfo) ScanStdin() {
	scanner := bufio.NewScanner(os.Stdin)
	log.Printf("wait stdin!")
	m.printGmExplain()
	for scanner.Scan() {
		m.receiveGmStr(scanner.Text())
	}
}

func (m *clientInfo) printGmExplain() {
	for _, gmStr := range gm.GmExplain {
		fmt.Println(gmStr)
	}
}

//登录
func (m *clientInfo) login(gmArr []string) {
	if len(gmArr) < 2 {
		log.Printf("login err,please input RoleName")
		return
	}
	req := &pb.CsAccountLogin{}
	req.RoleName = gmArr[1]
	m.RoleName = req.RoleName
	network.SendMsg(pbId.CsAccountLoginId, req, m.Conner)
}

//登出
func (m *clientInfo) logout() {
	req := &pb.CsAccountLogout{}
	network.SendMsg(pbId.CsAccountLogoutId, req, m.Conner)
}

func (m *clientInfo) sendClientSocketHeart() {
	log.Printf("account[%v] send heart msg", m.RoleName)
	network.SendMsg(pbId.ClientSocketHeartId, &pb.ClientHeartBeat{}, m.Conner)
}

// 公会管理
func (m *clientInfo) whichAlliance() {
	req := &pb.CsAllianceInfo{}
	network.SendMsg(pbId.CsAllianceInfoId, req, m.Conner)
}

func (m *clientInfo) createAlliance(gmArr []string) {
	if len(gmArr) < 2 {
		log.Printf("createAlliance err,please input AllianceName")
		return
	}
	if m.RoleName == "" || gmArr[1] == "" {
		log.Printf("createAlliance err, please login first")
		return
	}
	req := &pb.CsAllianceCreate{Name: gmArr[1]}
	network.SendMsg(pbId.CsAllianceCreateId, req, m.Conner)
}

func (m *clientInfo) allianceList() {
	req := &pb.CsAllianceList{}
	network.SendMsg(pbId.CsAllianceListId, req, m.Conner)
}

func (m *clientInfo) joinAlliance(gmArr []string) {
	if len(gmArr) < 2 {
		log.Printf("joinAlliance err,please input AllianceName")
		return
	}
	if m.RoleName == "" || gmArr[1] == "" {
		log.Printf("joinAlliance err, please login first")
		return
	}
	req := &pb.CsAllianceJoin{AllianceName: gmArr[1]}
	network.SendMsg(pbId.CsAllianceJoinId, req, m.Conner)
}

func (m *clientInfo) dismissAlliance() {
	if m.RoleName == "" {
		log.Printf("joinAlliance err, please login first")
		return
	}
	req := &pb.CsAllianceDismiss{}
	network.SendMsg(pbId.CsAllianceDismissId, req, m.Conner)
}

func (m *clientInfo) increaseCapacity() {
	if m.RoleName == "" {
		log.Printf("increaseCapacity err, please login first")
		return
	}
	req := &pb.CsAllianceDismiss{}
	network.SendMsg(pbId.CsAllianceDismissId, req, m.Conner)
}

func (m *clientInfo) storeItem(gmArr []string) {
	if m.RoleName == "" {
		log.Printf("storeItem err, please login first")
		return
	}

	if len(gmArr) < 4 {
		log.Printf("storeItem err,please check your inputs")
		return
	}

	id, number, index := util.StrToInt32(gmArr[1]), util.StrToInt32(gmArr[2]), util.StrToInt32(gmArr[3])
	if id == 0 || number == 0 || index == 0 {
		log.Printf("storeItem err,please check your inputs")
		return
	}

	req := &pb.CsAllianceStoreItem{Id: id, Number: number, Index: index}
	network.SendMsg(pbId.CsAllianceStoreItemId, req, m.Conner)
}

func (m *clientInfo) destroyItem(gmArr []string) {
	if m.RoleName == "" {
		log.Printf("destroyItem err, please login first")
		return
	}

	if len(gmArr) < 2 {
		log.Printf("destroyItem err, please check your index")
		return
	}

	req := &pb.CsAllianceDestroyItem{Index: util.StrToInt32(gmArr[1])}
	network.SendMsg(pbId.CsAllianceDestroyItemId, req, m.Conner)
}

func (m *clientInfo) clearup() {
	if m.RoleName == "" {
		log.Printf("clearup err, please login first")
		return
	}

	req := &pb.CsAllianceClearup{}
	network.SendMsg(pbId.CsAllianceClearupId, req, m.Conner)
}
