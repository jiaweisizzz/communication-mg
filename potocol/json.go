package potocol

import (
	"encoding/json"
	"fmt"
	"mg-coms/command"
	"mg-coms/config"
	"mg-coms/device"
	"strconv"
	"time"
)

//PotocolJSON 规约配置结构体
type PotocolJson struct {
	PotocolParm
	found      bool   //找到头为true
	frame      []byte //拼凑的帧
	size       int    //数据长度
	ThreadExit bool
	AllTime    time.Time
	AddTime    time.Time
	time       time.Time //上一次上送全部数据时间
}

//Connect 链接
type Connect struct {
	Cmd    string      `json:"cmd"` //INIT CONNECT YK YT
	Name   string      `json:"name"`
	Key    string      `json:"key"`
	Type   string      `json:"type"` // SELECT EXE CANCEL
	Dev    string      `json:"dev"`
	Dot    int         `json:"dot"`
	Val    interface{} `json:"val"`
	Result bool        `json:"result"`
}

type deviceConfig struct {
	Cmd     string `json:"cmd"`
	Devices `json:"devices"`
}

type Devices []deviceconfig

type deviceconfig struct {
	Name string                 `json:"name,omitempty"`
	Num  int                    `json:"num,omitempty"`
	Yc   map[string]interface{} `json:"yc,omitempty"`
	Yx   map[string]interface{} `json:"yx,omitempty"`
	Ym   map[string]interface{} `json:"ym,omitempty"`
	Yk   map[string]interface{} `json:"yk,omitempty"`
	Yt   map[string]interface{} `json:"yt,omitempty"`
}

func (d Devices) Len() int {
	return len(d)
}

func (d Devices) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

func (d Devices) Less(i, j int) bool {
	return d[j].Num > d[i].Num
}

func (j *PotocolJson) Created() {
	fmt.Println(j.ChannelConfig.Potocol.Data)
	var c = Connect{
		Cmd:  "CONNECT",
		Name: j.ChannelConfig.Potocol.Data.(map[string]interface{})["Name"].(string),
		Key:  j.ChannelConfig.Potocol.Data.(map[string]interface{})["Key"].(string),
	}

	var con, _ = json.Marshal(c)
	con = append(con, []byte("/0")...)
	time.Sleep(time.Millisecond * 1000)
	err := j.Port.Write(string(con), len(con))
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(c)

	go j.RD()
}

var a int

//Cmd JSON
func (j *PotocolJson) Cmd(msg command.Message) {
	cmd := msg.Cmd
	switch cmd {
	case command.YK:
		j.setYk(msg)
	case command.YT:
		j.setYt(msg)
	}
}

func (j *PotocolJson) setYk(msg command.Message) {
	for index, dev := range j.DeviceMap {
		if msg.Device == dev.Tnumber {
			j.DeviceMap[index].SetYk(msg.Point, msg.Value.(float32))
		}
	}
}

func (j *PotocolJson) setYt(msg command.Message) {
	for index, dev := range j.DeviceMap {
		if msg.Device == dev.Tnumber {
			j.DeviceMap[index].SetYt(msg.Point, msg.Value.(float32))
		}
	}
}

