package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

type Addrlist struct {
	AddrYc []Ypoint
	AddrYx []Ypoint
	AddrYm []Ypoint
	AddrYk []Ypoint
	AddrYt []Ypoint
}

type Ypoint struct {
	Addr      int    //数据地址
	Name      string //点名
	ID        int    //点号
	AcceptID  int    //104规约用的接收ID
	Funcode   int    //功能码
	DeviceID  int    //所属设备ID
	ChannelID int    //所属通道ID
}

func (c *Config) Creataddrlist() error {
	var y_c = 1
	var y_x = 1
	var y_m = 1
	var y_k = 1
	var y_t = 1
	var tmp = 1
	devs, _ := c.GetDeviceConfig()
	for _, dev := range devs {
		for _, yc := range dev.Ycs {
			yp := new(Ypoint)
			yp.Addr = tmp
			yp.DeviceID = dev.ID
			yp.ChannelID = dev.ChannelID
			yp.AcceptID = y_c
			fc, _ := strconv.Atoi(yc.Data.(map[string]interface{})["Funcode"].(string))
			yp.Funcode = fc
			yp.ID = yc.ID
			yp.Name = dev.Name + "_" + yc.Name
			c.AddrList.AddrYc = append(c.AddrList.AddrYc, *yp)
			y_c++
			tmp += 2
		}
	}
	for _, dev := range devs {
		for _, yx := range dev.Yxs {
			yp := new(Ypoint)
			yp.Addr = tmp
			yp.DeviceID = dev.ID
			yp.ChannelID = dev.ChannelID
			yp.AcceptID = y_x
			fc, _ := strconv.Atoi(yx.Data.(map[string]interface{})["Funcode"].(string))
			yp.Funcode = fc
			yp.ID = yx.ID
			yp.Name = dev.Name + "_" + yx.Name
			c.AddrList.AddrYx = append(c.AddrList.AddrYx, *yp)
			y_x++
			tmp++
		}
	}
	for _, dev := range devs {
		for _, ym := range dev.Yms {
			yp := new(Ypoint)
			yp.Addr = tmp
			yp.DeviceID = dev.ID
			yp.ChannelID = dev.ChannelID
			yp.AcceptID = y_m
			fc, _ := strconv.Atoi(ym.Data.(map[string]interface{})["Funcode"].(string))
			yp.Funcode = fc
			yp.ID = ym.ID
			yp.Name = dev.Name + "_" + ym.Name
			c.AddrList.AddrYm = append(c.AddrList.AddrYm, *yp)
			y_m++
			tmp += 2
		}
	}

	for _, dev := range devs {
		for _, yk := range dev.Yks {
			yp := new(Ypoint)
			yp.Addr = tmp
			yp.DeviceID = dev.ID
			yp.ChannelID = dev.ChannelID
			yp.AcceptID = y_k
			fc, _ := strconv.Atoi(yk.Data.(map[string]interface{})["Funcode"].(string))
			yp.Funcode = fc
			yp.ID = yk.ID
			yp.Name = dev.Name + "_" + yk.Name
			c.AddrList.AddrYk = append(c.AddrList.AddrYk, *yp)
			y_k++
			tmp++
		}
	}

	for _, dev := range devs {
		for _, yt := range dev.Yts {
			yp := new(Ypoint)
			yp.Addr = tmp
			yp.DeviceID = dev.ID
			yp.ChannelID = dev.ChannelID
			yp.AcceptID = y_t
			fc, _ := strconv.Atoi(yt.Data.(map[string]interface{})["Funcode"].(string))
			yp.Funcode = fc
			yp.ID = yt.ID
			yp.Name = dev.Name + "_" + yt.Name
			c.AddrList.AddrYt = append(c.AddrList.AddrYt, *yp)
			y_t++
			tmp += 2
		}
	}

	b, err := json.MarshalIndent(c.AddrList, "", "	")
	if err != nil {
		return errors.New("pointaddr Marshal error")
	}
	path := "point/addrlist.json"
	ioutil.WriteFile(path, b, os.ModePerm)
	return nil
}

