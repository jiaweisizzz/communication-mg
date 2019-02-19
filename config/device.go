package config

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

//DeviceConfig 设备配置结构体
type DeviceConfig struct {
	ID        int
	Tnumber   int //终端号
	Name      string
	Time      time.Time
	Channel   string
	ChannelID int
	Status    bool
	Ycs       []Point //遥测点表
	Yxs       []Point //遥信点表
	Yms       []Point //遥脉点表
	Yks       []Point //遥控点表
	Yts       []Point //遥调点表
}

type Point struct {
	ID    int
	Name  string
	Time  time.Time
	Scale float32
	Data  interface{}
}

func (p *Point) GetIntData(s string) int {
	i, _ := strconv.Atoi(p.Data.(map[string]interface{})[s].(string))
	return i
}

func (p *Point) GetUintData(s string) uint64 {
	i, _ := strconv.ParseUint(p.Data.(map[string]interface{})[s].(string), 0, 0)
	return i
}

//
// //Yc 遥测点表
// type Yc struct {
// 	Point
// }
//
// //Yx 遥信点表
// type Yx struct {
// 	Point
// }
//
// //Ym 遥脉点表
// type Ym struct {
// 	Point
// }
//
// //Yk 遥控点表
// type Yk struct {
// 	Point
// }
//
// //Yt 遥调点表
// type Yt struct {
// 	Point
// }

//GetDeviceByID　通过设备ID获取设备
func (c *Config) GetDeviceByID(name, id string) (*DeviceConfig, error) {
	ch, err := c.GetChannelByName(name)
	if err != nil {
		return nil, err
	}
	for index, dev := range ch.Device {
		if id == strconv.Itoa(dev.ID) {
			return &ch.Device[index], nil
		}
	}
	return nil, errors.New("设备不存在")
}

//GetChannelDevice 获得设备信息
func (c *Config) GetChannelDevice(name string) (Device []DeviceConfig, err error) {
	ch, err := c.GetChannelByName(name)
	if err != nil {
		return nil, err
	}
	return ch.Device, nil
}

//SetDeviceByID 通过设备ID设置设备
func (c *Config) SetDeviceByID(name, dname, template string, id, dTnumber int) (*DeviceConfig, error) {
	tp, err := c.GetTemplateByName(template)
	if err != nil {
		return nil, err
	}
	ch, err := c.GetChannelByName(name)
	fmt.Println(name)
	if err != nil {
		return nil, errors.New("找不到通道777")
	}
	for index, dev := range ch.Device {
		if dname == dev.Name {
			ch.Device[index].ID = id
			ch.Device[index].Tnumber = dTnumber
			ch.Device[index].Ycs = tp.Ycs
			ch.Device[index].Yxs = tp.Yxs
			ch.Device[index].Yms = tp.Yms
			ch.Device[index].Yks = tp.Yks
			ch.Device[index].Yts = tp.Yts
			c.SaveChannel(*ch)
			return &ch.Device[index], nil
		}
	}
	// dev.Name = value
	// c.SaveChannel(*ch)

	// b, err := json.Marshal(c.ChannelConfigList)
	// if err != nil {
	// 	return nil, errors.New("Marshal error")
	// }
	// ioutil.WriteFile("channel01.json", b, os.ModeAppend)
	return nil, errors.New("设备不存在")
}

func initDevice() *DeviceConfig {
	dev := new(DeviceConfig)
	dev.Ycs = make([]Point, 0, 100)
	dev.Yks = make([]Point, 0, 100)
	dev.Yms = make([]Point, 0, 100)
	dev.Yts = make([]Point, 0, 100)
	dev.Yxs = make([]Point, 0, 100)
	return dev
}

//AddDevice　增加设备
func (c *Config) AddDevice(channelname, name, template string, id, channelid int) error {
	tp, err := c.GetTemplateByName(template)
	if err != nil {
		return err
	}
	dev := initDevice()
	dev.Name = name
	dev.Channel = channelname
	dev.ChannelID = channelid
	// i, err := strconv.Atoi(id)
	// if err != nil {
	// 	return errors.New("atoi出错")
	// }
	dev.ID = id
	tnumber, _ := c.GetDeviceNum()
	dev.Tnumber = tnumber + 1
	dev.Status = true
	dev.Ycs = tp.Ycs
	dev.Yxs = tp.Yxs
	dev.Yms = tp.Yms
	dev.Yks = tp.Yks
	dev.Yts = tp.Yts
	for index, ch := range c.ChannelConfigList {
		if ch.Name == channelname {
			c.ChannelConfigList[index].Device = append(c.ChannelConfigList[index].Device, *dev)
			c.SaveChannel(c.ChannelConfigList[index])
			break
		}
	}
	// *d = append(*d, *dev)
	// b, err := json.Marshal(ch)
	// if err != nil {
	// 	return errors.New("Marshal error")
	// }
	// ioutil.WriteFile("channels/channel01.json", b, os.ModeAppend)
	return nil
}