func (j *PotocolJson) RD() {
	for {

		buf, _ := j.Port.Read()

		l := len(buf)
		if l != 0 {
			j.found = true
		}
		if j.size+l > 1024 {
			j.size = 0
		}
		for _, b := range buf {
			j.Buffer = append(j.Buffer, b)
		}
		j.size = len(j.Buffer)

		if j.found == true {
			j.found = false
			for index, _ := range j.Buffer {
				if index >= j.size {
					break
				}
				if j.Buffer[index] == '/' && j.Buffer[index+1] == '0' {
					j.frame = j.Buffer[:index]
					j.Buffer = nil
					if len(j.frame) == 0 {
					}
					fmt.Println(string(j.frame))
					var cmd = new(Connect)
					err := json.Unmarshal(j.frame, cmd)
					if err != nil {
						fmt.Println(err)
					}

					switch cmd.Cmd {
					case "CONNECT", "connect":
						fmt.Println("---------------------------CONNECT")
						time.Sleep(10 * time.Second)
						// j.LoadAllData()
						j.LoadDataAll("LEVELONE", false, true, false, true)
						j.LoadDataAll("LEVELTWO", true, false, false, true)
					case "INIT", "init":
						fmt.Println("INIT")
						j.configSend()
					case "YK", "yk":
						fmt.Println("YK")
						var d int
						var pid int
						var daddr int
						var dch string
						tmp, err := strconv.Atoi(cmd.Dev)
						if err != nil {
							fmt.Println("字符串转换成整数失败")
						}
						switch cmd.Type {
						case "SELECT":
							d = command.SELECT
						case "EXE":
							d = command.EXE
						case "CANCEL":
							d = command.CANCEL
						}
						for _, dev := range j.DeviceMap {
							if tmp == dev.Tnumber {
								for _, yk := range j.ChannelConfig.Uptable.AddrYk {
									if dev.Channel == j.DeviceMap[yk.ChannelID*10000+yk.DeviceID].Channel {
										if yk.AcceptID == cmd.Dot {
											pid = yk.ID
											daddr = dev.ID
											dch = dev.Channel
										}
									}
								}
							}
						}
						msg := command.Message{
							Cmd:     command.YK,
							Value:   cmd.Val.(float64) != 0,
							SubCmd:  d,
							Device:  daddr,
							Point:   pid,
							Channel: dch,
						}
						*j.Bus <- msg
					case "YT", "yt":
						fmt.Println("YT")
						var d int
						var pid int
						var daddr int
						var dch string
						tmp, err := strconv.Atoi(cmd.Dev)
						fmt.Println("tmp=============", tmp)
						if err != nil {
							fmt.Println("字符串转换成整数失败")
						}
						switch cmd.Type {
						case "SELECT":
							d = command.SELECT
						case "EXE":
							d = command.EXE
						case "CANCEL":
							d = command.CANCEL
						}
						for _, dev := range j.DeviceMap {
							if tmp == dev.Tnumber {
								for _, yt := range j.ChannelConfig.Uptable.AddrYt {
									if dev.Channel == j.DeviceMap[yt.ChannelID*10000+yt.DeviceID].Channel {
										if yt.AcceptID == cmd.Dot {
											pid = yt.ID
											daddr = dev.ID
											dch = dev.Channel
											break
										}
									}
								}
							}
						}
						msg := command.Message{
							Cmd:     command.YT,
							Value:   cmd.Val.(float64),
							SubCmd:  d,
							Device:  daddr,
							Point:   pid,
							Channel: dch,
						}
						fmt.Println("msg===================", msg)
						*j.Bus <- msg
					}
					break
				}

			}
		}
	}
}

//Do JSON
func (j *PotocolJson) Do() {
	t := time.Now()
	// time.Sleep(time.Millisecond * 200)
	if j.Port.Disconnected {
		j.Created()
		j.Port.Disconnected = false
	}
	//15分钟上送一次全遥信
	if j.AllTime.After(j.AddTime) {
		j.Last = t
		j.LoadDataAll("LEVELONE", false, true, false, true)
		j.LoadDataAll("LEVELTWO", true, false, false, true)
		j.AllTime = time.Now()
		h, _ := time.ParseDuration("15m")
		j.AddTime = j.AllTime.Add(h)
	}
	j.LoadData("LEVELONE", false, true, false, false)
	j.LoadData("LEVELTWO", true, false, false, false)
	// j.LoadData("ALLYM", false, false, true, true)
	j.Last = t
	j.AllTime = time.Now()
}

