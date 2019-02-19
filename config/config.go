package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

//Config 配置
type Config struct {
	Comsconfig        ComsConfig
	ChannelConfigList []ChannelConfig
	UsersConfigList   []UserConfig
	TemplateList      []TemplateConfig
	VersionConfigList []VersionConfig
	AddrList          Addrlist
	// UsersConfig Users
}

//PortConfig 端口配置结构体
type PortConfig struct {
	Ptype       string
	Conntype    string
	Name        string
	Baud        int
	Size        int
	ReadTimeout int
	Parity      string
	StopBits    int
	Channel     string
	Ip          string
	Portnum     string
}

//PotocolConfig 规约配置结构体
type PotocolConfig struct {
	Name string
	Data interface{}
}

//DeviceConfigList 设备配置列表
type DeviceConfigList []DeviceConfig

//Save 保存
func (c *Config) Save() {
	b, err := json.Marshal(c)
	if err != nil {
		fmt.Println("Marshal error")
		return
	}
	ioutil.WriteFile("data_test.json", b, os.ModePerm)
}

//GetVersionConfig 获得版本配置
func (c *Config) GetVersionConfig() (VersionConfig []VersionConfig, err error) {
	vl := c.VersionConfigList
	if vl == nil {
		err = errors.New("版本不存在")
		return nil, err
	}
	return vl, nil
}

//GetUserConfig 获得用户配置
func (c *Config) GetUserConfig() (Userconfig []UserConfig, err error) {
	usr := c.UsersConfigList
	if usr == nil {
		err = errors.New("用户不存在")
		return nil, err
	}
	return usr, nil
}

//GetTemplateConfig 获得模板配置
func (c *Config) GetTemplateConfig() (TemplateConfig []TemplateConfig, err error) {
	tl := c.TemplateList
	if tl == nil {
		err = errors.New("模板配置加载异常")
		return nil, err
	}
	return tl, nil
}

//GetChannelConfig 获得通道设置
func (c *Config) GetChannelConfig() (Channelconfig []ChannelConfig, err error) {
	ch := c.ChannelConfigList
	// fmt.Println(ch)
	if ch == nil {
		err = errors.New("通道配置加载异常")
		return nil, err
	}
	return ch, nil
}

//GetDeviceConfig 获得设备配置
func (c *Config) GetDeviceConfig() (Device []DeviceConfig, err error) {
	chs, err := c.GetChannelConfig()
	if err != nil {
		return nil, err
	}
	var result []DeviceConfig
	for _, ch := range chs {
		if ch.Source == true {
			for _, dev := range ch.Device {
				dev.Channel = ch.Name
				dev.ChannelID = ch.ID
				result = append(result, dev)
			}
		}
	}
	return result, nil
}

func (c *Config) GetDeviceNum() (num int, err error) {
	chs, err := c.GetChannelConfig()
	if err != nil {
		return 0, err
	}
	var i = 0
	for _, ch := range chs {
		for range ch.Device {
			i++
		}
	}
	return i, nil
}

//SetDeviceConfig 设置设备设置
func (c *Config) SetDeviceConfig(Device []DeviceConfig) error {
	chs, err := c.GetChannelConfig()
	if err != nil {
		err = errors.New("通道不存在")
		return err
	}

	for _, ch := range chs {
		ch.Device = Device
	}
	return nil
}

//GetPortConfig 获得端口配置
func (c *Config) GetPortConfig() ([]PortConfig, error) {
	var pt []PortConfig
	chs, err := c.GetChannelConfig()
	if err != nil {
		return nil, errors.New("获取Port时通道不存在")
	}
	for _, ch := range chs {
		pc := ch.Port
		pc.Channel = ch.Name
		pt = append(pt, pc)
	}
	return pt, nil
}

//SetPortConfig 设置端口配置
func (c *Config) SetPortConfig(Port PortConfig) error {
	chs, err := c.GetChannelConfig()
	if err != nil {
		err = errors.New("通道不存在")
		return err
	}

	for _, ch := range chs {
		ch.Port = Port
	}
	return nil
}

//GetPotocolConfig 获取规约配置
func (c *Config) GetPotocolConfig() ([]PotocolConfig, error) {
	var poc []PotocolConfig
	chs, err := c.GetChannelConfig()
	if err != nil {
		fmt.Println("通道不存在")
		return poc, errors.New("通道不存在")
	}
	for _, ch := range chs {
		po := ch.Potocol
		poc = append(poc, po)
	}
	return poc, nil
}

//SetPotocolConfig 设置规约配置
func (c *Config) SetPotocolConfig(Potocol PotocolConfig) error {
	chs, err := c.GetChannelConfig()
	if err != nil {
		err = errors.New("通道不存在")
		return err
	}

	for _, ch := range chs {
		ch.Potocol = Potocol
	}
	return nil
}

//Load 初始化配置
func (c *Config) Load() error {

	if err := loadComs(&c.Comsconfig); err != nil {
		return err
	}

	if err := loadUsers(&c.UsersConfigList); err != nil {
		return err
	}

	if err := loadChannels(&c.ChannelConfigList); err != nil {
		return err
	}

	if err := loadTempaltes(&c.TemplateList); err != nil {
		return err
	}
	if err := loadVersion(&c.VersionConfigList); err != nil {
		return err
	}

	return nil
}

//MakeFile 创建文件夹
func (c *Config) MakeFile(name string) error {
	err := os.Mkdir(name, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

//Empty 清空配置
func (c *Config) Empty() {
	c.TemplateList = nil
	c.ChannelConfigList = nil
	c.UsersConfigList = nil
	c.VersionConfigList = nil
}

//New 新建配置实例
// func New() (*Config, error) { return new(Config).Init() }
