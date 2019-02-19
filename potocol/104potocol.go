package potocol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"mg-coms/command"
	"mg-coms/config"
	"mg-coms/device"
	"time"
)

//发送序号
const (
	ADDR = 0x01
	STAR = 0x68 //启动符
	U    = 0x03 //U格式
	S    = 0x01 //S格式
	I    = 0x00 //I格式
)

//控制域
const (
	_uRxControl     = 0x07 //U帧接收
	_uTxControl     = 0x0B //U帧确认帧
	_uTxControlTest = 0x83 //U帧测试生效帧
	_uRxControlTest = 0x43 //接收U帧测试帧
)

//类型标识
const (
	_SinglePointYxKind = 0x01 //单点遥信类型标识
	_SinglePointYxTime = 0x02 //带时标单点遥信标识
	_DoublePointYx     = 0x03 //双点遥信类型标识
	_DoublePointYxTime = 0x04 //带时标双点遥信标识
	_SetYk             = 0x2E //遥控命令
	_CallAll           = 0x64 //总召唤
	_CallYm            = 0x65 //电能量
	_YcKind            = 0x09 //遥测类型标识
	_YmKind            = 0x0F //遥脉类型标识
	_ChangeYcKind      = 0x0D //遥测变位上送标识
	_ChangeYxKind      = 0x01 //遥信变位上送标识
	_SetPonitYk        = 0x2D //单点遥控类型标识
	_SetYt             = 0x32 //遥调类型标识
)

//原因
const (
	RxCallAllReason        = 0x06 //接收召唤传送原因
	TxCallAllReason        = 0x07 //发送召唤传送原因
	TxYmReason             = 0x05 //发送遥脉传送原因
	RxYkPresetReason       = 0x06 //接收遥控预置传送原因
	TxYkPresetReason       = 0x07 //发送遥控预置传送原因
	RxYkExecute            = 0x06 //接收遥控执行传送原因
	TxYkExecute            = 0x07 //发送遥控执行传送原因
	RxYKCancel             = 0x08 //接收遥控取消传送原因
	TxYKCannel             = 0x09 //发送遥控取消传送原因
	TxYcReason             = 0x03 //发送遥测传送原因
	TxYxReason             = 0x14 //发送遥信传送原因
	TxDisplacementYxReason = 0x03 //发送变位遥信传送原因
	_TxEndReason           = 0x0A //发送结束帧传送原因
	_ChangeYcReason        = 0x03 //变位上送传送原因
	_ChangeYxReason        = 0x03 //变位遥信传送原因
)

//Potocol104 104规约
type Potocol104 struct {
	PotocolParm
	imageFrame    *Frame
	frame         *Frame
	txNo          int //发送计数
	rxNo          int //接收计数
	txTime        time.Time
	ready         bool //端口是否正常
	changeReady   bool //
	number        int
	selected      bool
	changeAllTime time.Time
	changeAddTime time.Time
}

func (p *Potocol104) getTxNo() int {
	no := p.txNo
	p.txNo += 2
	return no
}

func (p *Potocol104) getRxNo() int {
	no := p.rxNo
	p.rxNo += 0
	return no
}

//Frame 接收帧

type Frame struct {
	start  byte   //启动符
	len    byte   //长度
	txnum  []byte //发送计数
	rxnum  []byte //接收计数
	kind   byte   //类型标识
	limit  byte   //可变结构限定词
	reason []byte //传动原因
	addr   []byte //公共地址
	data   []byte //内容
}

//Init 104
func (p *Potocol104) Init(cc config.ChannelConfig, dm device.DeviceMap) *Potocol104 {
	p.ChannelConfig = cc
	p.DeviceMap = dm
	p.ready = false

	return p
}

//Created 104规约初始化
func (p *Potocol104) Created() {
	p.ready = false
	p.changeReady = false
	p.frame = new(Frame)
	p.frame.initFrame()
	p.changeAllTime = time.Now()
	h, _ := time.ParseDuration("15m")
	p.changeAddTime = p.changeAllTime.Add(h)
}

