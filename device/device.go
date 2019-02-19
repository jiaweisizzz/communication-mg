package device

import (
	"fmt"
	"mg-coms/config"
	"mg-coms/notice"
	"sync"
	"time"
)

//DeviceMap 设备列表
type DeviceMap map[int]*Device

//Point 点基本类型
type Point struct {
	ID    int       //id
	Name  string    //名称
	Time  time.Time //更新时间
	Scale float32   //系数
	Data  float32
	fn    func() //Callback
	sync.RWMutex
}

//updateTime 更新时间
func (p *Point) updateTime() {
	p.Time = time.Now()
}

// //Yc 遥测
// type Yc struct {
// 	Point
// }
//
// //Yt 遥调
// type Yt struct {
// 	Point
// }
//
// //Yx 遥信
// type Yx struct {
// 	Point
// }
//
// //Yk 遥控
// type Yk struct {
// 	Point
// }
//
// //Ym 遥脉
// type Ym struct {
// 	Point
// }

// //Datapool 数据
// type Datapool struct {
// 	Ycs []Yc //遥测
// 	Yxs []Yx //遥信
// 	Yms []Ym //遥脉
// 	Yks []Yk //遥控
// 	Yts []Yt //遥调
// }

//Device 设备
type Device struct {
	ID      int       //id
	Tnumber int       //终端号
	Name    string    //名称
	Time    time.Time //更新时间
	Channel string    //所属通道
	Status  bool
	Ycs     []Point //遥测
	Yxs     []Point //遥信
	Yms     []Point //遥脉
	Yks     []Point //遥控
	Yts     []Point //遥调
}

func (d *Device) updateTime() {
	d.Time = time.Now()
}

//Start 启动设备
func (d *Device) Start() {

}

//GetYmByTime 获得Ym值， 如果超时则返回超时通知
func (d *Device) GetYmByTime(id int, lastTime time.Time) (*Point, bool) {
	pt := d.GetYm(id)
	need := lastTime.Before(pt.Time)

	return pt, need
}

//GetYcByTime 获得Yc值， 如果超时则返回超时通知
func (d *Device) GetYcByTime(id int, lastTime time.Time) (*Point, bool) {
	pt := d.GetYc(id)
	need := lastTime.Before(pt.Time)
	return pt, need
}

//GetYxByTime 获得Yx值，如超时则返回超时通知
func (d *Device) GetYxByTime(id int, lastTime time.Time) (*Point, bool) {
	pt := d.GetYx(id)
	need := lastTime.Before(pt.Time)
	return pt, need
}

//Print 打印设备表
func (dm *DeviceMap) Print() {
	for index, dev := range *dm {
		fmt.Println(index, *dev)
	}
}

//Init 设备初始化
func (dm DeviceMap) Init(list []config.DeviceConfig) DeviceMap {
	if dm == nil {
		dm = make(map[int]*Device, 0)
	}
	for _, dc := range list {
		dm[dc.ChannelID*10000+dc.ID] = new(Device).init(dc)
	}
	return dm
}

func (d *Device) init(dc config.DeviceConfig) *Device {
	d.ID = dc.ID
	d.Tnumber = dc.Tnumber
	d.Name = dc.Name
	d.Time = time.Now()
	d.Channel = dc.Channel
	d.Status = dc.Status
	d.Ycs = initYcs(dc.Ycs)
	d.Yxs = initYxs(dc.Yxs)
	d.Yks = initYks(dc.Yks)
	d.Yts = initYts(dc.Yts)
	d.Yms = initYms(dc.Yms)

	return d
}

//SetYc 更新遥测
func (d *Device) SetYc(num int, value float32) {
	for index := range d.Ycs {
		if d.Ycs[index].ID == num {
			d.Ycs[index].Lock()
			if d.Ycs[index].Data == value {
				d.Ycs[index].Unlock()
				break
			}
			d.Ycs[index].Data = value
			if d.Ycs[index].fn != nil {
				d.Ycs[index].fn()
			}
			d.Ycs[index].updateTime()
			d.updateTime()
			d.Ycs[index].Unlock()
			break
		}
	}
}

//GetYc  获得遥测数据
func (d *Device) GetYc(pid int) *Point {
	for index := range d.Ycs {
		if d.Ycs[index].ID == pid {
			return &d.Ycs[index]
		}
	}
	return &Point{}
}

func timetostring(t time.Time) string {
	// day := t.Unix()
	// d := time.Unix(day, 0)
	s := t.Format("2006-01-02 15:04:05")
	return s
}

//SetYx 更新遥测
func (d *Device) SetYx(num int, value bool) {
	for index := range d.Yxs {
		if d.Yxs[index].ID == num {
			d.Yxs[index].Lock()
			if value {
				if d.Yxs[index].Data != 1 {
					d.Yxs[index].Data = 1
					d.Yxs[index].updateTime()
					data := timetostring(d.Yxs[index].Time) + "," + d.Name + "," + d.Yxs[index].Name + "," + "报警" + "\n"
					notice.WriteMessage(data)
				} else {
					d.Yxs[index].Data = 1
					d.Yxs[index].Unlock()
					break
				}
			} else {
				if d.Yxs[index].Data != 0 {
					d.Yxs[index].Data = 0
					d.Yxs[index].updateTime()
					data := timetostring(d.Yxs[index].Time) + "," + d.Name + "," + d.Yxs[index].Name + "," + "恢复" + "\n"
					notice.WriteMessage(data)
				} else {
					d.Yxs[index].Data = 0
					d.Yxs[index].Unlock()
					break
				}
			}
			if d.Yxs[index].fn != nil {
				d.Yxs[index].fn()
			}

			d.updateTime()
			d.Yxs[index].Unlock()
			break
		}
	}
}

