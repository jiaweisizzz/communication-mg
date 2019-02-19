package channel

import (
	"fmt"
	"mg-coms/command"
	"mg-coms/config"
	"mg-coms/device"
	"mg-coms/porter"
	"mg-coms/potocol"
	"time"
)

//ChannelRun 真
//ChannelStop 假
const (
	ChannelRun  = true
	ChannelStop = false
)

//ChannelMap 通道列表
type ChannelMap map[string]*Channel

//Channel 通道配置结构体
type Channel struct {
	ID        int    //通道ID
	Name      string //通道名称\
	Crc       string
	Source    bool
	Running   bool                //运行状态
	Handle    *command.Handle     //当前通道消息句柄
	Bus       *command.Bus        //通讯管理机消息入口
	Port      *porter.Porter      //端口
	Potocol   potocol.PotocolFunc //规约
	DeviceMap device.DeviceMap    //设备
}

// //Send 通道消息入口
// func (c *ChannelMap) Send() {
//
// }

func (cl ChannelMap) SetCallBack(id int, t time.Time, f func(string, string)) {
	fmt.Println("ChannelMap", cl)
	for index, ch := range cl {
		if id == ch.ID {
			cl[index].Port.SetCallBack(t, f)
		}
	}
}

func (cl ChannelMap) RemoveCallBack(id int, t time.Time) {
	for index, ch := range cl {
		if id == ch.ID {
			cl[index].Port.RemoveCallBack(t)
		}
	}
}

//Init 通道初始化
func (cl ChannelMap) Init(list []config.ChannelConfig) ChannelMap {
	if cl == nil {
		cl = make(map[string]*Channel, 0)
	}
	for _, cc := range list {
		cl[cc.Name] = new(Channel).init(cc)
	}
	return cl
}

func (c *Channel) init(cc config.ChannelConfig) *Channel {
	c.Handle = command.NewHandle()
	c.ID = cc.ID
	c.Name = cc.Name
	c.Source = cc.Source
	pot := cc.GetPotocolConfig()
	c.Potocol = potocol.New(pot.Name, cc, c.DeviceMap)

	return c
}

//New 新建通道实例
func New(list []config.ChannelConfig) ChannelMap {
	return new(ChannelMap).Init(list)
}

//SetAllChannelDeviceMap 设置所有通道列表
func (cl *ChannelMap) SetAllChannelDeviceMap(dm device.DeviceMap) {
	for _, ch := range *cl {
		ch.SetDeviceMap(dm)
	}
}

//SetAllChannelBus 设置所有通道消息入口
func (cl *ChannelMap) SetAllChannelBus(b *command.Bus) {
	for _, ch := range *cl {
		ch.SetBus(b)
	}
}

//SetAllChannelPorter 设置所有通道端口
func (cl *ChannelMap) SetAllChannelPorter(pm porter.PorterMap) {
	for _, ch := range *cl {
		ch.SetPorter(pm[ch.Name])
	}
}

//SetPorter 设置通讯端口
func (c *Channel) SetPorter(p *porter.Porter) {
	c.Port = p
}

//SetDeviceMap 设置设备列表
func (c *Channel) SetDeviceMap(dm device.DeviceMap) {
	c.DeviceMap = dm
	// fmt.Println(c.DeviceMap)
}

//SetBus 设置消息入口
func (c *Channel) SetBus(b *command.Bus) {
	c.Bus = b
}

//Start 启动通道
func (c *Channel) Start() {
	// fmt.Println("cccccccc", c.Bus, c.Port, c.Potocol, c.DeviceMap)
	c.Potocol.Config(c.Bus, c.Port, c.Potocol, c.DeviceMap)
	fmt.Println(c.Name, "启动...")
	c.Potocol.Created()
	go c.start()
	// c.Potocol.Created()
}

//Stop 停止通道
func (c *Channel) Stop() {
	c.stop()
}

func (c *Channel) start() {
	if c.Source == false {
		time.Sleep(4 * time.Second)
	}
	// defer fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
	c.Running = ChannelRun
	for {
		// fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!", c.Name)
		if !c.Running {
			fmt.Println(c.Name, "停止...")
			// c.Port.Close()
			return
		}
		select {
		case msg := <-*c.Handle:
			c.Potocol.Cmd(msg)
		default:
			startime := time.Now()
			c.Potocol.Do()
			overtime := time.Now()
			if overtime.Before(startime.Add(200 * time.Millisecond)) {
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
}

func (c *Channel) stop() {
	c.Running = ChannelStop
}