//initFrame 初始化Frame
func (f *Frame) initFrame() {
	f.start = 0x68
	f.txnum = make([]byte, 2)
	f.rxnum = make([]byte, 2)
	f.reason = make([]byte, 2)
	f.addr = make([]byte, 2)
	f.data = make([]byte, 0)
}

//Cmd 命令处理
func (p *Potocol104) Cmd(msg command.Message) {
}

//Do 循环处理
func (p *Potocol104) Do() {

	//读取数据
	buf, err := p.Port.Read()

	if err == nil {
		p.ready = true
	} else {
		p.ready = false
	}

	if p.ready == false {
		time.Sleep(1 * time.Millisecond)
		return
	}

	p.Buffer = append(p.Buffer, buf...)

	if len(p.Buffer) != 0 {
		p.changeReady = false
	}

	// //没有后台数据进行YX/YC变位上送
	if p.changeReady == true {
		if len(p.Buffer) == 0 {
			if p.changeAllTime.After(p.changeAddTime) {
				now := time.Now()
				p.changeYx(false)
				p.changeYc(false)
				p.txTime = now
				p.changeAllTime = time.Now()
				h, _ := time.ParseDuration("15m")
				p.changeAddTime = p.changeAllTime.Add(h)
				return
			}
			now := time.Now()
			p.changeYx(true)
			p.changeYc(true)
			p.txTime = now
			p.changeAllTime = time.Now()
			return
		}
	}
	// labal: //解析帧
	tmp, err := p.parse(p.Buffer)
	if err != nil {
		return
	}

	p.Buffer = tmp

	switch p.frame.txnum[0] & 0x03 {
	case U: //U帧
		p.uFrame(p.frame)
	case S: //S帧
		p.changeReady = true
	default: //I帧
		p.iFrame(p.frame)
	}
	p.ready = true
	//清理现场
	p.frame = new(Frame)
	p.frame.initFrame()

}

//getFrame 获得完整帧
func (f *Frame) getFrame(b bool) []byte {
	tmp := make([]byte, 0)
	tmp = append(tmp, f.start)
	tmp = append(tmp, f.len)
	tmp = append(tmp, f.txnum...)
	tmp = append(tmp, f.rxnum...)
	if !b {
		tmp = append(tmp, f.kind)
		tmp = append(tmp, f.limit)
		tmp = append(tmp, f.reason...)
		tmp = append(tmp, f.addr...)
		tmp = append(tmp, f.data...)
	}

	return tmp
}

//parse 解析帧格式
func (p *Potocol104) parse(buf []byte) ([]byte, error) {
	found := false
	for i := range buf {
		//不断向前移位找到0x68
		if buf[i] == STAR {
			buf = buf[i:]
			found = true
			break
		}
	}

	if !found {
		return nil, errors.New("0x68 error")
	}

	if len(buf) < 2 {
		return buf, errors.New("buf length < 2 error")
	}

	if int(buf[1])+2 > len(buf) {
		return buf, errors.New("buf length < int(buf[1])+2 error")
	}

	p.frame.start = buf[0]
	p.frame.len = buf[1]
	p.frame.txnum = buf[2:4]
	p.frame.rxnum = buf[4:6]

	switch p.frame.txnum[0] & 0x01 {
	case I:
		p.frame.kind = buf[6]
		p.frame.limit = buf[7]
		p.frame.reason = buf[8:10]
		p.frame.addr = buf[10:12]
		p.frame.data = buf[12 : int(buf[1])+2]
	}

	return buf[int(buf[1])+2:], nil
}