//DeleteDevice 删除设备
func (c *Config) DeleteDevice(name, id string) error {
	var exist bool
	exist = false
	dev, err := c.GetDeviceByID(name, id)
	if err != nil {
		return errors.New("找不到设备2")
	}
	ch, err := c.GetChannelByName(name)
	if err != nil {
		return errors.New("找不到通道2")
	}
	// for _, ch := range c.ChannelConfigList {
	for index, device := range ch.Device {
		if device.ID == dev.ID {
			exist = true
			ch.Device = append(ch.Device[:index], ch.Device[index+1:]...)
			c.SaveChannel(*ch)
			break
		}
		// }
	}
	if exist == false {
		return errors.New("设备不存在，删除失败")
	}
	return nil
}

//SetYcsPointByID 通过id设置遥测点
func (c *Config) SetYcsPointByID(channelname, deviceid, pointid, value string) (*DeviceConfig, error) {
	ch, err := c.GetChannelByName(channelname)
	if err != nil {
		return nil, errors.New("配置ycs找不到通道")
	}
	f, err := strconv.ParseFloat(value, 32)
	if err != nil {
		return nil, errors.New("转换float出错")
	}
	for index, dev := range ch.Device {
		if deviceid == strconv.Itoa(dev.ID) {
			for i, yc := range dev.Ycs {
				if pointid == strconv.Itoa(yc.ID) {
					yc.Scale = float32(f)
					dev.Ycs[i] = yc
					ch.Device[index] = dev
					c.SaveChannel(*ch)
					return &dev, nil
				}
			}
		}
	}
	return nil, errors.New("设置ycs失败")
}

//SetYxsPointByID 通过id设置遥信点
func (c *Config) SetYxsPointByID(channelname, deviceid, pointid, value string) (*DeviceConfig, error) {
	ch, err := c.GetChannelByName(channelname)
	if err != nil {
		return nil, errors.New("配置ycs找不到通道")
	}
	f, err := strconv.ParseFloat(value, 32)
	if err != nil {
		return nil, errors.New("转换float出错")
	}
	for index, dev := range ch.Device {
		if deviceid == strconv.Itoa(dev.ID) {
			for i, yx := range dev.Yxs {
				if pointid == strconv.Itoa(yx.ID) {
					yx.Scale = float32(f)
					dev.Yxs[i] = yx
					ch.Device[index] = dev
					c.SaveChannel(*ch)
					return &dev, nil
				}
			}
		}
	}
	return nil, errors.New("设置ycs失败")
}

//SetYmsPointByID 通过id设置遥脉点
func (c *Config) SetYmsPointByID(channelname, deviceid, pointid, value string) (*DeviceConfig, error) {
	ch, err := c.GetChannelByName(channelname)
	if err != nil {
		return nil, errors.New("配置ycs找不到通道")
	}
	f, err := strconv.ParseFloat(value, 32)
	if err != nil {
		return nil, errors.New("转换float出错")
	}
	for index, dev := range ch.Device {
		if deviceid == strconv.Itoa(dev.ID) {
			for i, ym := range dev.Yms {
				if pointid == strconv.Itoa(ym.ID) {
					ym.Scale = float32(f)
					dev.Yms[i] = ym
					ch.Device[index] = dev
					c.SaveChannel(*ch)
					return &dev, nil
				}
			}
		}
	}
	return nil, errors.New("设置ycs失败")
}

//SetYksPointByID 通过id设置遥控点
func (c *Config) SetYksPointByID(channelname, deviceid, pointid, value string) (*DeviceConfig, error) {
	ch, err := c.GetChannelByName(channelname)
	if err != nil {
		return nil, errors.New("配置ycs找不到通道")
	}
	f, err := strconv.ParseFloat(value, 32)
	if err != nil {
		return nil, errors.New("转换float出错")
	}
	for index, dev := range ch.Device {
		if deviceid == strconv.Itoa(dev.ID) {
			for i, yk := range dev.Yks {
				if pointid == strconv.Itoa(yk.ID) {
					yk.Scale = float32(f)
					dev.Yks[i] = yk
					ch.Device[index] = dev
					c.SaveChannel(*ch)
					return &dev, nil
				}
			}
		}
	}
	return nil, errors.New("设置ycs失败")
}