//GetYx 获得遥信数据
func (d *Device) GetYx(pid int) *Point {
	for index := range d.Yxs {
		if d.Yxs[index].ID == pid {
			return &d.Yxs[index]
		}
	}
	return &Point{}
}

//SetYk 更新遥测updateTime
func (d *Device) SetYk(num int, value float32) {
	for index := range d.Yks {
		if d.Yks[index].ID == num {
			d.Yks[index].Lock()
			d.Yks[index].Scale = value
			if d.Yks[index].fn != nil {
				d.Yks[index].fn()
			}
			d.Yks[index].updateTime()
			d.updateTime()
			d.Yks[index].Unlock()
			break
		}
	}
}

//GetYk 获得遥控数据
func (d *Device) GetYk(pid int) Point {
	for index := range d.Yks {
		if d.Yks[index].ID == pid {
			return d.Yks[index]
		}
	}
	return Point{}
}

//SetYt 更新遥测
func (d *Device) SetYt(num int, value float32) {
	for index := range d.Yts {
		if d.Yts[index].ID == num {
			d.Yts[index].Lock()
			d.Yts[index].Scale = value
			if d.Yts[index].fn != nil {
				d.Yts[index].fn()
			}
			d.Yts[index].updateTime()
			d.updateTime()
			d.Yts[index].Unlock()
			break
		}
	}
}

//GetYt 获得遥调数据
func (d *Device) GetYt(pid int) Point {
	for index := range d.Yts {
		if d.Yts[index].ID == pid {
			return d.Yts[index]
		}
	}
	return Point{}
}

//SetYm 更新遥测
func (d *Device) SetYm(num int, value float32) {
	for index := range d.Yms {
		if d.Yms[index].ID == num {
			d.Yms[index].Lock()
			d.Yms[index].Data = value
			if d.Yms[index].fn != nil {
				d.Yms[index].fn()
			}
			d.Yms[index].updateTime()
			d.updateTime()
			d.Yms[index].Unlock()
			break
		}
	}
}

func (d *Device) GetYm(pid int) *Point {
	for index := range d.Yms {
		if d.Yms[index].ID == pid {
			return &d.Yms[index]
		}
	}
	return &Point{}
}

func initYcs(ycl []config.Point) []Point {
	var yl []Point
	if yl == nil {
		yl = make([]Point, 0, 100)
	}
	for _, yc := range ycl {
		yl = append(yl, newYc(yc.ID, yc.Name, time.Now(), yc.Scale, 0))
	}
	return yl
}

func newYc(id int, name string, time time.Time, scale float32, data float32) Point {
	yc := Point{}
	yc.ID = id
	yc.Name = name
	yc.Time = time
	yc.Scale = scale
	yc.Data = data
	return yc
}

func initYxs(yxl []config.Point) []Point {
	var yl []Point
	if yl == nil {
		yl = make([]Point, 0, 100)
	}
	for _, yx := range yxl {
		yl = append(yl, newYx(yx.ID, yx.Name, time.Now(), yx.Scale, 0))
	}
	return yl
}

func newYx(id int, name string, time time.Time, scale float32, data float32) Point {
	yx := Point{}
	yx.ID = id
	yx.Name = name
	yx.Time = time
	yx.Scale = scale
	yx.Data = data
	return yx
}

func initYks(ykl []config.Point) []Point {
	var yl []Point
	if yl == nil {
		yl = make([]Point, 0, 100)
	}
	for _, yk := range ykl {
		yl = append(yl, newYk(yk.ID, yk.Name, time.Now(), yk.Scale))
	}
	return yl
}

func newYk(id int, name string, time time.Time, scale float32) Point {
	yk := Point{}
	yk.ID = id
	yk.Name = name
	yk.Time = time
	yk.Scale = scale
	return yk
}

func initYts(ytl []config.Point) []Point {
	var yl []Point
	if yl == nil {
		yl = make([]Point, 0, 100)
	}
	for _, yt := range ytl {
		yl = append(yl, newYt(yt.ID, yt.Name, time.Now(), yt.Scale))
	}
	return yl
}

func newYt(id int, name string, time time.Time, scale float32) Point {
	yt := Point{}
	yt.ID = id
	yt.Name = name
	yt.Time = time
	yt.Scale = scale
	return yt
}

func initYms(yml []config.Point) []Point {
	var yl []Point
	if yl == nil {
		yl = make([]Point, 0, 100)
	}
	for _, ym := range yml {
		yl = append(yl, newYm(ym.ID, ym.Name, time.Now(), ym.Scale, 0))
	}
	return yl
}

func newYm(id int, name string, time time.Time, scale float32, data float32) Point {
	ym := Point{}
	ym.ID = id
	ym.Name = name
	ym.Time = time
	ym.Scale = scale
	ym.Data = data
	return ym
}

//New 新建设备实例
func New(list []config.DeviceConfig) DeviceMap { return new(DeviceMap).Init(list) }