//updateLength 获得完整帧
func (f *Frame) updateLength(b bool) {
	f.len = intTobyte(len(f.getFrame(b)) - 2)[3]
	if f.kind == _SinglePointYxKind && f.reason[0] == TxYxReason {
		f.limit = intTobyte((len(f.getFrame(b)) - 12) / 4)[3]
		return
	}
	if f.kind == _YcKind && f.reason[0] == TxYcReason {
		f.limit = intTobyte((len(f.getFrame(b)) - 12) / 6)[3]
		return
	}
	if f.kind == _YmKind && f.reason[0] == TxYmReason {
		f.limit = intTobyte((len(f.getFrame(b)) - 12) / 8)[3]
		return
	}
	if f.kind == _ChangeYcKind && f.reason[0] == _ChangeYcReason {
		f.limit = intTobyte((len(f.getFrame(b)) - 12) / 8)[3]
		return
	}
	if f.kind == _ChangeYxKind && f.reason[0] == _ChangeYxReason {
		f.limit = intTobyte((len(f.getFrame(b)) - 12) / 4)[3]
		return
	}

}

//iFrame 判断帧类型执行处理
func (p *Potocol104) iFrame(frame *Frame) {
	switch frame.kind {
	case _CallAll:
		p.callAll(frame)
		p.changeReady = true
	case _CallYm:
		p.callYm(frame)
		p.changeReady = true
	case _SetYk:
		p.setYk(frame)
	case _SetPonitYk:
		p.setPonitYk(frame)
	case _SetYt:
		p.setYt(frame)
	}
}

//uFrame 发送U帧数据
func (p *Potocol104) uFrame(f *Frame) {
	frame := *f
	switch frame.txnum[0] {
	case _uRxControl:
		frame.txnum[0] = _uTxControl
	case _uRxControlTest:
		p.changeReady = true
		frame.txnum[0] = _uTxControlTest
	}

	p.write(true, frame)
}

//sFrame 发送S帧数据
func (p *Potocol104) sFrame(f *Frame) {

	frame := *f
	rxnum := make([]byte, 2)
	rxnum[0] = frame.rxnum[1]
	rxnum[1] = frame.rxnum[0]
	rxnum = intTobyte(int(byteToint(rxnum)) + 0)
	frame.rxnum[0] = rxnum[3]
	frame.rxnum[1] = rxnum[2]

	rn := make([]byte, 2)
	rn[0] = f.rxnum[1]
	rn[1] = f.rxnum[0]
	p.rxNo = int(byteToint(rn))
	p.write(true, frame)
}

//callAll 总召唤帧处理
func (p *Potocol104) callAll(f *Frame) {
	ff := *f
	p.image(f)
	p.allYx(f)
	p.allYc(f)
	p.end(f, &ff)
	p.txTime = time.Now()
}

//image 发送召唤镜像
func (p *Potocol104) image(f *Frame) {
	f.setImageNum()
	frame := *f
	frame.reason[0] = TxCallAllReason
	p.write(false, frame)
}

//allYx 发送遥信帧
func (p *Potocol104) allYx(f *Frame) {
	frame := *f
	frame.kind = _SinglePointYxKind
	frame.reason[0] = TxYxReason
	data := p.fillYx()
	for {
		if len(data) > 240 {
			frame.data = data[:240]
			data = data[240:]
			f.addTxNum()
			p.write(false, frame)
		} else {
			frame.data = data
			f.addTxNum()
			p.write(false, frame)
			return
		}
	}

}

//allYc 发送遥测帧
func (p *Potocol104) allYc(f *Frame) {
	// f.addTxNum()
	frame := *f
	frame.kind = _YcKind
	frame.reason[0] = TxYcReason
	data := p.fillYc()
	for {
		if len(data) > 240 {
			f.addTxNum()
			frame.data = data[:240]
			data = data[240:]
			p.write(false, frame)
		} else {
			f.addTxNum()
			frame.data = data
			p.write(false, frame)
			return
		}
	}
}

//callYm 电度量遥脉帧处理
func (p *Potocol104) callYm(f *Frame) {
	ff := *f
	p.image(f)
	p.allYm(f)
	p.end(f, &ff)
	p.txTime = time.Now()
}