//SetYtsPointByID 通过id设置遥调点
func (c *Config) SetYtsPointByID(channelname, deviceid, pointid, value string) (*DeviceConfig, error) {
	ch, err := c.GetChannelByName(channelname)
	if err != nil {
		return nil, errors.New("配置ycs找不到通道")
	}
	f, err := strconv.ParseFloat(value, 32)
	if err != nil {
		return nil, errors.New("转换float出错")
	}
	for index, dev := range ch.Device {
		if deviceid == strconv.Itoa(dev.ID) {
			for i, yt := range dev.Yts {
				if pointid == strconv.Itoa(yt.ID) {
					yt.Scale = float32(f)
					dev.Yts[i] = yt
					ch.Device[index] = dev
					c.SaveChannel(*ch)
					return &dev, nil
				}
			}
		}
	}
	return nil, errors.New("设置ycs失败")
}

//AddYcsPoint 增加遥测点
func (c *Config) AddYcsPoint(channelname, deviceid, pname, pid string) error {
	yc := new(Point)
	yc.Name = pname
	i, err := strconv.Atoi(pid)
	if err != nil {
		return errors.New("atoi出错")
	}
	yc.ID = i
	for index, ch := range c.ChannelConfigList {
		if ch.Name == channelname {
			for _, dev := range ch.Device {
				if deviceid == strconv.Itoa(dev.ID) {
					dev.Ycs = append(dev.Ycs, *yc)
					ch.Device[index] = dev
					c.SaveChannel(ch)
					break
				}
			}
		}
	}
	return nil
}

//AddYxsPoint 增加遥信点
func (c *Config) AddYxsPoint(channelname, deviceid, pname, pid string) error {
	yx := new(Point)
	yx.Name = pname
	i, err := strconv.Atoi(pid)
	if err != nil {
		return errors.New("atoi出错")
	}
	yx.ID = i
	for index, ch := range c.ChannelConfigList {
		if ch.Name == channelname {
			for _, dev := range ch.Device {
				if deviceid == strconv.Itoa(dev.ID) {
					dev.Yxs = append(dev.Yxs, *yx)
					ch.Device[index] = dev
					c.SaveChannel(ch)
					break
				}
			}
		}
	}
	return nil
}

//AddYmsPoint 增加遥脉点
func (c *Config) AddYmsPoint(channelname, deviceid, pname, pid string) error {
	ym := new(Point)
	ym.Name = pname
	i, err := strconv.Atoi(pid)
	if err != nil {
		return errors.New("atoi出错")
	}
	ym.ID = i
	for index, ch := range c.ChannelConfigList {
		if ch.Name == channelname {
			for _, dev := range ch.Device {
				if deviceid == strconv.Itoa(dev.ID) {
					dev.Yms = append(dev.Yms, *ym)
					ch.Device[index] = dev
					c.SaveChannel(ch)
					break
				}
			}
		}
	}
	return nil
}

//AddYksPoint 增加遥控点
func (c *Config) AddYksPoint(channelname, deviceid, pname, pid string) error {
	yk := new(Point)
	yk.Name = pname
	i, err := strconv.Atoi(pid)
	if err != nil {
		return errors.New("atoi出错")
	}
	yk.ID = i
	for index, ch := range c.ChannelConfigList {
		if ch.Name == channelname {
			for _, dev := range ch.Device {
				if deviceid == strconv.Itoa(dev.ID) {
					dev.Yks = append(dev.Yks, *yk)
					ch.Device[index] = dev
					c.SaveChannel(ch)
					break
				}
			}
		}
	}
	return nil
}

//AddYtsPoint 增加遥调点
func (c *Config) AddYtsPoint(channelname, deviceid, pname, pid string) error {
	yt := new(Point)
	yt.Name = pname
	i, err := strconv.Atoi(pid)
	if err != nil {
		return errors.New("atoi出错")
	}
	yt.ID = i
	for index, ch := range c.ChannelConfigList {
		if ch.Name == channelname {
			for _, dev := range ch.Device {
				if deviceid == strconv.Itoa(dev.ID) {
					dev.Yts = append(dev.Yts, *yt)
					ch.Device[index] = dev
					c.SaveChannel(ch)
					break
				}
			}
		}
	}
	return nil
}