//Getaddrlist 加载地址点的映射
func (c *Config) Getaddrlist(a *Addrlist) (*Addrlist, error) {
	f, err := ioutil.ReadFile("point/addrlist.json")
	if err != nil {
		return nil, errors.New("Read addrlist.json File error")
	}
	err = json.Unmarshal(f, a)
	if err != nil {
		return nil, errors.New("地址映射配置文件错误")
	}
	fmt.Println("地址映射配置加载完成")
	return a, nil
}

func (c Config) GetAddrPoint(cname string) (*Addrlist, error) {
	ch, err := c.GetChannelByName(cname)
	if err != nil {
		return nil, errors.New("点的配置时获取通道失败")
	}
	return ch.Uptable, nil
}

// func (a *Addrlist) GetPointAddr(cname string) (*Addrlist, error) {
// 	f, err := ioutil.ReadFile("pointaddr/" + cname + ".json")
// 	if err != nil {
// 		return nil, errors.New("Read pointaddr File error")
// 	}
// 	err = json.Unmarshal(f, a)
// 	if err != nil {
// 		return nil, errors.New("读取新的地址配置错误")
// 	}
// 	return a, nil
// }

func (c *Config) SetPointAddr(pointaddryc, pointaddryx, pointaddrym, pointaddryk, pointaddryt []Ypoint, cname string) error {
	al := new(Addrlist)
	al.AddrYc = pointaddryc
	al.AddrYx = pointaddryx
	al.AddrYm = pointaddrym
	al.AddrYk = pointaddryk
	al.AddrYt = pointaddryt
	c.SavePointAddr(*al, cname)
	return nil
}

func (c *Config) SavePointAddr(al Addrlist, cname string) {
	for _, ch := range c.ChannelConfigList {
		if ch.Name == cname {
			ch.Uptable = &al
			ch.Save()
		}
	}
	// al.save(cname)
}

func (al *Addrlist) save(cname string) error {
	b, err := json.Marshal(*al)
	if err != nil {
		return errors.New("Marshal save error")
	}
	path := "channels/" + cname + ".json"
	ioutil.WriteFile(path, b, os.ModePerm)
	return nil
}

func (c *Config) DeletePointaddr(pointaddryc, pointaddryx, pointaddrym, pointaddryk, pointaddryt []Ypoint, cname string) error {
	fmt.Println("deletepointaddr", cname)
	a := new(Addrlist)
	a, err := c.GetAddrPoint(cname)
	fmt.Println("a", a)
	if err != nil {
		return errors.New("删除时获取配置出错")
	}
	for _, pyc := range pointaddryc {
		for index, yc := range a.AddrYc {
			if yc.ID == pyc.ID {
				a.AddrYc = append(a.AddrYc[:index], a.AddrYc[index+1:]...)
			}
		}
	}
	for _, pyx := range pointaddryx {
		for index, yx := range a.AddrYx {
			if yx.ID == pyx.ID {
				a.AddrYx = append(a.AddrYx[:index], a.AddrYx[index+1:]...)
			}
		}
	}
	for _, pym := range pointaddrym {
		for index, ym := range a.AddrYm {
			if ym.ID == pym.ID {
				a.AddrYm = append(a.AddrYm[:index], a.AddrYm[index+1:]...)
			}
		}
	}
	for _, pyk := range pointaddryk {
		for index, yk := range a.AddrYk {
			if yk.ID == pyk.ID {
				a.AddrYk = append(a.AddrYk[:index], a.AddrYk[index+1:]...)
			}
		}
	}
	for _, pyt := range pointaddryt {
		for index, yt := range a.AddrYt {
			if yt.ID == pyt.ID {
				a.AddrYt = append(a.AddrYt[:index], a.AddrYt[index+1:]...)
			}
		}
	}
	fmt.Println(a)
	c.SavePointAddr(*a, cname)
	return nil
}
