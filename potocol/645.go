package potocol

import (
	"container/list"
	"fmt"
	"mg-coms/command"
	"mg-coms/config"
	"mg-coms/device"
	"sort"
	"strconv"
	"strings"
	"time"
)

//Potocol104 104规约
type Potocol645 struct {
	PotocolParm
	Pm      potocol645Map    //分区
	Keys    []int            //分区号表
	Msgs    *list.List       //消息队列
	Message command.Message  //消息
	St      *Potocol645State //状态
	Ds      *Device645State
}

type potocol645Map map[int]*Potocol645Partition

//Potocol645Partition 规约分区
type Potocol645Partition struct {
	name         string    //分区名称
	terminalNum  int       //终端号
	partitionNum int       //分区号
	Addr         []byte    //分区地址
	funCode      byte      //控制码
	byteNum      int       //字节个数
	point        []data645 //点表
	timeout      string    //分区的超时时间
}

//data 分区内点的配置
type data645 struct {
	id    int     //点号
	name  string  //点名
	kind  int     //点的类型
	scale float32 //系数
	addr  []byte  //地址
	bit   byte    //比特位
	ptype int32   //数据类型
	on    []byte  //分闸
	off   []byte  //合闸
}

//State 状态
type Potocol645State struct {
	cursor           int       //当前下发的分区下标
	maxCursor        int       //游标最大值
	timeout          string    //超时时间
	done             bool      //发送或接收完成
	currentPartition int       //分区号
	frame            []byte    //内容存储
	rxTime           time.Time //接收时间
	txTime           time.Time //发送时间

}

type Device645State struct {
	ds map[int]*OnOffline645
}

type OnOffline645 struct {
	ID     int
	State  bool
	Number int
}

func (m *Device645State) init(c config.ChannelConfig) *Device645State {
	m.ds = make(map[int]*OnOffline645)
	for _, dev := range c.Device {
		m.ds[dev.ID] = new(OnOffline645).init(dev)
	}
	return m
}

func (m *Device645State) clearTimeoutNum(index int) {
	m.ds[index].Number = 0
}

func (m *Device645State) reset() {
	for _, ds := range m.ds {
		ds.Number = 0
		ds.State = true
	}
}

func (m *Device645State) isRun(index int) bool {
	return m.ds[index].State
}

func (m *Device645State) addTimeoutNum(nm *Potocol645) {
	index := nm.St.getDeviceNum()
	m.ds[index].Number++
	if m.ds[index].Number > 9 {
		m.ds[index].State = false
		nm.DeviceMap[nm.ChannelConfig.ID*10000+index].Status = false
	}
}

func (l *OnOffline645) init(d config.DeviceConfig) *OnOffline645 {
	l.ID = d.ID
	l.State = d.Status
	l.Number = 0
	return l
}

//Init 645
func (p *Potocol645) Init(cc config.ChannelConfig, dm device.DeviceMap) *Potocol645 {
	p.ChannelConfig = cc
	p.DeviceMap = dm
	return p
}

func (p *Potocol645) Cmd(msg command.Message) {

}

func (m *Potocol645) Created() {
	m.Msgs = list.New()
	m.Message = command.Message{}
	m.Pm = make(map[int]*Potocol645Partition)
	m.Ds = new(Device645State).init(m.ChannelConfig)
	m.initPartitionTable(m.ChannelConfig)
	m.partitionSort()
	m.St = new(Potocol645State).init(len(m.Keys))
	go m.timer()
}

func (s *Potocol645State) init(maxCursor int) *Potocol645State {
	s.frame = make([]byte, 0, 1024)
	s.maxCursor = maxCursor
	s.cursor = 0
	s.timeout = "500ms"
	return s
}

func (s *Potocol645State) reset() {
	s.done = false
	s.currentPartition = 0
	s.frame = nil
	s.rxTime = time.Time{}
	s.txTime = time.Time{}
}

//timer 定时器
func (m *Potocol645) timer() {
	//每30min重置设备在线状态
	for range time.Tick(5 * time.Minute) {
		m.Ds.reset()
		for index, _ := range m.DeviceMap {
			m.DeviceMap[index].Status = true
		}
	}
}

//initPartitionTable 格式化分区
func (m *Potocol645) initPartitionTable(c config.ChannelConfig) {
	for _, dev := range c.Device {
		m.partitionConfig(dev.Ycs, dev, 0)
		m.partitionConfig(dev.Yxs, dev, 1)
		m.partitionConfig(dev.Yms, dev, 2)
	}
}