//allYm 发送电度量遥脉数据
func (p *Potocol104) allYm(f *Frame) {
	frame := *f
	frame.kind = _YmKind
	frame.reason[0] = TxYmReason
	data := p.fillYm()
	for {
		if len(data) > 240 {
			f.addTxNum()
			frame.data = data[0:240]
			data = data[240:]
			p.write(false, frame)
		} else {
			f.addTxNum()
			frame.data = data
			p.write(false, frame)
			return
		}
	}

}

//end 发送结束帧
func (p *Potocol104) end(f, ff *Frame) {
	fmt.Println("总召唤========================================================", *ff)
	f.addTxNum()
	ff.txnum[0] = f.txnum[0]
	ff.rxnum[0] = f.rxnum[0]
	frame := *ff
	frame.reason[0] = 0x0A
	tn := make([]byte, 2)
	tn[0] = f.txnum[1]
	tn[1] = f.txnum[0]
	p.txNo = int(int(byteToint(tn))) + 2
	rn := make([]byte, 2)
	rn[0] = f.rxnum[1]
	rn[1] = f.rxnum[0]
	p.rxNo = int(byteToint(rn))
	p.write(false, frame)
}

//setYt 遥调帧处理
func (p *Potocol104) setYt(f *Frame) {
	if f.reason[0] == 0x06 && f.data[7] == 0x80 { //遥调激活确认
		p.ykPreset(f)
		p.executeYt(f)
		return
	}
	if f.reason[0] == 0x06 && f.data[7] == 0x00 { //遥调激活确认
		p.ykPreset(f)
		p.executeYt(f)
		p.ykEnd(f)
		p.changeReady = true
		p.txTime = time.Now()
		return
	}
	if f.reason[0] == 0x08 && f.data[7] == 0x80 { //遥调停止激活
		p.ykCancel(f)
		p.txTime = time.Now()
		p.changeReady = true

		return
	}
}

//setPonitYk 单点遥控帧处理
func (p *Potocol104) setPonitYk(f *Frame) {
	//根据传送原因和总召唤参数判断遥控类型
	if f.reason[0] == RxYkPresetReason && f.data[3] == 0x80 || f.reason[0] == RxYkPresetReason && f.data[3] == 0x81 { //0x82预置控合限定词
		p.ykPreset(f)
		p.executeYk(f)
		return
	}
	if f.reason[0] == RxYkExecute && f.data[3] == 0x00 || f.reason[0] == RxYkExecute && f.data[3] == 0x01 { //0x02执行控合限定词
		p.ykiFrame(f)
		p.executeYk(f)
		p.ykEnd(f)
		p.changeReady = true
		p.txTime = time.Now()
		return
	}
	if f.reason[0] == RxYKCancel && f.data[3] == 0x80 || f.reason[0] == RxYKCancel && f.data[3] == 0x81 { //0x08撤销限定词
		p.ykCancel(f)
		p.changeReady = true
		p.txTime = time.Now()
		return

	}
}

//setYk 双点遥控帧处理
func (p *Potocol104) setYk(f *Frame) {
	//根据传送原因和总召唤参数判断遥控类型
	if f.reason[0] == RxYkPresetReason && f.data[3] == 0x85 || f.reason[0] == RxYkPresetReason && f.data[3] == 0x86 { //0x82预置控合限定词
		p.ykPreset(f)
		p.executeYk(f)
		p.selected = true
		return
	}
	if f.reason[0] == RxYkExecute && f.data[3] == 0x05 || f.reason[0] == RxYkExecute && f.data[3] == 0x06 { //0x02执行控合限定词
		if !p.selected {
			p.fastExecuteYk(f)
			//自动选择后需要延时
			time.Sleep(time.Millisecond * 500)
		}
		p.ykiFrame(f)
		p.executeYk(f)
		p.ykEnd(f)
		p.changeReady = true
		p.txTime = time.Now()
		p.selected = false
		return
	}
	if f.reason[0] == RxYKCancel && f.data[3] == 0x85 || f.reason[0] == RxYKCancel && f.data[3] == 0x86 { //0x08撤销限定词
		p.ykCancel(f)
		p.txTime = time.Now()
		p.changeReady = true

		return

	}

}