//LoadAllData 通过Updata表加载全部数据
func (j *PotocolJson) LoadAllData() {
	fmt.Println("加载全部数据")
	ret := new(deviceConfig)
	ret.Cmd = "ALLDATA"
	ret.Devices = make([]deviceconfig, 0)
	for _, d := range j.DeviceMap {
		dev := new(deviceconfig)
		dev.Num = d.Tnumber
		dev.Yx = make(map[string]interface{}, 0)
		//全遥信
		if j.ChannelConfig.Uptable.AddrYx != nil {
			for _, yx := range j.ChannelConfig.Uptable.AddrYx {
				if d.Channel == j.DeviceMap[yx.ChannelID*10000+yx.DeviceID].Channel && yx.DeviceID == d.ID {
					for _, x := range d.Yxs {
						if x.ID == yx.ID {
							res := 0
							if x.Data == 1 {
								res = 1
							}
							dev.Yx[strconv.Itoa(yx.AcceptID)] = res
						}
					}
				}
			}
		}
		dev.Yc = make(map[string]interface{}, 0)
		//全遥测
		if j.ChannelConfig.Uptable.AddrYc != nil {
			for _, yc := range j.ChannelConfig.Uptable.AddrYc {
				if d.Channel == j.DeviceMap[yc.ChannelID*10000+yc.DeviceID].Channel && yc.DeviceID == d.ID {
					for _, c := range d.Ycs {
						if c.ID == yc.ID {
							dev.Yc[strconv.Itoa(yc.AcceptID)] = c.Data
						}
					}
				}
			}
		}
		dev.Ym = make(map[string]interface{}, 0)
		//全遥脉
		if j.ChannelConfig.Uptable.AddrYm != nil {
			for _, ym := range j.ChannelConfig.Uptable.AddrYm {
				if d.Channel == j.DeviceMap[ym.ChannelID*10000+ym.DeviceID].Channel && ym.DeviceID == d.ID {
					for _, m := range d.Yms {
						if m.ID == ym.ID {
							dev.Ym[strconv.Itoa(ym.AcceptID)] = m.Data
						}
					}
				}
			}
		}
		ret.Devices = append(ret.Devices, *dev)
	}
	if len(ret.Devices) == 0 {
		return
	}
	// fmt.Println(len(ret.Devices))
	dt, _ := json.Marshal(*ret)
	dt = append(dt, []byte("/0")...)
	err := j.Port.Write(string(dt), len(dt))
	if err != nil {
		fmt.Println(err)
	}
}

//loadAllData 通过DeviceMap加载全部数据
func (j *PotocolJson) loadAllData() {
	fmt.Println("加载全部数据")
	ret := new(deviceConfig)
	ret.Cmd = "ALLDATA"
	ret.Devices = make([]deviceconfig, 0)
	for _, d := range j.DeviceMap {
		dev := new(deviceconfig)
		dev.Num = d.Tnumber
		dev.Yx = make(map[string]interface{}, 0)
		for _, yx := range d.Yxs {
			res := 0
			if yx.Data == 1 {
				res = 1
			}
			dev.Yx[strconv.Itoa(yx.ID)] = res
		}

		dev.Yc = make(map[string]interface{}, 0)
		for _, yc := range d.Ycs {
			dev.Yc[strconv.Itoa(yc.ID)] = yc.Data
		}
		dev.Ym = make(map[string]interface{}, 0)
		for _, ym := range d.Yms {
			dev.Ym[strconv.Itoa(ym.ID)] = ym.Data
		}
		ret.Devices = append(ret.Devices, *dev)
	}
	if len(ret.Devices) == 0 {
		return
	}
	fmt.Println(ret)
	dt, _ := json.Marshal(*ret)
	dt = append(dt, []byte("/0")...)
	err := j.Port.Write(string(dt), len(dt))
	if err != nil {
		fmt.Println(err)
	}
}