//partitionConfig 分区配置
func (m *Potocol645) partitionConfig(pts []config.Point, dev config.DeviceConfig, kind int) {
	fmt.Println("ppppppp", pts)
	for _, p := range pts {
		d := data645{}
		d.kind = kind
		d.name = p.Name
		d.scale = p.Scale
		d.id = p.ID
		adr := strings.Split(p.Data.(map[string]interface{})["Addr"].(string), ",")
		for index, _ := range adr {
			d.addr = append(d.addr, intTobyte(int(getUintData("0x" + adr[index])))[3])
		}
		pbit := p.GetIntData("Bit")
		d.bit = byte(pbit)
		d.ptype = int32(p.GetIntData("Type"))

		pm := new(Potocol645Partition)
		pm.terminalNum = dev.ID
		pm.name = dev.Name
		pm.timeout = "500ms"
		staddr := strings.Split(p.Data.(map[string]interface{})["Addr"].(string), ",")
		for index, _ := range staddr {
			pm.Addr = append(pm.Addr, intTobyte(int(getUintData("0x" + staddr[index])))[3])
		}
		pm.partitionNum = p.GetIntData("Partition")
		fc := p.GetIntData("Funcode")
		pm.funCode = intTobyte(int(getUintData("0x" + strconv.Itoa(fc))))[3]
		pm.byteNum = 4
		pm.point = append(pm.point, d)
		m.Pm[pm.terminalNum+pm.partitionNum] = pm
		sort.Stable(m.Pm[pm.terminalNum+pm.partitionNum])
	}
}

func (s Potocol645Partition) Len() int      { return len(s.point) }
func (s Potocol645Partition) Swap(i, j int) { s.point[i], s.point[j] = s.point[j], s.point[i] }
func (s Potocol645Partition) Less(i, j int) bool {
	return bytetoint(s.point[i].addr) < bytetoint(s.point[j].addr)
}

