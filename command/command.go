package command

import (
	"fmt"
	"mg-coms/device"
)

var (
	commander Command
)

//Bus Bus
type Bus chan Message

//Handle 消息
type Handle chan Message

//Command 控制中心
type Command struct {
	Bus         Bus
	RegisterMap map[string]*Handle
	DeviceMap   device.DeviceMap
}

//Register 注册列表
func (cmd *Command) Register(channelName string, handle *Handle) {
	if cmd.RegisterMap == nil {
		cmd.RegisterMap = make(map[string]*Handle)
		fmt.Println("初始化注册列表")
	}
	cmd.RegisterMap[channelName] = handle
}

//UnRegister  删除列表
func (cmd *Command) UnRegister(channelName string) {
	delete(cmd.RegisterMap, channelName)
}

//Send 向消息中心发送命令
func (cmd *Command) Send(msg Message) {
	fmt.Println("开始向所选通道发送消息 [ 类型:", msg.Cmd, "子类型:", msg.SubCmd, "通道名称:", msg.Channel, "设备号:", msg.Device, "点号:", msg.Point, "目标值:", msg.Value)
	// for name, handle := range cmd.RegisterMap {
	// 	fmt.Println("向通道", name, "发送消息", msg)
	// 	*handle <- msg
	// }

	// for _, dev := range *cmd.DeviceMap {
	// 	if dev. == msg.Device {
	// 		cmd.exeCmd(msg)
	// 	}
	// }
	cmd.exeCmd(msg)
}

//Start 控制中心初始化
func (cmd Command) Start() {
	go cmd.start()
}

func (cmd *Command) start() {
	fmt.Println("启动消息中心，当前注册通道为", cmd.RegisterMap)
	for {
		select {
		case msg := <-cmd.Bus:
			fmt.Println("通讯管理机收到消息->", msg)
			cmd.exeCmd(msg)
		}
	}
}

func (cmd *Command) exeCmd(msg Message) {
	// c := msg.Cmd
	fmt.Println("通讯管理机开始执行消息->", msg)
	// switch c {
	// case UPDATEDATA:
	// 	cmd.updateData(msg)
	// default:
	cmd.SendChannel(msg)
	// }
}

// func (cmd *Command) updateData(msg Message) error {
// 	dev, ok := cmd.DeviceMap[msg.Device]
// 	if !ok {
// 		fmt.Println("发送结束", cmd.DeviceMap)
//
// 		return errors.New("未找到设备")
// 	}
// 	handle, ok := cmd.RegisterMap[dev.Channel]
// 	if !ok {
// 		return errors.New("未找到通道")
// 	}
// 	fmt.Println("发送心跳结束")
// 	*handle <- msg
// 	return nil
// }

//SendDevice 向设备发送数据
func (cmd *Command) SendDevice(msg Message) {

}

//SendChannel 向设备发送数据
func (cmd *Command) SendChannel(msg Message) {
	if msg.Channel == "" {
		for _, dev := range cmd.DeviceMap {
			if dev.Tnumber == msg.Device {
				msg.Channel = dev.Channel
			}
		}
	}
	*cmd.RegisterMap[msg.Channel] <- msg
}

//SetDeviceMap 设置设备表
func (cmd *Command) SetDeviceMap(dm device.DeviceMap) {
	cmd.DeviceMap = dm
}

//Init 消息中心初始化
func (cmd *Command) Init() {
	if cmd.Bus == nil {
		cmd.Bus = make(chan Message, 100)
	}
}

//NewHandle 创建Handle实例
func NewHandle() *Handle {
	return new(Handle).Init()
}

//Init 初始化Handle
func (h Handle) Init() *Handle {
	h = make(chan Message, 100)
	return &h
}