//LoadData 变位数据上送
func (j *PotocolJson) LoadDataAll(cmd string, hasYc, hasYx, hasYm, all bool) {
	ret := new(deviceConfig)
	ret.Cmd = cmd
	ret.Devices = make([]deviceconfig, 0)
	for _, d := range j.DeviceMap {
		dev := new(deviceconfig)
		dev.Num = d.Tnumber
		if hasYc {
			dev.Yc = make(map[string]interface{}, 0)
			if all {
				//全遥测
				if j.ChannelConfig.Uptable.AddrYc != nil {
					for _, yc := range j.ChannelConfig.Uptable.AddrYc {
						if d.Channel == j.DeviceMap[yc.ChannelID*10000+yc.DeviceID].Channel && yc.DeviceID == d.ID {
							for _, c := range d.Ycs {
								if c.ID == yc.ID {
									dev.Yc[strconv.Itoa(yc.AcceptID)] = c.Data
								}
							}
						}

					}
				}
			} else {
				//变位遥测
				if j.ChannelConfig.Uptable.AddrYc != nil {
					for _, yc := range j.ChannelConfig.Uptable.AddrYc {
						if d.Channel == j.DeviceMap[yc.ChannelID*10000+yc.DeviceID].Channel && yc.DeviceID == d.ID {
							for _, c := range d.Ycs {
								if j.Last.Before(c.Time) && c.ID == yc.ID {
									dev.Yc[strconv.Itoa(yc.AcceptID)] = c.Data
								}
							}
						}
					}
				}
			}
		}

		if hasYx {
			dev.Yx = make(map[string]interface{}, 0)
			if all {
				//全遥信
				if j.ChannelConfig.Uptable.AddrYx != nil {
					for _, yx := range j.ChannelConfig.Uptable.AddrYx {
						if d.Channel == j.DeviceMap[yx.ChannelID*10000+yx.DeviceID].Channel && yx.DeviceID == d.ID {
							for _, x := range d.Yxs {
								if x.ID == yx.ID {
									res := 0
									if x.Data == 1 {
										res = 1
									}
									dev.Yx[strconv.Itoa(yx.AcceptID)] = res
								}
							}
						}
					}
				}
			} else {
				//变位遥信
				if j.ChannelConfig.Uptable.AddrYx != nil {
					for _, yx := range j.ChannelConfig.Uptable.AddrYx {
						if d.Channel == j.DeviceMap[yx.ChannelID*10000+yx.DeviceID].Channel && yx.DeviceID == d.ID {
							for _, x := range d.Yxs {
								if j.Last.Before(x.Time) && x.ID == yx.ID {
									// fmt.Println("遥信现在值xxxxxxxx", x.Data)
									res := 0
									if x.Data != 0 {
										res = 1
									}
									dev.Yx[strconv.Itoa(yx.AcceptID)] = res
								}
							}
						}
					}
				}
			}
		}

		if hasYm {
			dev.Ym = make(map[string]interface{}, 0)
			if all {
				//全遥脉
				if j.ChannelConfig.Uptable.AddrYm != nil {
					for _, ym := range j.ChannelConfig.Uptable.AddrYm {
						if d.Channel == j.DeviceMap[ym.ChannelID*10000+ym.DeviceID].Channel && ym.DeviceID == d.ID {
							for _, m := range d.Yms {
								if m.ID == ym.ID {
									dev.Ym[strconv.Itoa(ym.AcceptID)] = m.Data
								}
							}
						}
					}
				}
			} else {
				//变位遥脉
				if j.ChannelConfig.Uptable.AddrYm != nil {
					for _, ym := range j.ChannelConfig.Uptable.AddrYm {
						if d.Channel == j.DeviceMap[ym.ChannelID*10000+ym.DeviceID].Channel && ym.DeviceID == d.ID {
							for _, m := range d.Yms {
								if j.Last.Before(m.Time) && m.ID == ym.ID {
									dev.Ym[strconv.Itoa(ym.AcceptID)] = m.Data
								}
							}
						}
					}
				}
			}
		}
		if len(dev.Yc) == 0 && len(dev.Yx) == 0 {
			continue
		}
		ret.Devices = append(ret.Devices, *dev)
		if len(ret.Devices) == 0 {
			return
		}
		dt, _ := json.Marshal(*ret)
		dt = append(dt, []byte("/0")...)
		err := j.Port.Write(string(dt), len(dt))
		if err != nil {
			fmt.Println(err)
		}
		ret.Devices = nil
	}

	// if len(ret.Devices) == 0 {
	// 	return
	// }
	// dt, _ := json.Marshal(*ret)
	// dt = append(dt, []byte("/0")...)
	// err := j.Port.Write(string(dt), len(dt))
	// if err != nil {
	// 	fmt.Println(err)
	// }
}