//partitionSort  对分区进行排序获得键值
func (m *Potocol645) partitionSort() {
	var keys []int
	// fmt.Println("获得键值时", m.Pm)
	for k := range m.Pm {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	m.Keys = keys
}

func (s *Potocol645State) isDone() bool {
	return s.done
}

func (s *Potocol645State) getCursor() int {
	return s.cursor
}

func getUintData(s string) uint64 {
	i, _ := strconv.ParseUint(s, 0, 0)
	return i
}

func (s *Potocol645State) setFrame(b []byte) {
	s.frame = b
}

func (s *Potocol645State) getFrameData(index int) byte {
	return s.frame[index]
}

func (s *Potocol645State) addCursor() {
	s.cursor++
	if s.cursor == s.maxCursor {
		s.cursor = 0
	}
}

func (s *Potocol645State) setDone(b bool) {
	s.done = b
}

func (s *Potocol645State) updateTxTime() {
	s.txTime = time.Now()
}

func (s *Potocol645State) updateRxTime() {
	s.rxTime = time.Now()
}

func (s *Potocol645State) getTimeout() string {
	trsub := s.rxTime.Sub(s.txTime)
	ad, _ := time.ParseDuration("350ms")
	newt := trsub + ad
	s.timeout = newt.String()
	return s.timeout
	// return "500ms"
}

func (s *Potocol645State) getCurrentPartition() int {
	return s.currentPartition
}

func (s *Potocol645State) isTimeout() bool {
	s.updateRxTime()
	t, _ := time.ParseDuration(s.timeout)
	return s.rxTime.After(s.txTime.Add(t))
}

func (s *Potocol645State) getDeviceNum() int {
	var num []int
	for i := 10; i > 0; i-- {
		s := strconv.FormatInt(int64(s.frame[i]), 16)
		tmp, _ := strconv.Atoi(s)
		num = append(num, tmp)
	}
	return num[0]*100000*100000 + num[1]*100000000 + num[2]*1000000 + num[3]*10000 + num[4]*100 + num[5]
}

func (p *Potocol645) Do() {
	p.tx()
	p.rx()
	time.Sleep(time.Millisecond * 50)
}

func (m *Potocol645) rx() {
	if !m.St.isDone() {
		return
	}
	buf, err := m.Port.Read()
	if err != nil {
		return
	}
	for _, b := range buf {
		m.Buffer = append(m.Buffer, b)
	}
	fmt.Println("rxrxrxrxrxrxrxrxrxrxrxrx", m.Buffer)
	if m.St.isTimeout() {
		// fmt.Println("1s", m.St.getDeviceNum())
		if m.Ds.isRun(m.St.getDeviceNum()) {
			m.Ds.addTimeoutNum(m)
		}
		// m.St.setDone(false)
		m.Buffer = nil
		m.St.reset()
		fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
		return
	}

	if len(m.Buffer) < 24 {
		fmt.Println("shujubugouchang   &&&&&&&")
		return
	}

	found := false
	index := 0
	for index, _ = range m.Buffer {
		if m.Buffer[index+4] == m.St.getFrameData(4) && m.Buffer[index+5] == m.St.getFrameData(5) && m.Buffer[index+6] == m.St.getFrameData(6) && m.Buffer[index+7] == m.St.getFrameData(7) && m.Buffer[index+8] == m.St.getFrameData(8) && m.Buffer[index+9] == m.St.getFrameData(9) && m.Buffer[index+10] == m.St.getFrameData(10) {
			if m.Buffer[index+14] == m.St.getFrameData(14) && m.Buffer[index+15] == m.St.getFrameData(15) && m.Buffer[index+16] == m.St.getFrameData(16) && m.Buffer[index+17] == m.St.getFrameData(17) {
				found = true
				break
			}
		} else {
			m.Buffer = m.Buffer[index+1:]
			break
		}
	}
	//没有完整帧返回
	if !found {
		// m.St.reset()
		fmt.Println(" weizhaodaowanzhengzhen   $$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$")
		return
	}

	var md *Potocol645Partition

	//找到该帧对应分区
	for _, mp := range m.Pm {
		if m.St.getCurrentPartition() == (mp.terminalNum + mp.partitionNum) {
			md = mp
			break
		}
	}

	//没找到分区返回
	if md == nil {
		fmt.Println("************************************************")
		m.St.reset()
		return
	}
	for _, pt := range md.point {
		var sdata []string
		for i := 21; i > 17; i-- {
			s := strconv.FormatInt(int64(m.Buffer[i]-byte(51)), 16)
			sdata = append(sdata, s)
		}
		str0 := sdata[0]
		str1 := sdata[1]
		str2 := sdata[2]
		str3 := sdata[3]
		data0, _ := strconv.Atoi(str0)
		data0 = data0 * 1000
		data1, _ := strconv.Atoi(str1)
		data1 = data1 * 100
		data2, _ := strconv.Atoi(str2)
		data2 = data2 * 1
		datatmp := data0 + data1 + data2
		data3, _ := strconv.Atoi(str3)
		// data := float32(float32(datatmp)+float32(data3)/float32(math.Pow10(len(sdata[3])))) * pt.scale
		data := float32(float32(datatmp)+float32(data3)/float32(100)) * pt.scale //小数点后值保留两位 除以100
		m.DeviceMap[m.ChannelConfig.ID*10000+md.terminalNum].SetYc(pt.id, data)
	}
	if m.Ds.isRun(m.St.getDeviceNum()) {
		m.Ds.clearTimeoutNum(m.St.getDeviceNum())
	}
	m.Buffer = m.Buffer[index+24:]

	md.timeout = m.St.getTimeout()
	m.St.reset()
	return
}

func (m *Potocol645) tx() {
	if m.St.isDone() {
		return
	}
	var frame = make([]byte, 0, 1024)
	keys := m.Keys
	cur := m.St.getCursor()
	if len(m.Pm) == 0 {
		return
	}
	pm := m.Pm[keys[cur]]
	m.St.timeout = pm.timeout
	m.St.currentPartition = keys[cur]
	for i := 0; i < 4; i++ {
		frame = append(frame, byte(254))
	}
	frame = append(frame, byte(104))
	sternum := strconv.Itoa(pm.terminalNum)
	s := strings.Split(sternum, "")
	var str []string
	for i := len(s) - 1; i > 0; {
		str = append(str, s[i-1]+s[i])
		i -= 2
	}
	if len(str) == 5 {
		str = append(str, "00")
	}
	for index, _ := range str {
		frame = append(frame, intTobyte(int(getUintData("0x" + str[index])))[3])
	}
	frame = append(frame, byte(104))
	frame = append(frame, pm.funCode)
	frame = append(frame, byte(pm.byteNum))
	for i := 3; i >= 0; i-- {
		frame = append(frame, pm.Addr[i])
	}
	var check byte
	for i := 4; i < len(frame); i++ {
		check += frame[i]
	}
	frame = append(frame, check)
	frame = append(frame, byte(22))
	m.St.setFrame(frame)
	if m.Ds.isRun(pm.terminalNum) == false {
		m.St.addCursor()
		return
	}
	fmt.Println(frame)
	err := m.Port.Write(string(frame), len(frame))
	if err != nil {
		fmt.Println(err)
	}
	m.St.setDone(true)
	m.St.updateTxTime()
	m.St.addCursor()
}