//setImageNum 设置镜像帧发送和接收序列号
func (f *Frame) setImageNum() {
	rxnum := make([]byte, 2)

	rxnum[0] = f.txnum[1]
	rxnum[1] = f.txnum[0]
	f.txnum[0] = f.rxnum[0]
	f.txnum[1] = f.rxnum[1]

	rxnum = intTobyte(int(byteToint(rxnum)) + 2)
	f.rxnum[0] = rxnum[3]
	f.rxnum[1] = rxnum[2]
}

//addTxNum 添加发送序列号
func (f *Frame) addTxNum() {
	txnum := make([]byte, 2)
	txnum[0] = f.txnum[1]
	txnum[1] = f.txnum[0]
	txnum = intTobyte(int(byteToint(txnum)) + 2)
	f.txnum[0] = txnum[3]
	f.txnum[1] = txnum[2]
}

//changeYc 发送变位遥信
func (p *Potocol104) changeYc(mark bool) {
	if !p.ready {
		return
	}

	frame := new(Frame)
	frame.initFrame()
	frame.kind = _ChangeYcKind
	frame.reason[0] = _ChangeYcReason
	frame.addr[0] = 0x00
	data := p.fillChangeYc(mark)
	if len(data) == 0 {
		return
	}
	for {
		if len(data) > 240 {
			frame.setTxNum(p.getTxNo())
			frame.setRxNum(p.getRxNo())
			frame.data = data[:240]
			data = data[240:]
			p.write(false, *frame)
		} else {
			frame.setTxNum(p.getTxNo())
			frame.setRxNum(p.getRxNo())
			frame.data = data
			p.write(false, *frame)
			return
		}
	}
}

//changeYx 发送变位遥信
func (p *Potocol104) changeYx(mark bool) {
	if !p.ready {
		return
	}
	frame := new(Frame)
	frame.initFrame()
	frame.kind = _ChangeYxKind
	frame.reason[0] = _ChangeYxReason
	frame.addr[0] = 0x01
	buf := make([]byte, 0)
	buf = p.fillChangeYx(mark)
	if len(buf) == 0 {
		return
	}
	for {
		if len(buf) > 240 {
			frame.data = buf[:240]
			frame.setTxNum(p.getTxNo())
			frame.setRxNum(p.getRxNo())
			p.write(false, *frame)
			buf = buf[240:]
			frame.data = []byte{}
		} else {
			frame.data = append(frame.data, buf...)
			frame.setTxNum(p.getTxNo())
			frame.setRxNum(p.getRxNo())
			p.write(false, *frame)

			break
		}

	}
}

//ykEnd 遥控执行结束
func (p *Potocol104) ykEnd(f *Frame) {
	frame := f
	frame.reason[0] = 0x0a
	f.addTxNum()
	tn := make([]byte, 2)
	tn[0] = f.txnum[1]
	tn[1] = f.txnum[0]
	p.txNo = int(byteToint(tn)) + 2
	rn := make([]byte, 2)
	rn[0] = f.rxnum[1]
	rn[1] = f.rxnum[0]
	p.rxNo = int(byteToint(rn))
	p.write(false, *frame)
}

//ykPreset 发送遥控预置返校帧
func (p *Potocol104) ykPreset(f *Frame) {
	frame := f
	frame.reason[0] = TxYkPresetReason
	f.setImageNum()
	tn := make([]byte, 2)
	tn[0] = f.txnum[1]
	tn[1] = f.txnum[0]
	p.txNo = int(byteToint(tn)) + 2
	rn := make([]byte, 2)
	rn[0] = f.rxnum[1]
	rn[1] = f.rxnum[0]
	p.rxNo = int(byteToint(rn))
	p.write(false, *frame)

}