func (j *PotocolJson) LoadData(cmd string, hasYc, hasYx, hasYm, all bool) {
	ret := new(deviceConfig)
	ret.Cmd = cmd
	ret.Devices = make([]deviceconfig, 0)
	for _, d := range j.DeviceMap {
		dev := new(deviceconfig)
		dev.Num = d.Tnumber
		if hasYc {
			dev.Yc = make(map[string]interface{}, 0)
			if all {
				//全遥测
				if j.ChannelConfig.Uptable.AddrYc != nil {
					for _, yc := range j.ChannelConfig.Uptable.AddrYc {
						if d.Channel == j.DeviceMap[yc.ChannelID*10000+yc.DeviceID].Channel && yc.DeviceID == d.ID {
							for _, c := range d.Ycs {
								if c.ID == yc.ID {
									dev.Yc[strconv.Itoa(yc.AcceptID)] = c.Data
								}
							}
						}

					}
				}
			} else {
				//变位遥测
				if j.ChannelConfig.Uptable.AddrYc != nil {
					for _, yc := range j.ChannelConfig.Uptable.AddrYc {
						if d.Channel == j.DeviceMap[yc.ChannelID*10000+yc.DeviceID].Channel && yc.DeviceID == d.ID {
							for _, c := range d.Ycs {
								if j.Last.Before(c.Time) && c.ID == yc.ID {
									dev.Yc[strconv.Itoa(yc.AcceptID)] = c.Data
								}
							}
						}
					}
				}
			}
		}

		if hasYx {
			dev.Yx = make(map[string]interface{}, 0)
			if all {
				//全遥信
				if j.ChannelConfig.Uptable.AddrYx != nil {
					for _, yx := range j.ChannelConfig.Uptable.AddrYx {
						if d.Channel == j.DeviceMap[yx.ChannelID*10000+yx.DeviceID].Channel && yx.DeviceID == d.ID {
							for _, x := range d.Yxs {
								if x.ID == yx.ID {
									res := 0
									if x.Data == 1 {
										res = 1
									}
									dev.Yx[strconv.Itoa(yx.AcceptID)] = res
								}
							}
						}
					}
				}
			} else {
				//变位遥信
				if j.ChannelConfig.Uptable.AddrYx != nil {
					for _, yx := range j.ChannelConfig.Uptable.AddrYx {
						if d.Channel == j.DeviceMap[yx.ChannelID*10000+yx.DeviceID].Channel && yx.DeviceID == d.ID {
							for _, x := range d.Yxs {
								if j.Last.Before(x.Time) && x.ID == yx.ID {
									// fmt.Println("遥信现在值xxxxxxxx", x.Data)
									res := 0
									if x.Data != 0 {
										res = 1
									}
									dev.Yx[strconv.Itoa(yx.AcceptID)] = res
								}
							}
						}
					}
				}
			}
		}

		if hasYm {
			dev.Ym = make(map[string]interface{}, 0)
			if all {
				//全遥脉
				if j.ChannelConfig.Uptable.AddrYm != nil {
					for _, ym := range j.ChannelConfig.Uptable.AddrYm {
						if d.Channel == j.DeviceMap[ym.ChannelID*10000+ym.DeviceID].Channel && ym.DeviceID == d.ID {
							for _, m := range d.Yms {
								if m.ID == ym.ID {
									dev.Ym[strconv.Itoa(ym.AcceptID)] = m.Data
								}
							}
						}
					}
				}
			} else {
				//变位遥脉
				if j.ChannelConfig.Uptable.AddrYm != nil {
					for _, ym := range j.ChannelConfig.Uptable.AddrYm {
						if d.Channel == j.DeviceMap[ym.ChannelID*10000+ym.DeviceID].Channel && ym.DeviceID == d.ID {
							for _, m := range d.Yms {
								if j.Last.Before(m.Time) && m.ID == ym.ID {
									dev.Ym[strconv.Itoa(ym.AcceptID)] = m.Data
								}
							}
						}
					}
				}
			}
		}
		if len(dev.Yc) == 0 && len(dev.Yx) == 0 {
			continue
		}
		ret.Devices = append(ret.Devices, *dev)
		// if len(ret.Devices) == 0 {
		// 	return
		// }
		// dt, _ := json.Marshal(*ret)
		// dt = append(dt, []byte("/0")...)
		// err := j.Port.Write(string(dt), len(dt))
		// if err != nil {
		// 	fmt.Println(err)
		// }
		// ret.Devices = nil
	}

	if len(ret.Devices) == 0 {
		return
	}
	dt, _ := json.Marshal(*ret)
	dt = append(dt, []byte("/0")...)
	err := j.Port.Write(string(dt), len(dt))
	if err != nil {
		fmt.Println(err)
	}
}