//DeleteYcsPoint 删除遥测点
func (c *Config) DeleteYcsPoint(channelname, deviceid, pointid string) error {
	var exist bool
	exist = false
	// dev, err := c.GetDeviceByID(channelname, deviceid)
	// if err != nil {
	// 	return errors.New("找不到设备")
	// }
	ch, err := c.GetChannelByName(channelname)
	if err != nil {
		return errors.New("找不到通道")
	}
	for i, device := range ch.Device {
		if deviceid == strconv.Itoa(device.ID) {
			for index, yc := range device.Ycs {
				if pointid == strconv.Itoa(yc.ID) {
					exist = true
					device.Ycs = append(device.Ycs[:index], device.Ycs[index+1:]...)
					ch.Device[i] = device
					c.SaveChannel(*ch)
					break
				}
			}
		}
	}
	if exist == false {
		return errors.New("Ycs不存在，删除失败")
	}
	return nil
}

//DeleteYxsPoint 删除遥信点
func (c *Config) DeleteYxsPoint(channelname, deviceid, pointid string) error {
	var exist bool
	exist = false
	// dev, err := c.GetDeviceByID(channelname, deviceid)
	// if err != nil {
	// 	return errors.New("找不到设备")
	// }
	ch, err := c.GetChannelByName(channelname)
	if err != nil {
		return errors.New("找不到通道")
	}
	for i, device := range ch.Device {
		if deviceid == strconv.Itoa(device.ID) {
			for index, yx := range device.Yxs {
				if pointid == strconv.Itoa(yx.ID) {
					exist = true
					device.Yxs = append(device.Yxs[:index], device.Yxs[index+1:]...)
					ch.Device[i] = device
					c.SaveChannel(*ch)
					break
				}
			}
		}
	}
	if exist == false {
		return errors.New("Ycs不存在，删除失败")
	}
	return nil
}

//DeleteYmsPoint 删除遥脉点
func (c *Config) DeleteYmsPoint(channelname, deviceid, pointid string) error {
	var exist bool
	exist = false
	// dev, err := c.GetDeviceByID(channelname, deviceid)
	// if err != nil {
	// 	return errors.New("找不到设备")
	// }
	ch, err := c.GetChannelByName(channelname)
	if err != nil {
		return errors.New("找不到通道")
	}
	for i, device := range ch.Device {
		if deviceid == strconv.Itoa(device.ID) {
			for index, ym := range device.Yms {
				if pointid == strconv.Itoa(ym.ID) {
					exist = true
					device.Yms = append(device.Yms[:index], device.Yms[index+1:]...)
					ch.Device[i] = device
					c.SaveChannel(*ch)
					break
				}
			}
		}
	}
	if exist == false {
		return errors.New("Ycs不存在，删除失败")
	}
	return nil
}

//DeleteYksPoint 删除遥控点
func (c *Config) DeleteYksPoint(channelname, deviceid, pointid string) error {
	var exist bool
	exist = false
	// dev, err := c.GetDeviceByID(channelname, deviceid)
	// if err != nil {
	// 	return errors.New("找不到设备")
	// }
	ch, err := c.GetChannelByName(channelname)
	if err != nil {
		return errors.New("找不到通道")
	}
	for i, device := range ch.Device {
		if deviceid == strconv.Itoa(device.ID) {
			for index, yk := range device.Yks {
				if pointid == strconv.Itoa(yk.ID) {
					exist = true
					device.Yks = append(device.Yks[:index], device.Yks[index+1:]...)
					ch.Device[i] = device
					c.SaveChannel(*ch)
					break
				}
			}
		}
	}
	if exist == false {
		return errors.New("Ycs不存在，删除失败")
	}
	return nil
}

//DeleteYtsPoint 删除遥调点
func (c *Config) DeleteYtsPoint(channelname, deviceid, pointid string) error {
	var exist bool
	exist = false
	// dev, err := c.GetDeviceByID(channelname, deviceid)
	// if err != nil {
	// 	return errors.New("找不到设备")
	// }
	ch, err := c.GetChannelByName(channelname)
	if err != nil {
		return errors.New("找不到通道")
	}
	for i, device := range ch.Device {
		if deviceid == strconv.Itoa(device.ID) {
			for index, yt := range device.Yts {
				if pointid == strconv.Itoa(yt.ID) {
					exist = true
					device.Yts = append(device.Yts[:index], device.Yts[index+1:]...)
					ch.Device[i] = device
					c.SaveChannel(*ch)
					break
				}
			}
		}
	}
	if exist == false {
		return errors.New("Ycs不存在，删除失败")
	}
	return nil
}
