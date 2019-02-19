package potocol

import (
	"mg-coms/command"
	"mg-coms/config"
	"mg-coms/device"
	"mg-coms/porter"
	"time"
)

//PotocolParm 规约配置结构体
type PotocolParm struct {
	Name          string           //规约名
	Data          interface{}      //内容
	Bus           *command.Bus     //通讯管理机消息入口
	Port          *porter.Porter   //端口
	Potocol       PotocolFunc      //规约
	DeviceMap     device.DeviceMap //设备
	Buffer        []byte           //缓存
	Index         int              //缓存结束位置
	Last          time.Time
	ChannelConfig config.ChannelConfig
}

//Config 规约初始化
func (p *PotocolParm) Config(bus *command.Bus,
	port *porter.Porter,
	potocol PotocolFunc,
	deviceMap device.DeviceMap) {
	p.Buffer = make([]byte, 0, 1024)
	p.Bus = bus
	p.Port = port
	p.Potocol = potocol
	p.DeviceMap = deviceMap
}

//PotocolFunc 规约模板
type PotocolFunc interface {
	Cmd(msg command.Message)
	Do()
	Created()
	Config(
		bus *command.Bus,
		port *porter.Porter,
		potocol PotocolFunc,
		deviceMap device.DeviceMap,
	)
}

//New 新建通讯管理机实例
func New(name string, cc config.ChannelConfig, dm device.DeviceMap) PotocolFunc {
	var pot PotocolFunc
	switch name {
	case "MODBUS":
		pot = new(PotocolNewModbus).Init(cc, dm)
	case "MODBUSTCP":
		pot = new(PotocolModbusTcp).Init(cc, dm)
	case "JSON":
		pot = new(PotocolJson).Init(cc, dm)
	case "FORWARD":
		pot = new(PotocolForward).Init(cc, dm)
	case "104":
		pot = new(Potocol104).Init(cc, dm)
	case "645":
		pot = new(Potocol645).Init(cc, dm)
	}
	return pot
}