//ykiFrame 发送遥控执行确认帧
func (p *Potocol104) ykiFrame(f *Frame) {
	frame := f
	frame.reason[0] = TxYkPresetReason
	f.setImageNum()
	tn := make([]byte, 2)
	tn[0] = f.txnum[1]
	tn[1] = f.txnum[0]
	p.txNo = int(byteToint(tn)) + 2
	rn := make([]byte, 2)
	rn[0] = f.rxnum[1]
	rn[1] = f.rxnum[0]
	p.rxNo = int(byteToint(rn))

	p.write(false, *frame)
}

//ykCancel 发送遥控命令取消帧
func (p *Potocol104) ykCancel(f *Frame) {
	frame := f
	frame.reason[0] = TxYKCannel
	f.setImageNum()
	tn := make([]byte, 2)
	tn[0] = f.txnum[1]
	tn[1] = f.txnum[0]
	p.txNo = int(byteToint(tn)) + 2
	rn := make([]byte, 2)
	rn[0] = f.rxnum[1]
	rn[1] = f.rxnum[0]
	p.rxNo = int(byteToint(rn))

	p.write(false, *frame)

}

//setRxNum 接收序列号
func (f *Frame) setRxNum(count int) {
	tmp := intTobyte(count)
	f.rxnum[0] = tmp[3]
	f.rxnum[1] = tmp[2]
}

//setTxNum 发送序列号
func (f *Frame) setTxNum(count int) {
	tmp := intTobyte(count)
	f.txnum[0] = tmp[3]
	f.txnum[1] = tmp[2]
}

//fillYm 填充变位遥脉数据
func (p *Potocol104) fillYm() []byte {
	data := make([]byte, 0)
	//填写数据
	if p.ChannelConfig.Uptable.AddrYm != nil {
		for _, pt := range p.ChannelConfig.Uptable.AddrYm {
			dev := p.DeviceMap[pt.ChannelID*10000+pt.DeviceID]
			point, _ := dev.GetYmByTime(pt.ID, p.txTime)
			data = append(data, addYmData(pt.AcceptID, point.Data)...)
		}
	}

	return data
}

//fillYx 填充变位遥信数据
func (p *Potocol104) fillYx() []byte {
	data := make([]byte, 0)
	//填写数据
	if p.ChannelConfig.Uptable.AddrYx != nil {
		for _, pt := range p.ChannelConfig.Uptable.AddrYx {
			dev := p.DeviceMap[pt.ChannelID*10000+pt.DeviceID]
			point, _ := dev.GetYxByTime(pt.ID, p.txTime)

			data = append(data, addYxData(pt.AcceptID, point.Data)...)
		}
	}

	return data
}

func (p *Potocol104) fillChangeYx(mark bool) []byte {
	data := make([]byte, 0)
	//填写数据
	if p.ChannelConfig.Uptable.AddrYx != nil {
		for _, pt := range p.ChannelConfig.Uptable.AddrYx {
			dev := p.DeviceMap[pt.ChannelID*10000+pt.DeviceID]
			point, need := dev.GetYxByTime(pt.ID, p.txTime)
			if need == false && mark == true {
				continue
			}
			data = append(data, addYxData(pt.AcceptID, point.Data)...)
		}
	}

	return data
}