//configSend 按上行点表发送配置
func (j *PotocolJson) configSend() {
	dc := new(deviceConfig)
	dc.Cmd = "CONFIG"
	dc.Devices = make([]deviceconfig, 0)
	for _, d := range j.DeviceMap {
		dev := new(deviceconfig)
		dev.Name = d.Name
		dev.Num = d.Tnumber
		dev.Yc = make(map[string]interface{}, 0)
		dev.Yx = make(map[string]interface{}, 0)
		dev.Ym = make(map[string]interface{}, 0)
		dev.Yk = make(map[string]interface{}, 0)
		dev.Yt = make(map[string]interface{}, 0)
		if j.ChannelConfig.Uptable.AddrYc != nil {
			for _, yc := range j.ChannelConfig.Uptable.AddrYc {
				if d.Channel == j.DeviceMap[yc.ChannelID*10000+yc.DeviceID].Channel && yc.DeviceID == d.ID {
					dev.Yc[yc.Name] = yc.AcceptID
				}
			}
		}
		if j.ChannelConfig.Uptable.AddrYx != nil {
			for _, yx := range j.ChannelConfig.Uptable.AddrYx {
				if d.Channel == j.DeviceMap[yx.ChannelID*10000+yx.DeviceID].Channel && yx.DeviceID == d.ID {
					dev.Yx[yx.Name] = yx.AcceptID
				}
			}
		}
		if j.ChannelConfig.Uptable.AddrYm != nil {
			for _, ym := range j.ChannelConfig.Uptable.AddrYm {
				if d.Channel == j.DeviceMap[ym.ChannelID*10000+ym.DeviceID].Channel && ym.DeviceID == d.ID {
					dev.Ym[ym.Name] = ym.AcceptID
				}
			}
		}
		if j.ChannelConfig.Uptable.AddrYk != nil {
			for _, yk := range j.ChannelConfig.Uptable.AddrYk {
				if d.Channel == j.DeviceMap[yk.ChannelID*10000+yk.DeviceID].Channel && yk.DeviceID == d.ID {
					dev.Yk[yk.Name] = yk.AcceptID
				}
			}
		}
		if j.ChannelConfig.Uptable.AddrYt != nil {
			for _, yt := range j.ChannelConfig.Uptable.AddrYt {
				if d.Channel == j.DeviceMap[yt.ChannelID*10000+yt.DeviceID].Channel && yt.DeviceID == d.ID {
					dev.Yt[yt.Name] = yt.AcceptID
				}
			}
		}
		if dev.Yc == nil && dev.Yx == nil && dev.Ym == nil && dev.Yk == nil && dev.Yt == nil {
			continue
		}
		dc.Devices = append(dc.Devices, *dev)
	}
	fmt.Println(dc)
	dt, _ := json.Marshal(dc)
	dt = append(dt, []byte("/0")...)
	err := j.Port.Write(string(dt), len(dt))
	if err != nil {
		fmt.Println(err)
	}
}

//configsend 按DeviceMap发送配置
func (j *PotocolJson) configsend() {
	dc := new(deviceConfig)
	dc.Cmd = "CONFIG"
	dc.Devices = make([]deviceconfig, 0)

	for _, d := range j.DeviceMap {
		dev := new(deviceconfig)
		dev.Name = d.Name
		dev.Num = d.Tnumber
		dev.Yc = make(map[string]interface{}, 0)
		dev.Yx = make(map[string]interface{}, 0)
		dev.Ym = make(map[string]interface{}, 0)
		dev.Yk = make(map[string]interface{}, 0)
		dev.Yt = make(map[string]interface{}, 0)
		for _, cyc := range d.Ycs {
			dev.Yc[cyc.Name] = cyc.ID
		}
		for _, cyx := range d.Yxs {
			dev.Yx[cyx.Name] = cyx.ID
		}
		for _, cym := range d.Yms {
			dev.Ym[cym.Name] = cym.ID
		}
		for _, cyk := range d.Yks {
			dev.Yk[cyk.Name] = cyk.ID
		}
		for _, cyt := range d.Yts {
			dev.Yt[cyt.Name] = cyt.ID
		}
		dc.Devices = append(dc.Devices, *dev)
	}
	fmt.Println(dc)
	dt, _ := json.Marshal(dc)
	dt = append(dt, []byte("/0")...)
	err := j.Port.Write(string(dt), len(dt))
	if err != nil {
		fmt.Println(err)
	}
}

//Init JSON
func (j *PotocolJson) Init(cc config.ChannelConfig, dm device.DeviceMap) *PotocolJson {
	j.ChannelConfig = cc
	j.DeviceMap = dm
	j.AllTime = time.Now()
	h, _ := time.ParseDuration("15m")
	j.AddTime = j.AllTime.Add(h)
	return j
}
