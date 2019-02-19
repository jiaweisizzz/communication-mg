package potocol

import (
	"bytes"
	"container/list"
	"encoding/binary"
	"fmt"
	"math"
	"mg-coms/command"
	"mg-coms/config"
	"mg-coms/device"
	"sort"
	"strconv"
	"time"
)

//PotocolNewModbus  modbus规约
type PotocolNewModbus struct {
	PotocolParm
	Pm      ModbusMap       //分区
	Keys    []int           //分区号表
	Msgs    *list.List      //消息队列
	Message command.Message //消息
	St      *State          //状态
	Ds      *DeviceState
}

//ModbusPartition 规约分区
type ModbusPartition struct {
	name         string //分区名称
	terminalNum  int    //终端号
	partitionNum int    //分区号
	startAddr    []byte //起始地址
	endAddr      []byte //终止地址
	funCode      byte   //功能码
	registerNum  []byte //寄存器个数
	point        []data //点表
	timeout      string //分区的超时时间
}

//data 分区内点的配置
type data struct {
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

//ModbusMap 分区列表
type ModbusMap map[int]*ModbusPartition

//DeviceState 设备状态
type DeviceState struct {
	ds map[int]*OnOffline
}

func (m *DeviceState) init(c config.ChannelConfig) *DeviceState {
	m.ds = make(map[int]*OnOffline)
	for _, dev := range c.Device {
		m.ds[dev.ID] = new(OnOffline).init(dev)
	}
	return m
}

func (m *DeviceState) reset() {
	for _, ds := range m.ds {
		ds.Number = 0
		ds.State = true
	}
}

func (m *DeviceState) isRun(index int) bool {
	return m.ds[index].State
}

func (m *DeviceState) getTimeoutNum(index int) int {
	return m.ds[index].Number
}

func (m *DeviceState) addTimeoutNum(nm *PotocolNewModbus) {
	index := nm.St.getDeviceNum()
	m.ds[index].Number++
	if m.ds[index].Number > 9 {
		m.ds[index].State = false
		nm.DeviceMap[nm.ChannelConfig.ID*10000+index].Status = false
	}
}

func (m *DeviceState) clearTimeoutNum(index int) {
	m.ds[index].Number = 0
}

type OnOffline struct {
	ID     int
	State  bool
	Number int
}

func (l *OnOffline) init(d config.DeviceConfig) *OnOffline {
	l.ID = d.ID
	l.State = d.Status
	l.Number = 0
	return l
}

//State 状态
type State struct {
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

func (s *State) init(maxCursor int) *State {
	s.frame = make([]byte, 0, 1024)
	s.maxCursor = maxCursor
	s.cursor = 0
	s.timeout = "500ms"
	return s
}

func (s *State) reset() {
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

func (s *State) isDone() bool {
	return s.done
}

func (s *State) updateRxTime() {
	s.rxTime = time.Now()
}

func (s *State) updateTxTime() {
	s.txTime = time.Now()
}

func (s *State) setDone(b bool) {
	s.done = b
}
func (s *State) isCmd() bool {
	return s.cmdState
}
func (s *State) setCmd(b bool) {
	s.cmdState = b
}

func (s *State) setFrame(b []byte) {
	s.frame = b
}

func (s *State) getDeviceNum() int {
	return int(s.frame[0])
}

func (s *State) getCursor() int {
	return s.cursor
}

func (s *State) addCursor() {
	s.cursor++
	if s.cursor == s.maxCursor {
		s.cursor = 0
	}
}

func (s *State) getFrameData(index int) byte {
	return s.frame[index]
}

func (s *State) getCurrentPartition() int {
	return s.currentPartition
}

func (s *State) isTimeout() bool {
	s.updateRxTime()
	t, _ := time.ParseDuration(s.timeout)
	// fmt.Println("tttttttttttttttttttimeout", t)
	// fmt.Println("istimeout=rx-tx", s.rxTime.Sub(s.txTime).String())
	return s.rxTime.After(s.txTime.Add(t))
}

func (s *State) getTimeout() string {
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
func (m *PotocolNewModbus) Created() {
	m.Msgs = list.New()
	// m.Msgs = nil
	m.Message = command.Message{}
	m.Pm = make(map[int]*ModbusPartition)
	m.Ds = new(DeviceState).init(m.ChannelConfig)
	m.initPartitionTable(m.ChannelConfig)
	m.partitionSort()
	m.St = new(State).init(len(m.Keys))
	go m.timer()
}

//timer 定时器
func (m *PotocolNewModbus) timer() {
	//每30min重置设备在线状态
	for range time.Tick(5 * time.Minute) {
		m.Ds.reset()
		for index, _ := range m.DeviceMap {
			m.DeviceMap[index].Status = true
		}
	}
}

//initPartitionTable 格式化分区
func (m *PotocolNewModbus) initPartitionTable(c config.ChannelConfig) {
	for _, dev := range c.Device {
		m.partitionConfig(dev.Ycs, dev, 0)
		m.partitionConfig(dev.Yxs, dev, 1)
		m.partitionConfig(dev.Yms, dev, 2)
	}
}

//partitionConfig 分区配置
func (m *PotocolNewModbus) partitionConfig(pts []config.Point, dev config.DeviceConfig, kind int) {
	for _, p := range pts {
		d := data{}
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

func (s ModbusPartition) Len() int      { return len(s.point) }
func (s ModbusPartition) Swap(i, j int) { s.point[i], s.point[j] = s.point[j], s.point[i] }
func (s ModbusPartition) Less(i, j int) bool {
	return bytetoint(s.point[i].addr) < bytetoint(s.point[j].addr)
}

//updatePartition 更新分区
func (m *PotocolNewModbus) updatePartition(p config.Point, d data, pnum int) {
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
func (m *PotocolNewModbus) newPartition(p config.Point, dev config.DeviceConfig, d data) {
	pm := new(ModbusPartition)
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

//Cmd消息处理
func (m *PotocolNewModbus) Cmd(msg command.Message) {
	m.Msgs.PushBack(msg)
	fmt.Println("#####################################################")
}

//Do 轮询收发
func (m *PotocolNewModbus) Do() {
	m.tx()
	m.rx()
	time.Sleep(time.Millisecond * 50)
}

//rx 接收
func (m *PotocolNewModbus) rx() {

	//未完成返回
	if !m.St.isDone() {
		return
	}

	//超时返回

	buf, err := m.Port.Read()

	if err != nil {
		// fmt.Println("err", err)
		return
	}
	// fmt.Println("buf====", buf)

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
		if len(m.Buffer) < 6 {
			fmt.Println("&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&")
			return
		}

		//寻找完整一帧数据
		found := false
		index := 0
		for index, _ = range m.Buffer {
			//寻找首尾
			if m.Buffer[index] == m.St.getFrameData(0) && m.Buffer[index+1] == m.St.getFrameData(1) {
				//判断crc
				crc := make([]byte, 0, 10)
				if len(m.Buffer) < index+int(m.Buffer[index+2])+5 {
					break
				}
				v := intTobyte(int(crc16(m.Buffer[index:index+int(m.Buffer[index+2])+3], len(m.Buffer[index:index+int(m.Buffer[index+2])+3]))))
				if m.ChannelConfig.Crc == "高" {
					crc = append(crc, v[2])
					crc = append(crc, v[3])
				} else {
					crc = append(crc, v[3])
					crc = append(crc, v[2])
				}
				recrc := m.Buffer[index+int(m.Buffer[index+2])+3 : index+int(m.Buffer[index+2])+5]
				if crc[0] != recrc[0] && crc[1] != recrc[1] {
					m.Buffer = m.Buffer[index+1:]
					break
				}
				if (int(m.Buffer[index+2])) != (int(m.St.getFrameData(5))*2) && int(m.Buffer[index+1]) != 1 && int(m.Buffer[index+1]) != 2 {
					m.Buffer = nil
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

		var md *ModbusPartition

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
		a := index + 3
		var b int
		for _, dt := range md.point {
			// LOOP:
			if int(byteToint(dt.addr)) == mark {
				b = a + byteNum(dt)
				if b >= len(m.Buffer) {
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
					var strs []string
					if int(m.Buffer[1]) == 1 || int(m.Buffer[1]) == 2 {
						// str := bytetoBinaryString(m.Buffer[a+1]) + bytetoBinaryString(m.Buffer[a])  //青岛版本
						str := bytetoBinaryString(m.Buffer[a]) + bytetoBinaryString(m.Buffer[a+1]) //民勤版本
						for i := 15; i >= 0; i-- {
							strs = append(strs, string(str[i]))
						}
					} else {
						str := bytetoBinaryString(m.Buffer[a]) + bytetoBinaryString(m.Buffer[a+1])
						for i := 15; i >= 0; i-- {
							strs = append(strs, string(str[i]))
						}
					}
					// fmt.Println("strs", strs)
					for index, b := range strs {
						for _, d := range md.point {
							if int(byteToint(d.addr)) == mark && int(d.bit) == index {
								if string(b) == "1" {
									m.DeviceMap[m.ChannelConfig.ID*10000+md.terminalNum].SetYx(d.id, true)
								} else {
									m.DeviceMap[m.ChannelConfig.ID*10000+md.terminalNum].SetYx(d.id, false)
								}
							}
						}
					}
					break
				case 2:
					m.DeviceMap[m.ChannelConfig.ID*10000+md.terminalNum].SetYm(dt.id, data)
					break
				}
			} else {
				continue
			}
			mark = mark + byteNum(dt)/2
			a = b
		}
		if m.Ds.isRun(m.St.getDeviceNum()) {
			m.Ds.clearTimeoutNum(m.St.getDeviceNum())
		}
		m.Buffer = m.Buffer[index+int(m.Buffer[index+2])+5:]

		md.timeout = m.St.getTimeout()
		m.St.reset()

		//帧数据处理完成返回
		return
	}

	//处理命令
	// fmt.Println("___________________________________", time.Now(), m.St.selecttime)

	fmt.Println("msgcontrol", m.St.msgcontrol.SubCmd)
	fmt.Println("MESSSSSSSSSSSS", m.Message)

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

// var i = 0
// var message command.Message
var heartnum = 0

//tx 发送
func (m *PotocolNewModbus) tx() {
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
									frame = append(frame, byte(dev.ID))
									fc := p.GetIntData("Funcode")
									frame = append(frame, byte(fc))
									paddr := p.GetIntData("Addr")
									for index, st := range intTobyte(paddr) {
										if index > 1 {
											frame = append(frame, st)
										}
									}
									ptype := p.GetIntData("Type")
									if fc == 16 {
										if 0 < ptype && ptype < 7 {
											for index, rt := range intTobyte(registerNum(ptype)) {
												if index > 1 {
													frame = append(frame, rt)
												}
											}
											frame = append(frame, byte(2))
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
									crc := intTobyte(int(crc16(frame, len(frame))))
									if m.ChannelConfig.Crc == "高" {
										frame = append(frame, crc[2])
										frame = append(frame, crc[3])
									} else {
										frame = append(frame, crc[3])
										frame = append(frame, crc[2])
									}
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
									frame = append(frame, byte(dev.ID))
									fc := p.GetIntData("Funcode")
									frame = append(frame, byte(fc))
									paddr := p.GetIntData("Addr")
									for index, st := range intTobyte(paddr) {
										if index > 1 {
											frame = append(frame, st)
										}
									}
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
									for _, f := range m.byteSort(ptype, float64(p.Scale)) {
										frame = append(frame, f)
									}
									crc := intTobyte(int(crc16(frame, len(frame))))
									if m.ChannelConfig.Crc == "高" {
										frame = append(frame, crc[2])
										frame = append(frame, crc[3])
									} else {
										frame = append(frame, crc[3])
										frame = append(frame, crc[2])
									}
									err := m.Port.Write(string(frame), len(frame))
									if err != nil {
										fmt.Println(err)
									}
									// m.St.updateTxTime()
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
		if len(m.Pm) == 0 {
			return
		}
		pm := m.Pm[keys[cur]]
		m.St.timeout = pm.timeout //------
		// fmt.Println("i==", cur)
		// fmt.Println("最大长度", len(keys))
		//
		fmt.Println("分区", keys[cur])
		m.St.currentPartition = keys[cur]
		frame = append(frame, byte(pm.terminalNum))
		frame = append(frame, pm.funCode)
		for _, staddr := range pm.startAddr {
			frame = append(frame, staddr)
		}
		switch int(pm.funCode) {
		case 1:
			redisternum := intTobyte(16)
			for index, rtnum := range redisternum {
				if index > 1 {
					frame = append(frame, rtnum)
				}
			}
		case 2:
			redisternum := intTobyte(16)
			for index, rtnum := range redisternum {
				if index > 1 {
					frame = append(frame, rtnum)
				}
			}
		case 6:
			if heartnum == 0 {
				for index, rtnum := range intTobyte(0) {
					if index > 1 {
						frame = append(frame, rtnum)
					}
				}
				heartnum = 1
			}
			if heartnum == 1 {
				for index, rtnum := range intTobyte(1) {
					if index > 1 {
						frame = append(frame, rtnum)
					}
				}
				heartnum = 0
			}
		default:
			for _, rtnum := range pm.registerNum {
				frame = append(frame, rtnum)
			}
		}

		crc := intTobyte(int(crc16(frame, len(frame))))
		if m.ChannelConfig.Crc == "高" {
			frame = append(frame, crc[2])
			frame = append(frame, crc[3])
		} else {
			frame = append(frame, crc[3])
			frame = append(frame, crc[2])
		}
		m.St.setFrame(frame)
		if m.Ds.isRun(pm.terminalNum) == false {
			m.St.addCursor()
			return
		}
		fmt.Println("frame-----", frame)
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

//partitionSort  对分区进行排序获得键值
func (m *PotocolNewModbus) partitionSort() {
	var keys []int
	// fmt.Println("获得键值时", m.Pm)
	for k := range m.Pm {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	m.Keys = keys
}

//byteSort 字节排序
func (m *PotocolNewModbus) byteSort(ptype int, scale float64) []byte {
	var frame = make([]byte, 0, 1024)
	switch ptype {
	case 0:
		msg := int(m.Message.Value.(float64) * scale)
		for index, v := range intTobyte(msg) {
			if index > 1 {
				frame = append(frame, v)
			}
		}
		break
	case 1:
		msg := int(m.Message.Value.(float64) * scale)
		for index, v := range intTobyte(msg) {
			if index > 1 {
				frame = append(frame, v)
			}
		}
		break
	case 2:
		msg := int(m.Message.Value.(float64) * scale)
		v := intTobyte(msg)
		frame = append(frame, v[3])
		frame = append(frame, v[2])
		break
	case 3:
		msg := int(m.Message.Value.(float64) * scale)
		for index, v := range intTobyte(msg) {
			if index > 1 {
				frame = append(frame, v)
			}
		}
		break
	case 4:
		msg := int(m.Message.Value.(float64) * scale)
		v := intTobyte(msg)
		frame = append(frame, v[3])
		frame = append(frame, v[2])
		break
	case 5:
		msg := int(m.Message.Value.(float64) * scale)
		for index, v := range intTobyte(msg) {
			if index > 1 {
				frame = append(frame, v)
			}
		}
		break
	case 6:
		msg := int(m.Message.Value.(float64) * scale)
		v := intTobyte(msg)
		frame = append(frame, v[3])
		frame = append(frame, v[2])
		break
	case 7:
		msg := int(m.Message.Value.(float64) * scale)
		v := intTobyte(msg)
		frame = append(frame, v[0])
		frame = append(frame, v[1])
		frame = append(frame, v[2])
		frame = append(frame, v[3])
		break
	case 8:
		msg := int(m.Message.Value.(float64) * scale)
		v := intTobyte(msg)
		frame = append(frame, v[1])
		frame = append(frame, v[0])
		frame = append(frame, v[3])
		frame = append(frame, v[2])
		break
	case 9:
		msg := int(m.Message.Value.(float64) * scale)
		v := intTobyte(msg)
		frame = append(frame, v[2])
		frame = append(frame, v[3])
		frame = append(frame, v[0])
		frame = append(frame, v[1])
		break
	case 10:
		msg := int(m.Message.Value.(float64) * scale)
		v := intTobyte(msg)
		frame = append(frame, v[3])
		frame = append(frame, v[2])
		frame = append(frame, v[1])
		frame = append(frame, v[0])
		break
	case 11:
		msg := int(m.Message.Value.(float64) * scale)
		v := intTobyte(msg)
		frame = append(frame, v[0])
		frame = append(frame, v[1])
		frame = append(frame, v[2])
		frame = append(frame, v[3])
		break
	case 12:
		msg := int(m.Message.Value.(float64) * scale)
		v := intTobyte(msg)
		frame = append(frame, v[1])
		frame = append(frame, v[0])
		frame = append(frame, v[3])
		frame = append(frame, v[2])
		fmt.Printf("%x", frame)
		break
	case 13:
		msg := int(m.Message.Value.(float64) * scale)
		v := intTobyte(msg)
		frame = append(frame, v[2])
		frame = append(frame, v[3])
		frame = append(frame, v[0])
		frame = append(frame, v[1])
		break
	case 14:
		msg := int(m.Message.Value.(float64) * scale)
		v := intTobyte(msg)
		frame = append(frame, v[3])
		frame = append(frame, v[2])
		frame = append(frame, v[1])
		frame = append(frame, v[0])
		break
	case 15:
		msg := float32(m.Message.Value.(float64) * scale)
		v := float32Tobyte(msg)
		fmt.Println(v)
		frame = append(frame, v[0])
		frame = append(frame, v[1])
		frame = append(frame, v[2])
		frame = append(frame, v[3])
		break
	case 16:
		msg := float32(m.Message.Value.(float64) * scale)
		v := float32Tobyte(msg)
		frame = append(frame, v[1])
		frame = append(frame, v[0])
		frame = append(frame, v[3])
		frame = append(frame, v[2])
		break
	case 17:
		msg := float32(m.Message.Value.(float64) * scale)
		v := float32Tobyte(msg)
		frame = append(frame, v[2])
		frame = append(frame, v[3])
		frame = append(frame, v[0])
		frame = append(frame, v[1])
		break
	case 18:
		msg := float32(m.Message.Value.(float64) * scale)
		v := float32Tobyte(msg)
		frame = append(frame, v[3])
		frame = append(frame, v[2])
		frame = append(frame, v[1])
		frame = append(frame, v[0])
		break
	default:
		msg := int(m.Message.Value.(float64) * scale)
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
func (m *PotocolNewModbus) calculate(buf []byte, ptype int, scale float32) float32 {
	// fmt.Println("ptype", ptype)
	var mmm int16
	var tmpa uint16
	var tmpaa int16
	var tmpb uint16
	var tmpc uint32
	var tmpd uint32
	var tmpe = make([]byte, 10)
	var tmpf uint64
	var value float32
	switch ptype {
	case 0:
	case 1:
		mmm = int16(buf[0])
		mmm = mmm << 8
		mmm += int16(buf[1])
		tt := float32(int16(mmm))
		value = tt * scale
		return value

		// mmm = int16(buf[0])
		// mmm = mmm << 8
		// mmm += int16(buf[1])
		// tt := float32(int16(mmm))
		// if tt < 0 {
		// 	value = -32768 - tt
		// } else {
		// 	value = tt
		// }
		// value = value * scale
		// return value
	case 2:
		mmm = int16(buf[1])
		mmm = mmm << 8
		mmm += int16(buf[0])
		tt := float32(int16(mmm))
		value = tt * scale
		return value
		// tmpa = uint16(buf[1])
		// tmpa = (tmpa << 8) & 0xff00
		// tmpa += uint16(buf[0])
		// tt := int16(tmpa)
		// value = float32(tt) * scale
		// return value
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
		// tmpa = uint16(buf[0])
		// tmpa = (tmpa << 8) & 0xff00
		// tmpa += uint16(buf[1])
		// value = float32(int(tmpa)) * scale
		tmpaa = int16(buf[0])
		tmpaa = tmpaa << 8
		tmpaa += int16(buf[1])
		// tt := float32(int16(mmm))
		// if tt < 0 {
		// 	value = -32768 - tt
		// } else {
		// 	value = tt
		// }
		value = float32(int16(tmpaa)) * scale
		return value
	case 6:
		tmpb = uint16(buf[1])
		tmpb = (tmpb << 8) & 0xff00
		tmpb += uint16(buf[0])
		value = float32(int(tmpb)) * scale
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
		fmt.Println(value)
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
		//寄存器由高到低 3 2 1 0
		tmpe[0] = buf[3]
		tmpe[1] = buf[2]
		tmpe[2] = buf[1]
		tmpe[3] = buf[0]
		value = byteTofloat32(tmpe) * scale
		return value
	case 16:
		tmpe[0] = buf[2]
		tmpe[1] = buf[3]
		tmpe[2] = buf[0]
		tmpe[3] = buf[1]
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
		tmpe[0] = buf[0]
		tmpe[1] = buf[1]
		tmpe[2] = buf[2]
		tmpe[3] = buf[3]
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

//registerNum 获得寄存器个数
func registerNum(ptype int) int {
	if 0 <= ptype && ptype <= 6 {
		return 1
	} else if 6 < ptype && ptype < 19 {
		return 2
	} else {
		return 4
	}
}

//byteNum 获得字节个数
func byteNum(p data) int {
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

//crc16 crc校验
func crc16(data []byte, len int) uint16 {
	var crc, x, y uint16
	crc = 0xffff
	if len == 0 {
		return crc
	}
	var temp byte
	len--
	for x = 0; int(x) <= len; x++ {
		temp = (byte(crc) & 0x00ff) ^ data[x]
		crc = (crc & 0xff00) + uint16(temp)
		for y = 0; y <= 7; y++ {
			if (crc & 1) != 0 {
				crc = (crc >> 1) ^ 0xA001
			} else {
				crc = crc >> 1
			}
		}
	}
	return crc
}

//byteToint 字节转整型
func byteToint(b []byte) int16 {
	buf := bytes.NewBuffer(b)
	var x int16
	binary.Read(buf, binary.BigEndian, &x)
	return x
}

//byteToint32 字节转int32
func byteToint32(b []byte) int32 {
	buf := bytes.NewBuffer(b)
	var x int32
	binary.Read(buf, binary.BigEndian, &x)
	return x
}

//intTobyte 整型转字节
func intTobyte(x int) []byte {
	n := int32(x)
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, n)
	return buf.Bytes()
}

//int64Tobyte int64转byte
func int64Tobyte(x int) []byte {
	n := int64(x)
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, n)
	return buf.Bytes()
}

//stringTofloat32 字符串转float32
func stringTofloat32(s string) float32 {
	f, _ := strconv.ParseFloat(s, 32)
	return float32(f)
}

//float64Tobyte float64转字节
func float64Tobyte(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)
	return bytes
}

//float32Tobyte float32转byte
func float32Tobyte(float float32) []byte {
	bits := math.Float32bits(float)
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, bits)
	return bytes
}

//byteTofloat32 字节转float32
func byteTofloat32(bytes []byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes)
	return math.Float32frombits(bits)
}

func bytetoBinaryString(data byte) string {
	var a byte
	var str string
	for i := 0; i < 8; i++ {
		a = data
		data <<= 1
		data >>= 1

		switch a {
		case data:
			str += "0"
		default:
			str += "1"
		}
		data <<= 1
	}
	return str
}

//Init modbus规约初始化
func (m *PotocolNewModbus) Init(cc config.ChannelConfig, dm device.DeviceMap) *PotocolNewModbus {
	m.ChannelConfig = cc
	m.DeviceMap = dm
	return m
}