//executeYt执行遥调命令
func (p *Potocol104) executeYt(f *Frame) {
	fmt.Println("f.data", f.data)
	addr := make([]byte, 0)
	addr = append(addr, 0)
	addr = append(addr, f.data[2])
	addr = append(addr, f.data[1])
	addr = append(addr, f.data[0])
	data := make([]byte, 0)
	data = append(data, f.data[3])
	data = append(data, f.data[4])
	data = append(data, f.data[5])
	data = append(data, f.data[6])
	tmp := int(byteToint32(addr)) - 0x6200
	for _, pt := range p.ChannelConfig.Uptable.AddrYt {
		if tmp == pt.AcceptID {
			if f.data[7] == 0x80 && f.reason[0] == 0x07 {
				msg := command.Message{
					Cmd:     command.YT,
					SubCmd:  4,
					Value:   float64(byteTofloat32(data)),
					Device:  pt.DeviceID,
					Point:   pt.ID,
					Channel: p.DeviceMap[pt.ChannelID*10000+pt.DeviceID].Channel,
				}
				*p.Bus <- msg
			}
			if f.data[7] == 0x00 && f.reason[0] == 0x07 {
				msg := command.Message{
					Cmd:     command.YT,
					SubCmd:  5,
					Value:   float64(byteTofloat32(data)),
					Device:  pt.DeviceID,
					Point:   pt.ID,
					Channel: p.DeviceMap[pt.ChannelID*10000+pt.DeviceID].Channel,
				}
				*p.Bus <- msg
			}

		}
	}
}

//fastExecuteYk 超时执行遥控命令
func (p *Potocol104) fastExecuteYk(f *Frame) {
	addr := make([]byte, 0)
	addr = append(addr, 0)
	addr = append(addr, f.data[2])
	addr = append(addr, f.data[1])
	addr = append(addr, f.data[0])
	tmp := int(byteToint32(addr)) - 0x6000
	for _, pt := range p.ChannelConfig.Uptable.AddrYk {
		if tmp == pt.AcceptID {
			if f.kind == 0x2e {
				msg := command.Message{
					Cmd:     command.YK,
					SubCmd:  4,
					Value:   f.data[3] != 0x86,
					Device:  pt.DeviceID,
					Point:   pt.ID,
					Channel: p.DeviceMap[pt.ChannelID*10000+pt.DeviceID].Channel,
				}
				*p.Bus <- msg
				break
			}
		}
	}
}

//executeYk 执行遥控命令
func (p *Potocol104) executeYk(f *Frame) {
	addr := make([]byte, 0)
	addr = append(addr, 0)
	addr = append(addr, f.data[2])
	addr = append(addr, f.data[1])
	addr = append(addr, f.data[0])
	tmp := int(byteToint32(addr)) - 0x6000
	for _, pt := range p.ChannelConfig.Uptable.AddrYk {
		if tmp == pt.AcceptID {
			if f.kind == 0x2e {
				if (f.data[3] == 0x85 || f.data[3] == 0x86) && f.reason[0] == 0x07 {
					msg := command.Message{
						Cmd:     command.YK,
						SubCmd:  4,
						Value:   f.data[3] == 0x86,
						Device:  pt.DeviceID,
						Point:   pt.ID,
						Channel: p.DeviceMap[pt.ChannelID*10000+pt.DeviceID].Channel,
					}
					*p.Bus <- msg
					break
				}
				if (f.data[3] == 0x05 || f.data[3] == 0x06) && f.reason[0] == 0x07 {
					msg := command.Message{
						Cmd:     command.YK,
						SubCmd:  5,
						Value:   f.data[3] == 0x06,
						Device:  pt.DeviceID,
						Point:   pt.ID,
						Channel: p.DeviceMap[pt.ChannelID*10000+pt.DeviceID].Channel,
					}
					*p.Bus <- msg
					break
				}
			}
			if f.kind == 0x2d {
				if (f.data[3] == 0x80 || f.data[3] == 0x81) && f.reason[0] == 0x07 {
					msg := command.Message{
						Cmd:     command.YK,
						SubCmd:  4,
						Value:   f.data[3] != 0x81,
						Device:  pt.DeviceID,
						Point:   pt.ID,
						Channel: p.DeviceMap[pt.ChannelID*10000+pt.DeviceID].Channel,
					}
					*p.Bus <- msg
					break
				}
				if (f.data[3] == 0x00 || f.data[3] == 0x01) && f.reason[0] == 0x07 {
					msg := command.Message{
						Cmd:     command.YK,
						SubCmd:  5,
						Value:   f.data[3] != 0x01,
						Device:  pt.DeviceID,
						Point:   pt.ID,
						Channel: p.DeviceMap[pt.ChannelID*10000+pt.DeviceID].Channel,
					}
					*p.Bus <- msg
					break
				}
				break
			}

			break
		}
	}
}

