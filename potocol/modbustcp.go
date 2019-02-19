package potocol

import (
	"container/list"
	"fmt"
	"mg-coms/command"
	"mg-coms/config"
	"mg-coms/device"
	"sort"
	"time"
)

type PotocolModbusTcp struct {
	PotocolParm
	Pm      ModbusTcpMap    //分区
	Keys    []int           //分区号表
	Msgs    *list.List      //消息队列
	Message command.Message //消息
	St      *ModbusTcpState //状态
	Ds      *TcpDeviceState
}

//ModbusTcpPartition 规约分区
type ModbusTcpPartition struct {
	name         string          //分区名称
	terminalNum  int             //终端号
	partitionNum int             //分区号
	startAddr    []byte          //起始地址
	endAddr      []byte          //终止地址
	funCode      byte            //功能码
	registerNum  []byte          //寄存器个数
	point        []modbustcpdata //点表
	timeout      string          //分区的超时时间
}

//data 分区内点的配置
type modbustcpdata struct {
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

type TcpdevOnOffline struct {
	ID     int
	State  bool
	Number int
}

//ModbusTcpMap 分区列表
type ModbusTcpMap map[int]*ModbusTcpPartition

//TcpDeviceState 设备状态
type TcpDeviceState struct {
	ds map[int]*TcpdevOnOffline
}

func (m *TcpDeviceState) init(c config.ChannelConfig) *TcpDeviceState {
	m.ds = make(map[int]*TcpdevOnOffline)
	for _, dev := range c.Device {
		m.ds[dev.ID] = new(TcpdevOnOffline).init(dev)
	}
	return m
}

//reset 重置设备状态
func (m *TcpDeviceState) reset() {
	for _, ds := range m.ds {
		ds.Number = 0
		ds.State = true
	}
}

//isRun 设备运行状态
func (m *TcpDeviceState) isRun(index int) bool {
	return m.ds[index].State
}

//getTimeoutNum 获得设备超时次数
func (m *TcpDeviceState) getTimeoutNum(index int) int {
	return m.ds[index].Number
}

//addTimeoutNum 设备超时累加超时次数
func (m *TcpDeviceState) addTimeoutNum(nm *PotocolModbusTcp) {
	index := nm.St.getDeviceNum()
	m.ds[index].Number++
	if m.ds[index].Number > 9 {
		m.ds[index].State = false
		nm.DeviceMap[nm.ChannelConfig.ID*10000+index].Status = false
	}
}

//clearTimeoutNum 清空设备超时次数
func (m *TcpDeviceState) clearTimeoutNum(index int) {
	m.ds[index].Number = 0
}

func (l *TcpdevOnOffline) init(d config.DeviceConfig) *TcpdevOnOffline {
	l.ID = d.ID
	l.State = d.Status
	l.Number = 0
	return l
}

//ModbusTcpState 状态
type ModbusTcpState struct {
	cursor           int    //当前下发的分区下标
	maxCursor        int    //游标最大值
	timeout          string //超时时间
	done             bool   //发送或接收完成
	msgcontrol       command.Message
	selecttime       time.Time
	cmdtime          time.Time
	currentPartition int       //分区号
	cmdState         bool      //是否是控制命令
	frame            []byte    //内容存储
	rxTime           time.Time //接收时间
	txTime           time.Time //发送时间

}

func (s *ModbusTcpState) init(maxCursor int) *ModbusTcpState {
	s.frame = make([]byte, 0, 1024)
	s.maxCursor = maxCursor
	s.cursor = 0
	s.timeout = "500ms"
	return s
}

//reset
func (s *ModbusTcpState) reset() {
	s.done = false
	s.currentPartition = 0
	s.cmdState = false
	s.frame = nil
	// s.timeout = "500ms"
	s.rxTime = time.Time{}
	s.txTime = time.Time{}
	s.msgcontrol = command.Message{}
	s.selecttime = time.Time{}
	s.cmdtime = time.Time{}
}

func (s *ModbusTcpState) isDone() bool {
	return s.done
}

func (s *ModbusTcpState) updateRxTime() {
	s.rxTime = time.Now()
}

func (s *ModbusTcpState) updateTxTime() {
	s.txTime = time.Now()
}

func (s *ModbusTcpState) setDone(b bool) {
	s.done = b
}
func (s *ModbusTcpState) isCmd() bool {
	return s.cmdState
}
func (s *ModbusTcpState) setCmd(b bool) {
	s.cmdState = b
}

func (s *ModbusTcpState) setFrame(b []byte) {
	s.frame = b
}

func (s *ModbusTcpState) getDeviceNum() int {
	return int(s.frame[6])
}

func (s *ModbusTcpState) getCursor() int {
	return s.cursor
}

func (s *ModbusTcpState) addCursor() {
	s.cursor++
	if s.cursor == s.maxCursor {
		s.cursor = 0
	}
}

func (s *ModbusTcpState) getFrameData(index int) byte {
	return s.frame[index]
}

func (s *ModbusTcpState) getCurrentPartition() int {
	return s.currentPartition
}

func (s *ModbusTcpState) isTimeout() bool {
	s.updateRxTime()
	t, _ := time.ParseDuration(s.timeout)
	// fmt.Println("tttttttttttttttttttimeout", t)
	// fmt.Println("istimeout=rx-tx", s.rxTime.Sub(s.txTime).String())
	return s.rxTime.After(s.txTime.Add(t))
}

func (s *ModbusTcpState) getTimeout() string {
	trsub := s.rxTime.Sub(s.txTime)
	// fmt.Println("rxrxrxrxrx", s.rxTime)
	// fmt.Println("txtxtxtxtx", s.txTime)
	// fmt.Println("gettimeout=rx-tx", trsub.String())
	ad, _ := time.ParseDuration("350ms")
	// if trsub < ad {
	// 	s.timeout = "500ms"
	// 	return s.timeout
	// }
	newt := trsub + ad
	s.timeout = newt.String()
	return s.timeout
	// return "500ms"
}

//Created 初始化规约缓存
func (m *PotocolModbusTcp) Created() {
	m.Msgs = list.New()
	// m.Msgs = nil
	m.Message = command.Message{}
	m.Pm = make(map[int]*ModbusTcpPartition)
	m.Ds = new(TcpDeviceState).init(m.ChannelConfig)
	m.initPartitionTable(m.ChannelConfig)
	m.partitionSort()
	m.St = new(ModbusTcpState).init(len(m.Keys))
	go m.timer()
}

//timer 定时器
func (m *PotocolModbusTcp) timer() {
	//每30min重置设备在线状态
	for range time.Tick(5 * time.Minute) {
		m.Ds.reset()
		for index, _ := range m.DeviceMap {
			m.DeviceMap[index].Status = true
		}
	}
}

//initPartitionTable 格式化分区
func (m *PotocolModbusTcp) initPartitionTable(c config.ChannelConfig) {
	for _, dev := range c.Device {
		m.partitionConfig(dev.Ycs, dev, 0)
		m.partitionConfig(dev.Yxs, dev, 1)
		m.partitionConfig(dev.Yms, dev, 2)
	}
}

//partitionConfig 分区配置
func (m *PotocolModbusTcp) partitionConfig(pts []config.Point, dev config.DeviceConfig, kind int) {
	for _, p := range pts {
		d := modbustcpdata{}
		d.kind = kind
		d.name = p.Name
		d.scale = p.Scale
		d.id = p.ID
		adr := p.GetIntData("Addr")
		for index, ad := range intTobyte(adr) {
			if index > 1 {
				d.addr = append(d.addr, ad)
			}
		}
		pbit := p.GetIntData("Bit")
		d.bit = byte(pbit)
		d.ptype = int32(p.GetIntData("Type"))
		ptmp := p.GetIntData("Partition")
		var pnum = dev.ID*100 + ptmp
		_, ok := m.Pm[pnum]
		if ok {
			m.updatePartition(p, d, pnum)
		} else {
			m.newPartition(p, dev, d)
		}
	}
}

func (s ModbusTcpPartition) Len() int      { return len(s.point) }
func (s ModbusTcpPartition) Swap(i, j int) { s.point[i], s.point[j] = s.point[j], s.point[i] }
func (s ModbusTcpPartition) Less(i, j int) bool {
	return bytetoint(s.point[i].addr) < bytetoint(s.point[j].addr)
}

//updatePartition 更新分区
func (m *PotocolModbusTcp) updatePartition(p config.Point, d modbustcpdata, pnum int) {
	var have = false
	adr := p.GetIntData("Addr")
	if int(byteToint(m.Pm[pnum].startAddr)) > adr {
		ad := intTobyte(adr)
		m.Pm[pnum].startAddr = ad[2:4]
	}
	for _, dt := range m.Pm[pnum].point {
		if adr == int(byteToint(dt.addr)) {
			have = true
			break
		}
	}
	if have == false {
		m.Pm[pnum].endAddr = intTobyte(int(byteToint32(m.Pm[pnum].endAddr)) + registerNum(int(d.ptype)))
		m.Pm[pnum].registerNum = m.Pm[pnum].endAddr[2:4]
	}
	m.Pm[pnum].point = append(m.Pm[pnum].point, d)
	sort.Stable(m.Pm[pnum])
}

//newPartition 创建分区
func (m *PotocolModbusTcp) newPartition(p config.Point, dev config.DeviceConfig, d modbustcpdata) {
	pm := new(ModbusTcpPartition)
	pm.terminalNum = dev.ID
	pm.name = dev.Name
	pm.timeout = "500ms"
	staddr := p.GetIntData("Addr")
	for index, st := range intTobyte(staddr) {
		if index > 1 {
			pm.startAddr = append(pm.startAddr, st)
		}
	}
	pm.partitionNum = p.GetIntData("Partition")
	fc := p.GetIntData("Funcode")
	pm.funCode = byte(fc)
	pm.endAddr = intTobyte(registerNum(int(d.ptype)))
	pm.registerNum = pm.endAddr[2:4]
	pm.point = append(pm.point, d)
	m.Pm[pm.terminalNum*100+pm.partitionNum] = pm
}

//partitionSort  对分区进行排序获得键值
func (m *PotocolModbusTcp) partitionSort() {
	var keys []int
	for k := range m.Pm {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	m.Keys = keys
}

//Cmd消息处理
func (m *PotocolModbusTcp) Cmd(msg command.Message) {
	fmt.Println("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^", msg)
	m.Msgs.PushBack(msg)
}

//Do 轮询收发
func (m *PotocolModbusTcp) Do() {
	m.tx()
	m.rx()
	time.Sleep(time.Millisecond * 50)
}

//rx 接收
func (m *PotocolModbusTcp) rx() {

	//未完成返回
	if !m.St.isDone() {
		return
	}

	//超时返回

	buf, err := m.Port.Read()

	if err != nil {
		fmt.Println("err", err)
		return
	}
	// fmt.Println("buf====")
	// fmt.Printf("% x", buf)

	for _, b := range buf {
		m.Buffer = append(m.Buffer, b)
	}
	// fmt.Println(time.Now().Sub(m.St.txTime).String())
	////不足最短帧数据量返回

	//处理数据
	if !m.St.isCmd() {
		// fmt.Println("BUFFER===", m.Buffer)
		// m.St.updateRxTime()
		if m.St.isTimeout() {
			// fmt.Println("1s", m.St.getDeviceNum())
			if m.Ds.isRun(m.St.getDeviceNum()) {
				m.Ds.addTimeoutNum(m)
			}
			// m.St.setDone(false)
			m.Buffer = nil
			m.St.reset()
			// m.St.timeout = "500ms"
			fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
			return
		}
		if len(m.Buffer) < 11 {
			fmt.Println("&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&")
			return
		}

		//寻找完整一帧数据
		found := false
		index := 0
		for index, _ = range m.Buffer {
			//寻找首尾
			if m.Buffer[index+6] == m.St.getFrameData(6) && m.Buffer[index+7] == m.St.getFrameData(7) {
				if len(m.Buffer)-6 < int(m.Buffer[5]) {
					break
				}
				found = true
				break
			} else {
				m.Buffer = m.Buffer[index+1:]
				break
			}
		}

		//没有完整帧返回
		if !found {
			// m.St.reset()
			fmt.Println("$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$")
			return
		}

		var md *ModbusTcpPartition

		//找到该帧对应分区
		for _, mp := range m.Pm {
			if m.St.getCurrentPartition() == (mp.terminalNum*100 + mp.partitionNum) {
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

		//开始处理帧数据
		mark := int(byteToint(md.startAddr))
		a := index + 9
		var b int
		for _, dt := range md.point {
			if int(byteToint(dt.addr)) == mark {
				b = a + byteNumTcp(dt)
				if b >= (len(m.Buffer) + 1) {
					m.Buffer = nil
					if m.Ds.isRun(m.St.getDeviceNum()) {
						m.Ds.clearTimeoutNum(m.St.getDeviceNum())
					}
					m.St.reset()
					return
				}
				data := m.calculate(m.Buffer[a:b], int(dt.ptype), dt.scale)
				switch dt.kind {
				case 0:
					m.DeviceMap[m.ChannelConfig.ID*10000+md.terminalNum].SetYc(dt.id, data)
					break
				case 1:
					str := bytetoBinaryString(m.Buffer[a]) + bytetoBinaryString(m.Buffer[a+1])
					var strs []string
					for i := 15; i >= 0; i-- {
						strs = append(strs, string(str[i]))
					}
					for index, bt := range strs {
						for _, d := range md.point {
							if int(byteToint(d.addr)) == mark && int(d.bit) == index {
								if string(bt) == "1" {
									m.DeviceMap[m.ChannelConfig.ID*10000+md.terminalNum].SetYx(d.id, true)
								} else {
									m.DeviceMap[m.ChannelConfig.ID*10000+md.terminalNum].SetYx(d.id, false)
								}
							}
						}
					}
					fmt.Println("tcptcp##########", strs)
					break
				case 2:
					m.DeviceMap[m.ChannelConfig.ID*10000+md.terminalNum].SetYm(dt.id, data)
					break
				}
			} else {
				continue
			}
			mark = mark + byteNumTcp(dt)/2
			a = b
			// if pindex < len(md.point)-1 {
			// 	if bytetoint(dt.addr) == bytetoint(md.point[pindex+1].addr) {
			// 		b = a
			// 	} else {
			// 		a = b
			// 	}
			// }
		}
		if m.Ds.isRun(m.St.getDeviceNum()) {
			m.Ds.clearTimeoutNum(m.St.getDeviceNum())
		}
		m.Buffer = m.Buffer[index+int(m.Buffer[index+8])+9:]

		md.timeout = m.St.getTimeout()
		// fmt.Println("outoutoutoutout------", md.timeout)
		m.St.reset()

		//帧数据处理完成返回
		return
	}

	//处理命令
	// fmt.Println("___________________________________", time.Now(), m.St.selecttime)

	if m.St.msgcontrol.SubCmd == 4 {
		ad, _ := time.ParseDuration("5s")
		if time.Now().After(m.St.selecttime.Add(ad)) {
			m.St.msgcontrol = command.Message{}
			// m.Msgs.Remove(m1)
			m.Message = command.Message{}
			m.St.reset()
			return
		}
		m.St.setDone(false)
		return
	}

	if len(m.Buffer) < 8 {
		tcmd, _ := time.ParseDuration("5s")
		if time.Now().After(m.St.cmdtime.Add(tcmd)) {
			m.St.reset()
			m.Buffer = nil
			return
		}
		return
	}

	// m.St.setCmd(false)
	// m.St.setDone(false)

	fmt.Println("BBBBBBBBBBBBBBBBBBBBBBBB", m.Buffer)
	m.Buffer = nil
	m.St.reset()
	// for index, _ := range m.Buffer {
	// 	if m.Buffer[index] == m.St.getFrameData(0) && m.Buffer[index+1] == m.St.getFrameData(1) && m.Buffer[index+2] == m.St.getFrameData(2) {
	// 		crc := make([]byte, 0, 10)
	//
	// 		v := intTobyte(int(crc16(m.Buffer[index:index+6], 6)))
	// 		crc = append(crc, v[3])
	// 		crc = append(crc, v[2])
	// 		recrc := m.Buffer[index+6 : index+8]
	// 		if crc[0] != recrc[0] && crc[1] != recrc[1] {
	// 			m.Buffer = m.Buffer[index+1:]
	// 			break
	// 		}
	// 		// msg := command.Message{
	// 		// 	Cmd:     m.Message.Cmd,
	// 		// 	Value:   m.Message.Value,
	// 		// 	SubCmd:  m.Message.SubCmd,
	// 		// 	Device:  m.Message.Device,
	// 		// 	Point:   m.Message.Point,
	// 		// 	Success: true,
	// 		// }
	// 		// *m.Bus <- msg
	// 		// m.Message = command.Message{}
	// 		m.St.reset()
	// 		m.Buffer = nil
	// 		return
	// 	}
	// }
}

//tx 发送
func (m *PotocolModbusTcp) tx() {
	if m.St.isDone() {
		return
	}

	var frame = make([]byte, 0, 1024)

	//有消息优先处理消息
	if m.Msgs.Len() > 0 && !m.St.isDone() {

		m.St.setCmd(true)
		fmt.Println("Msgs-+-+-+", m.Msgs)
		m1 := m.Msgs.Front()
		m.Message = m1.Value.(command.Message)
		if m.St.msgcontrol.SubCmd != 4 {
			if m.Message.SubCmd == 4 {
				m.St.msgcontrol = m.Message
				m.St.selecttime = time.Now()
				m.Message = command.Message{}
			}
			m.Msgs.Remove(m1)
			m.St.setDone(true)
			return
		}
		fmt.Println("m.MESSSSSSS", m.Message)
		switch m.Message.SubCmd {
		case 4:
			m.St.msgcontrol = m.Message
			m.St.selecttime = time.Now()
			m.Message = command.Message{}
			m.Msgs.Remove(m1)
			break
		case 5:
			fmt.Println("msgcon", m.St.msgcontrol)
			fmt.Println("MESS", m.Message)
			if m.St.msgcontrol.Device == m.Message.Device && m.St.msgcontrol.Point == m.Message.Point {
				switch m.Message.Cmd {
				case 10:
					fmt.Println("YK")
					for _, dev := range m.ChannelConfig.Device {
						if dev.ID == m.Message.Device {
							for _, p := range dev.Yks {
								if p.ID == m.Message.Point {
									for a := 0; a < 6; a++ {
										frame = append(frame, 0)
									}
									frame = append(frame, byte(dev.ID))
									fc := p.GetIntData("Funcode")
									frame = append(frame, byte(fc))
									paddr := p.GetIntData("Addr")
									for index, st := range intTobyte(paddr) {
										if index > 1 {
											frame = append(frame, st)
										}
									}
									if m.Message.Value.(bool) == true {
										n := p.GetUintData("On")
										for index, on := range intTobyte(int(n)) {
											if index > 1 {
												frame = append(frame, on)
											}
										}
									} else {
										f := p.GetUintData("Off")
										for index, off := range intTobyte(int(f)) {
											if index > 1 {
												frame = append(frame, off)
											}
										}
									}
									frame[5] = intTobyte(len(frame) - 6)[3]
									err := m.Port.Write(string(frame), len(frame))
									if err != nil {
										fmt.Println(err)
									}
									m.St.cmdtime = time.Now()
									fmt.Println("yk success")
									break
								}
							}
							break
						}
					}
					m.St.setFrame(frame)
					m.St.msgcontrol = command.Message{}
					m.Msgs.Remove(m1)
					break
				case 11: //yt消息处理
					fmt.Println("YT")
					for _, dev := range m.ChannelConfig.Device {
						if dev.ID == m.Message.Device {
							for _, p := range dev.Yts {
								if p.ID == m.Message.Point {
									for a := 0; a < 6; a++ {
										frame = append(frame, 0)
									}
									frame = append(frame, byte(dev.ID))
									fc := p.GetIntData("Funcode")
									frame = append(frame, byte(fc))
									paddr := p.GetIntData("Addr")
									for index, st := range intTobyte(paddr) {
										if index > 1 {
											frame = append(frame, st)
										}
									}
									scale := p.Scale
									ptype := p.GetIntData("Type")
									if fc == 16 {
										if 6 < ptype && ptype < 19 {
											for index, rt := range intTobyte(registerNum(ptype)) {
												if index > 1 {
													frame = append(frame, rt)
												}
											}
											frame = append(frame, byte(4))
										}
										if ptype == 19 {
											for index, rt := range intTobyte(registerNum(ptype)) {
												if index > 1 {
													frame = append(frame, rt)
												}
											}
											frame = append(frame, byte(8))
										}
									}
									for _, f := range m.byteSort(ptype, scale) {
										frame = append(frame, f)
									}
									frame[5] = intTobyte(len(frame) - 6)[3]
									err := m.Port.Write(string(frame), len(frame))
									if err != nil {
										fmt.Println(err)
									}
									m.St.cmdtime = time.Now()
									fmt.Println("yt success")
									break
								}
							}
							break
						}
					}
					break
				}
			}
			m.St.setFrame(frame)
			m.St.msgcontrol = command.Message{}
			m.Msgs.Remove(m1)
			break
		case 6:
			if m.St.msgcontrol.Device == m.Message.Device && m.St.msgcontrol.Point == m.Message.Point {
				m.St.msgcontrol = command.Message{}
				m.Message = command.Message{}
			}
			m.Msgs.Remove(m1)
			break
		}
		m.St.setDone(true)
		return

		//--------------------------------------------------------------------------------------------------
	}

	//未接收完成返回

	if m.St.cmdState == false {
		//没有消息类型则开始处理数据
		keys := m.Keys
		cur := m.St.getCursor()
		pm := m.Pm[keys[cur]]
		m.St.timeout = pm.timeout //------

		fmt.Println("分区", keys[cur])
		m.St.currentPartition = keys[cur]
		for a := 0; a < 6; a++ {
			frame = append(frame, 0)
		}
		frame = append(frame, byte(pm.terminalNum))
		frame = append(frame, pm.funCode)
		for _, staddr := range pm.startAddr {
			frame = append(frame, staddr)
		}
		for _, rtnum := range pm.registerNum {
			frame = append(frame, rtnum)
		}
		m.St.setFrame(frame)
		if m.Ds.isRun(pm.terminalNum) == false {
			m.St.addCursor()
			return
		}
		frame[5] = intTobyte(len(frame) - 6)[3]
		fmt.Println("frame-----")
		fmt.Printf("% x", frame)
		err := m.Port.Write(string(frame), len(frame))
		if err != nil {
			fmt.Println(err)
		}
		m.St.setDone(true)
		m.St.updateTxTime()
		// }
		m.St.addCursor()
	} else {
		m.St.setDone(true)
	}
}

//byteSort 字节排序
func (m *PotocolModbusTcp) byteSort(ptype int, scale float32) []byte {
	var frame = make([]byte, 0, 1024)
	switch ptype {
	case 0:
		msg := int(m.Message.Value.(float64))
		msg = msg * int(scale)
		for index, v := range intTobyte(msg) {
			if index > 1 {
				frame = append(frame, v)
			}
		}
		break
	case 1:
		msg := int(m.Message.Value.(float64))
		msg = msg * int(scale)
		for index, v := range intTobyte(msg) {
			if index > 1 {
				frame = append(frame, v)
			}
		}
		break
	case 2:
		msg := int(m.Message.Value.(float64))
		msg = msg * int(scale)
		v := intTobyte(msg)
		frame = append(frame, v[3])
		frame = append(frame, v[2])
		break
	case 3:
		msg := int(m.Message.Value.(float64))
		msg = msg * int(scale)
		for index, v := range intTobyte(msg) {
			if index > 1 {
				frame = append(frame, v)
			}
		}
		break
	case 4:
		msg := int(m.Message.Value.(float64))
		msg = msg * int(scale)
		v := intTobyte(msg)
		frame = append(frame, v[3])
		frame = append(frame, v[2])
		break
	case 5:
		msg := int(m.Message.Value.(float64))
		msg = msg * int(scale)
		for index, v := range intTobyte(msg) {
			if index > 1 {
				frame = append(frame, v)
			}
		}
		break
	case 6:
		msg := int(m.Message.Value.(float64))
		msg = msg * int(scale)
		v := intTobyte(msg)
		frame = append(frame, v[3])
		frame = append(frame, v[2])
		break
	case 7:
		msg := int(m.Message.Value.(float64))
		msg = msg * int(scale)
		v := intTobyte(msg)
		frame = append(frame, v[0])
		frame = append(frame, v[1])
		frame = append(frame, v[2])
		frame = append(frame, v[3])
		break
	case 8:
		msg := int(m.Message.Value.(float64))
		msg = msg * int(scale)
		v := intTobyte(msg)
		frame = append(frame, v[1])
		frame = append(frame, v[0])
		frame = append(frame, v[3])
		frame = append(frame, v[2])
		break
	case 9:
		msg := int(m.Message.Value.(float64))
		msg = msg * int(scale)
		v := intTobyte(msg)
		frame = append(frame, v[2])
		frame = append(frame, v[3])
		frame = append(frame, v[0])
		frame = append(frame, v[1])
		break
	case 10:
		msg := int(m.Message.Value.(float64))
		msg = msg * int(scale)
		v := intTobyte(msg)
		frame = append(frame, v[3])
		frame = append(frame, v[2])
		frame = append(frame, v[1])
		frame = append(frame, v[0])
		break
	case 11:
		msg := int(m.Message.Value.(float64))
		msg = msg * int(scale)
		v := intTobyte(msg)
		frame = append(frame, v[0])
		frame = append(frame, v[1])
		frame = append(frame, v[2])
		frame = append(frame, v[3])
		break
	case 12:
		msg := int(m.Message.Value.(float64))
		msg = msg * int(scale)
		v := intTobyte(msg)
		frame = append(frame, v[1])
		frame = append(frame, v[0])
		frame = append(frame, v[3])
		frame = append(frame, v[2])
		fmt.Printf("%x", frame)
		break
	case 13:
		msg := int(m.Message.Value.(float64))
		msg = msg * int(scale)
		v := intTobyte(msg)
		frame = append(frame, v[2])
		frame = append(frame, v[3])
		frame = append(frame, v[0])
		frame = append(frame, v[1])
		break
	case 14:
		msg := int(m.Message.Value.(float64))
		msg = msg * int(scale)
		v := intTobyte(msg)
		frame = append(frame, v[3])
		frame = append(frame, v[2])
		frame = append(frame, v[1])
		frame = append(frame, v[0])
		break
	case 15:
		msg := float32(m.Message.Value.(float64))
		msg = msg * float32(scale)
		v := float32Tobyte(msg)
		fmt.Println(v)
		frame = append(frame, v[0])
		frame = append(frame, v[1])
		frame = append(frame, v[2])
		frame = append(frame, v[3])
		break
	case 16:
		msg := float32(m.Message.Value.(float64))
		msg = msg * float32(scale)
		v := float32Tobyte(msg)
		frame = append(frame, v[1])
		frame = append(frame, v[0])
		frame = append(frame, v[3])
		frame = append(frame, v[2])
		break
	case 17:
		msg := float32(m.Message.Value.(float64))
		msg = msg * float32(scale)
		v := float32Tobyte(msg)
		frame = append(frame, v[2])
		frame = append(frame, v[3])
		frame = append(frame, v[0])
		frame = append(frame, v[1])
		break
	case 18:
		msg := float32(m.Message.Value.(float64))
		msg = msg * float32(scale)
		v := float32Tobyte(msg)
		frame = append(frame, v[3])
		frame = append(frame, v[2])
		frame = append(frame, v[1])
		frame = append(frame, v[0])
		break
	default:
		msg := int(m.Message.Value.(float64))
		msg = msg * int(scale)
		v := int64Tobyte(msg)
		fmt.Println(v, len(v))
		frame = append(frame, v[0])
		frame = append(frame, v[1])
		frame = append(frame, v[2])
		frame = append(frame, v[3])
		frame = append(frame, v[4])
		frame = append(frame, v[5])
		frame = append(frame, v[6])
		frame = append(frame, v[7])
		break
	}
	return frame
}

//calculate 计算接收到的数据
func (m *PotocolModbusTcp) calculate(buf []byte, ptype int, scale float32) float32 {
	// fmt.Println("ptype", ptype)
	var tmpa uint16
	var tmpb uint16
	var tmpc uint32
	var tmpd uint32
	var tmpe = make([]byte, 10)
	var tmpf uint64
	var value float32
	switch ptype {
	case 0:
	case 1:
		tmpa = uint16(buf[0])
		tmpa = (tmpa << 8) & 0xff00
		tmpa += uint16(buf[1])
		tt := int16(tmpa)
		value = float32(tt) * scale
		return value
	case 2:
		tmpa = uint16(buf[1])
		tmpa = (tmpa << 8) & 0xff00
		tmpa += uint16(buf[0])
		tt := int16(tmpa)
		value = float32(tt) * scale
		return value
	case 3:
		tmpa = uint16(buf[0])
		tmpa = (tmpa << 8) & 0xff00
		tmpa += uint16(buf[1])
		value = float32(int(tmpa)) * scale
		return value
	case 4:
		tmpa = uint16(buf[1])
		tmpa = (tmpa << 8) & 0xff00
		tmpa += uint16(buf[0])
		value = float32(int(tmpa)) * scale
		return value
	case 5:
		tmpb = uint16(buf[0])
		tmpb = (tmpb << 8) & 0xff00
		tmpb += uint16(buf[1])
		value = float32(tmpb) * scale
		return value
	case 6:
		tmpb = uint16(buf[1])
		tmpb = (tmpb << 8) & 0xff00
		tmpb += uint16(buf[0])
		value = float32(tmpb) * scale
		return value
	case 8:
		tmpc = uint32(buf[1])
		tmpc = (tmpc << 8) & 0xff00
		tmpc += uint32(buf[0])
		tmpc = (tmpc << 8) & 0x00ffff00
		tmpc += uint32(buf[3])
		tmpc = (tmpc << 8) & 0xffffff00
		tmpc += uint32(buf[2])
		tt := int32(tmpc)
		value = float32(tt) * scale
		return value
	case 7:
		tmpc = uint32(buf[0])
		tmpc = (tmpc << 8) & 0xff00
		tmpc += uint32(buf[1])
		tmpc = (tmpc << 8) & 0x00ffff00
		tmpc += uint32(buf[2])
		tmpc = (tmpc << 8) & 0xffffff00
		tmpc += uint32(buf[3])
		tt := int32(tmpc)
		value = float32(tt) * scale
		return value
	case 9:
		tmpc = uint32(buf[2])
		tmpc = (tmpc << 8) & 0xff00
		tmpc += uint32(buf[3])
		tmpc = (tmpc << 8) & 0x00ffff00
		tmpc += uint32(buf[0])
		tmpc = (tmpc << 8) & 0xffffff00
		tmpc += uint32(buf[1])
		tt := int32(tmpc)
		value = float32(tt) * scale
		return value
	case 10:
		tmpc = uint32(buf[3])
		tmpc = (tmpc << 8) & 0xff00
		tmpc += uint32(buf[2])
		tmpc = (tmpc << 8) & 0x00ffff00
		tmpc += uint32(buf[1])
		tmpc = (tmpc << 8) & 0xffffff00
		tmpc += uint32(buf[0])
		tt := int32(tmpc)
		value = float32(tt) * scale
		return value
	case 11:
		tmpd = uint32(buf[0])
		tmpd = (tmpd << 8) & 0xff00
		tmpd += uint32(buf[1])
		tmpd = (tmpd << 8) & 0x00ffff00
		tmpd += uint32(buf[2])
		tmpd = (tmpd << 8) & 0xffffff00
		tmpd += uint32(buf[3])
		value = float32(tmpd) * scale
		return value
	case 12:
		tmpd = uint32(buf[1])
		tmpd = (tmpd << 8) & 0xff00
		tmpd += uint32(buf[0])
		tmpd = (tmpd << 8) & 0x00ffff00
		tmpd += uint32(buf[3])
		tmpd = (tmpd << 8) & 0xffffff00
		tmpd += uint32(buf[2])
		value = float32(tmpd) * scale
		return value
	case 13:
		tmpd = uint32(buf[2])
		tmpd = (tmpd << 8) & 0xff00
		tmpd += uint32(buf[3])
		tmpd = (tmpd << 8) & 0x00ffff00
		tmpd += uint32(buf[0])
		tmpd = (tmpd << 8) & 0xffffff00
		tmpd += uint32(buf[1])
		value = float32(tmpd) * scale
		return value
	case 14:
		tmpd = uint32(buf[3])
		tmpd = (tmpd << 8) & 0xff00
		tmpd += uint32(buf[2])
		tmpd = (tmpd << 8) & 0x00ffff00
		tmpd += uint32(buf[1])
		tmpd = (tmpd << 8) & 0xffffff00
		tmpd += uint32(buf[0])
		value = float32(tmpd) * scale
		return value
	case 15:
		tmpe[0] = buf[3]
		tmpe[1] = buf[2]
		tmpe[2] = buf[1]
		tmpe[3] = buf[0]
		value = byteTofloat32(tmpe) * scale
		return value
	case 16:
		tmpe[0] = buf[1]
		tmpe[1] = buf[0]
		tmpe[2] = buf[3]
		tmpe[3] = buf[2]
		value = byteTofloat32(tmpe) * scale
		return value
	case 17:
		tmpe[0] = buf[1]
		tmpe[1] = buf[0]
		tmpe[2] = buf[3]
		tmpe[3] = buf[2]
		value = byteTofloat32(tmpe) * scale
		return value
	case 18:
		tmpe[3] = buf[3]
		tmpe[2] = buf[2]
		tmpe[1] = buf[1]
		tmpe[0] = buf[0]
		value = byteTofloat32(tmpe) * scale
		return value
	case 19:
		tmpf = uint64(buf[0])
		tmpf = (tmpf << 8) & 0xffffffffffffff00
		tmpf += uint64(buf[1])
		tmpf = (tmpf << 8) & 0xffffffffffffff00
		tmpf += uint64(buf[2])
		tmpf = (tmpf << 8) & 0xffffffffffffff00
		tmpf += uint64(buf[3])
		tmpf = (tmpf << 8) & 0xffffffffffffff00
		tmpf += uint64(buf[4])
		tmpf = (tmpf << 8) & 0xffffffffffffff00
		tmpf += uint64(buf[5])
		tmpf = (tmpf << 8) & 0xffffffffffffff00
		tmpf += uint64(buf[6])
		tmpf = (tmpf << 8) & 0xffffffffffffff00
		tmpf += uint64(buf[7])
		tt := int64(tmpf)
		value = float32(tt) * scale
		return value
	default:
		return value
	}
	//0~19
	return value
}

//byteNum 获得字节个数
func byteNumTcp(p modbustcpdata) int {
	pt := p.ptype
	if pt == 0 {
		return 2
	} else if 0 < pt && pt <= 6 {
		return 2
	} else if 6 < pt && pt <= 18 {
		return 4
	} else {
		return 8
	}
}

//Init modbus规约初始化
func (m *PotocolModbusTcp) Init(cc config.ChannelConfig, dm device.DeviceMap) *PotocolModbusTcp {
	m.ChannelConfig = cc
	m.DeviceMap = dm
	return m
}
