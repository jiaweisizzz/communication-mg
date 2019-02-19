package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
)

//ChannelConfig 通道配置结构体
type ChannelConfig struct {
	ID      int
	Name    string
	Crc     string
	Source  bool
	Port    PortConfig
	Potocol PotocolConfig
	Device  []DeviceConfig
	Uptable *Addrlist
}

type Form map[int]*List

type List struct {
	Type    string
	Addr    byte
	Funcode int
	Rnum    int
	Points  interface{}
}

// func (c *ChannelConfig) MakeAllList(t TemplateConfig) {
// 	tm := make(map[int]*List, 0)
// 	for _, dev := range c.Device {
// 		for index, n := range dev.Ycs {
//
// 		}
// 	}
// }

//GetChannelByName　通过通道名字获取通道
func (c *Config) GetChannelByName(name string) (*ChannelConfig, error) {
	for index, ch := range c.ChannelConfigList {
		if name == ch.Name {
			return &c.ChannelConfigList[index], nil
		}
	}
	return nil, errors.New("通道不存在")
}

//SetChannelByName 通过通道名设置通道
func (c *Config) SetPortChannelByName(name, crc, Pname, Pparitye string, Pbaud, Psize, Preadtimeout, Pstopbits int) error {
	ch, err := c.GetChannelByName(name)
	if err != nil {
		return errors.New("找不到通道999")
	}
	// b, err := strconv.Atoi(Pbaud)
	// if err != nil {
	// 	return errors.New("转换失败2")
	// }
	// d, err := strconv.Atoi(Psize)
	// if err != nil {
	// 	return errors.New("转换失败")
	// }
	// e, err := strconv.Atoi(Pstopbits)
	// if err != nil {
	// 	return errors.New("转换失败")
	// }
	ch.Crc = crc
	ch.Port.Name = c.ttyorcom(Pname)
	ch.Port.Baud = Pbaud
	ch.Port.Size = Psize
	ch.Port.ReadTimeout = Preadtimeout
	ch.Port.StopBits = Pstopbits
	ch.Port.Parity = Pparitye
	c.SaveChannel(*ch)
	return nil
}

func (c *Config) SetTcpChannelByName(name, pconntype, pip, portnum string, data interface{}) error {
	ch, err := c.GetChannelByName(name)
	if err != nil {
		return errors.New("找不到通道888")
	}
	ch.Port.Conntype = pconntype
	ch.Port.Ip = pip
	ch.Port.Portnum = portnum
	ch.Potocol.Data = data
	c.SaveChannel(*ch)
	return nil
}

//AddPortChannel 添加端口通道
func (c *Config) AddPortChannel(name, crc, ptype, Pname, Pparity, channelPotocol string, id, Pbaud, Psize, Preadtimeout, Pstopbits int, source bool) error {
	ch := new(ChannelConfig)
	ch.Name = name
	ch.ID = id
	ch.Crc = crc
	ch.Source = source
	ch.Potocol.Name = channelPotocol
	ch.Port.Ptype = ptype
	ch.Port.Name = c.ttyorcom(Pname)
	ch.Port.Baud = Pbaud
	ch.Port.Size = Psize
	ch.Port.ReadTimeout = Preadtimeout
	ch.Port.StopBits = Pstopbits
	ch.Port.Parity = Pparity
	c.ChannelConfigList = append(c.ChannelConfigList, *ch)
	c.SaveChannel(*ch)
	return nil
}

func (c *Config) ttyorcom(pname string) string {
	if runtime.GOOS == "linux" {
		switch pname {
		case "COM1":
			return "/dev/ttyS0"
		case "COM2":
			return "/dev/ttyS1"
		case "COM3":
			return "/dev/ttyS2"
		case "COM4":
			return "/dev/ttyS3"
		case "COM5":
			return "/dev/ttyS4"
		case "COM6":
			return "/dev/ttyS5"
		case "COM7":
			return "/dev/ttyS6"
		case "COM8":
			return "/dev/ttyS7"
		case "COM9":
			return "/dev/ttyS8"
		case "COM10":
			return "/dev/ttyS9"
		case "COM11":
			return "/dev/ttyS10"
		case "COM12":
			return "/dev/ttyS11"
		case "COM13":
			return "/dev/ttyS12"
		case "COM14":
			return "/dev/ttyS13"
		case "COM15":
			return "/dev/ttyS14"
		case "COM16":
			return "/dev/ttyS15"
		case "COM17":
			return "/dev/ttyS16"
		case "COM18":
			return "/dev/ttyS17"
		}
	}
	if runtime.GOOS == "windows" {
		return pname
	}
	return ""
}

//AddTcpChannel 添加tcp通道
func (c *Config) AddTcpChannel(name, ptype, channelPotocol, ctype, pip, portnum string, id int, source bool, data interface{}) error {
	ch := new(ChannelConfig)
	ch.Name = name
	ch.ID = id
	ch.Source = source
	ch.Potocol.Name = channelPotocol
	ch.Potocol.Data = data
	ch.Port.Ptype = ptype
	ch.Port.Conntype = ctype
	ch.Port.Ip = pip
	ch.Port.Portnum = portnum
	c.ChannelConfigList = append(c.ChannelConfigList, *ch)
	c.SaveChannel(*ch)
	return nil
}

//DeleteChannel 删除通道
func (c *Config) DeleteChannel(name string) error {
	var exist bool
	exist = false
	ch, err := c.GetChannelByName(name)
	if err != nil {
		return errors.New("找不到设备")
	}
	for index, channel := range c.ChannelConfigList {
		if channel.Name == ch.Name {
			exist = true
			c.ChannelConfigList = append(c.ChannelConfigList[:index], c.ChannelConfigList[index+1:]...)
			err := channel.Delete()
			if err != nil {
				return errors.New("删除文件失败")
			}
		}
	}
	if exist == false {
		return errors.New("设备不存在，删除失败")
	}
	return nil
}

//loadChannels 加载所有通道
func loadChannels(c *[]ChannelConfig) error {
	fs, err := ListDir("channels", "json")
	if err != nil {
		return errors.New("目录加载错误")
	}
	for _, file := range fs {
		if err := loadChannel(file, c); err != nil {
			return err
		}
	}
	fmt.Println("通道配置加载完成")
	return nil
}

//loadChannel 加载通道
func loadChannel(file string, c *[]ChannelConfig) error {
	f, err := ioutil.ReadFile(file)
	if err != nil {
		return errors.New("读取通道配置文件错误")
	}
	cc := new(ChannelConfig)
	err = json.Unmarshal(f, cc)
	if err != nil {
		return errors.New("Unmarshal error2211")
	}

	*c = append(*c, *cc)
	return nil
}

//SaveChannel 保存通道
func (c *Config) SaveChannel(ch ChannelConfig) {
	ch.Save()
}

//Save 保存通道
func (c *ChannelConfig) Save() error {
	b, err := json.MarshalIndent(*c, "", "	")
	if err != nil {
		return errors.New("Marshal error")
	}
	path := "channels/" + c.Name + ".json"
	ioutil.WriteFile(path, b, os.ModePerm)

	return nil
}

//Delete 删除通道
func (c *ChannelConfig) Delete() error {
	path := "channels/" + c.Name + ".json"
	err := os.Remove(path)
	if err != nil {
		return errors.New("删除失败")
	}
	return nil
}

//GetPotocolConfig 获得规约配置
func (c *ChannelConfig) GetPotocolConfig() PotocolConfig {
	return c.Potocol
}