//fillYx 全遥测数据上送
func (p *Potocol104) fillYc() []byte {
	data := make([]byte, 0)
	//填写数据
	if p.ChannelConfig.Uptable.AddrYc != nil {
		for _, pt := range p.ChannelConfig.Uptable.AddrYc {
			dev := p.DeviceMap[pt.ChannelID*10000+pt.DeviceID]
			point, _ := dev.GetYcByTime(pt.ID, p.txTime)
			data = append(data, addYcData(pt.AcceptID, point.Data)...)

		}
	}

	return data
}

//fillYx 填充变位遥测上送数据
func (p *Potocol104) fillChangeYc(mark bool) []byte {
	data := make([]byte, 0)
	//填写数据
	if p.ChannelConfig.Uptable.AddrYc != nil {
		for _, pt := range p.ChannelConfig.Uptable.AddrYc {
			point, need := p.DeviceMap[pt.ChannelID*10000+pt.DeviceID].GetYcByTime(pt.ID, p.txTime)
			if need == false && mark == true {
				fmt.Println("++++++++++++++++++++++++++++")
				continue
			}
			data = append(data, addChangeYcData(pt.AcceptID, point.Data)...)
		}
	}

	return data
}

//添加变位遥测信息体地址和数据
func addChangeYcData(AcceptID int, pdata float32) []byte {
	data := make([]byte, 0)
	bufAddr := byteSort(AcceptID + 0x4000)
	bufData := float32Tobyte(pdata)
	data = append(data, bufAddr[0], bufAddr[1], bufAddr[3])
	data = append(data, bufData[0], bufData[1], bufData[2], bufData[3], 0)
	return data
}

//添加遥测信息体地址和数据
func addYcData(AcceptID int, pdata float32) []byte {
	data := make([]byte, 0)
	bufAddr := byteSort(AcceptID + 0x4000)
	bufData := float32Tobyte(pdata)
	data = append(data, bufAddr[0], bufAddr[1], bufAddr[3])
	data = append(data, bufData[0], bufData[1], 0)
	return data
}

//addYcValue添加遥测值数据
func addYcValue(AcceptID int, pdata float32) []byte {
	data := make([]byte, 0)
	bufData := float32Tobyte(pdata)
	data = append(data, bufData[0], bufData[1], 0)
	return data
}

//添加遥信信息体地址和数据
func addYxData(ptAcceptID int, pdata float32) []byte {
	data := make([]byte, 0)
	buf := byteSort(ptAcceptID)
	data = append(data, buf[:3]...)
	if pdata > 0 {
		data = append(data, 0x01)
	} else {
		data = append(data, 0x00)
	}

	return data
}

//添加遥脉信息体地址和数据
func addYmData(ptAcceptID int, pdata float32) []byte {
	data := make([]byte, 0)
	buf := byteSort(ptAcceptID + 0x6400)
	data = append(data, buf[:3]...)
	data = append(data, float32Tobyte(pdata)...)
	data = append(data, 0)
	return data

}

//byteSort 字节排序
func byteSort(i int) []byte {
	n := int32(i)
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.LittleEndian, n)
	return buf.Bytes()
}

//write 发送报文
func (p *Potocol104) write(b bool, f Frame) {

	f.updateLength(b)
	tmp := string(f.getFrame(b))
	// fmt.Println("tmp------", f.getFrame(b))
	err := p.Port.Write(tmp, len(tmp))
	if err != nil {
		// p.ready = false
		// p.rxNo = 0
		// p.txNo = 0
	}
	time.Sleep(time.Millisecond)
}
