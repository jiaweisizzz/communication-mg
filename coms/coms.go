package coms

import (
	"errors"
	"fmt"
	"mg-coms/channel"
	"mg-coms/command"
	"mg-coms/config"
	"mg-coms/device"
	"mg-coms/porter"
	"mg-coms/process"
	"strconv"
)

//Coms 通讯管理机结构
type Coms struct {
	DeviceMap  device.DeviceMap   //设备表map[int]string
	ChannelMap channel.ChannelMap //通道表map[int]string
	Commander  command.Command    // 消息中心
	Config     config.Config      //配置
	PorterMap  porter.PorterMap
	Process    process.Process
}

func (c *Coms) GetProcess() (temp, data float64) {
	sum, num := c.Process.ProcessGet()
	return sum, num
}

//GetChannelByName　通过通道名字获取通道
func (c *Coms) GetChannelByName(name string) (*channel.Channel, error) {
	for _, ch := range c.ChannelMap {
		if name == ch.Name {
			return ch, nil
		}
	}
	return nil, errors.New("通道不存在")
}

func (c *Coms) GetAllDevices() device.DeviceMap {
	return c.DeviceMap
}

//GetDeviceByID　通过设备ID获取设备
func (c *Coms) GetDeviceByID(name, id string) (*device.Device, error) {
	ch, err := c.GetChannelByName(name)
	if err != nil {
		return nil, err
	}
	i, _ := strconv.Atoi(id)
	for index, dev := range ch.DeviceMap {
		if ch.ID*10000+i == index {
			return dev, nil
		}
	}
	return nil, errors.New("设备不存在")
}

//Start 启动通讯管理机
func (c *Coms) Start() {
	c.PorterMap.OpenAllPorter()
	c.StartAllChannel()
	c.Commander.Start()

}

//Stop 停止通讯管理机
func (c *Coms) Stop() {
	c.StopAllChannel()
	c.CloseAllPorter()
}

//StartAllChannel 启动全部通道
func (c *Coms) StartAllChannel() {
	for _, ch := range c.ChannelMap {
		c.StartChannel(ch)
	}
	println("通讯管理机已启动")
}

//StartChannel 启动通道
func (c *Coms) StartChannel(ch *channel.Channel) {
	ch.Start()
	c.Commander.Register(ch.Name, ch.Handle)
}

//StopAllChannel 停止全部通道
func (c *Coms) StopAllChannel() {
	for index, ch := range c.ChannelMap {
		c.StopChannel(ch)
		delete(c.ChannelMap, index)
	}
	println("通讯管理机已停止")
}

//StopAllChannel 停止全部通道
func (c *Coms) CloseAllPorter() {
	for index, _ := range c.PorterMap {
		c.PorterMap[index].Close()
		delete(c.PorterMap, index)
	}
	fmt.Println("suoyouduankouyitingzhi!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
}

//StopChannel 停止通道
func (c *Coms) StopChannel(ch *channel.Channel) {
	c.Commander.UnRegister(ch.Name)
	ch.Stop()
}

//Init 初始化通讯管理机实例
func (c *Coms) Init() (*Coms, error) {

	c.Commander.Init()
	//通讯管理机的配置
	err := c.Config.Load()
	if err != nil {
		return nil, err
	}
	err = c.Config.Creataddrlist()
	if err != nil {
		return nil, err
	}

	p, err := c.Config.GetPortConfig()

	if err != nil {
		// return nil, err
		fmt.Println(err)
	}
	c.PorterMap = porter.New(p)

	dcl, err := c.Config.GetDeviceConfig()
	if err != nil {
		// return nil, err
		fmt.Println("-----", err)
	}
	// fmt.Println(dcl)
	c.DeviceMap = device.New(dcl)
	// fmt.Println("+++++++++++++", c.DeviceMap)
	// c.DeviceMap.Print()
	c.Commander.SetDeviceMap(c.DeviceMap)

	cc, err := c.Config.GetChannelConfig()
	if err != nil {
		// return nil, err
		fmt.Println(err)
	}
	c.ChannelMap = channel.New(cc)

	c.ChannelMap.SetAllChannelDeviceMap(c.DeviceMap)
	c.ChannelMap.SetAllChannelBus(&c.Commander.Bus)
	c.ChannelMap.SetAllChannelPorter(c.PorterMap)
	// c.ChannelMap.StartAll()
	return c, nil
}

// //TestSend 发送消息
// func (c *Coms) TestSend() {
// 	c.Commander.Send(command.Message{Cmd: command.UPDATEDATA, Device: "设备02"})
// }

//New 新建通讯管理机实例
func New() (*Coms, error) { return new(Coms).Init() }
