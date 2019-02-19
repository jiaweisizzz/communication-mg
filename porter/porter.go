package porter

import (
	"errors"
	"fmt"
	"mg-coms/config"
	"mg-coms/serial"
	"net"
	"sync"
	"time"
)

type PorterMap map[string]*Porter

//Porter 端口管理
type Porter struct {
	Ptype        string    //端口类型
	Ctype        string    //连接类型
	Ip           string    //IP地址
	Portnum      string    //端口号
	Portname     string    //端口名
	Baud         int       //波特率
	ReadTimeout  int       //超时时间
	Size         int       //字节大小
	Parity       string    //校验位
	StopBits     int       //停止位
	Port         ReadWrite //读写接口
	Ok           bool      //是否准备就绪
	exit         bool      //退出标志
	Channel      string    //所属通道
	buff         []byte    //读写缓存
	Disconnected bool      //断开连接标志
	Callback     map[time.Time]func(string, string)
	sync.RWMutex
}

func (p *Porter) SetCallBack(t time.Time, f func(string, string)) {
	p.Callback[t] = f
}

func (p *Porter) RemoveCallBack(t time.Time) {
	delete(p.Callback, t)
}

//ReadWrite 读写接口
type ReadWrite interface {
	Read(b []byte) (n int, err error)
	Write(b []byte) (n int, err error)
	Close() error
	SetReadDeadline(t time.Time) error
	SetWriteDeadline(t time.Time) error
}

//Init 通道初始化
func (pl PorterMap) Init(p []config.PortConfig) PorterMap {
	if pl == nil {
		pl = make(map[string]*Porter, 0)
	}
	for _, pp := range p {
		pl[pp.Channel] = new(Porter).init(pp)
	}
	return pl
}

func (pt *Porter) init(pp config.PortConfig) *Porter {
	if pp.Ptype == "com" {
		pt.portinit(pp)
	} else {
		pt.tcpinit(pp)
	}
	return pt
}

//portinit 端口初始化
func (pt *Porter) portinit(pp config.PortConfig) {
	pt.Ptype = pp.Ptype
	pt.Portname = pp.Name
	pt.Baud = pp.Baud
	pt.ReadTimeout = pp.ReadTimeout
	pt.Size = pp.Size
	pt.Parity = pp.Parity
	pt.StopBits = pp.StopBits
	pt.Channel = pp.Channel
	pt.Callback = make(map[time.Time]func(string, string))
	pt.buff = make([]byte, 1024)
}

//tcpinit tcp初始化
func (pt *Porter) tcpinit(pp config.PortConfig) {
	pt.Ctype = pp.Conntype
	pt.Ip = pp.Ip
	pt.Portnum = pp.Portnum
	pt.Callback = make(map[time.Time]func(string, string))
	pt.buff = make([]byte, 1024)
}

//New 新建端口实例
func New(p []config.PortConfig) PorterMap {
	return new(PorterMap).Init(p)
}

//Open 打开端口
func (pt *Porter) Open() error {
	var err error
	if pt.Ptype == "com" {
		err = pt.portopen()
	} else {
		pt.tcpopen()
	}
	return err
}

//portopen 打开端口
func (pt *Porter) portopen() error {
	var parity serial.Parity
	var stop serial.StopBits
	switch pt.StopBits {
	case 1:
		stop = serial.Stop1
	case 15:
		stop = serial.Stop1Half
	case 2:
		stop = serial.Stop2
	default:
		pt.Ok = false
		return errors.New("端口配置错误1")
	}

	switch pt.Parity {
	case "N":
		parity = serial.ParityNone
	case "O":
		parity = serial.ParityOdd
	case "E":
		parity = serial.ParityEven
	case "M":
		parity = serial.ParityMark
	case "S":
		parity = serial.ParitySpace
	default:
		pt.Ok = false
		return errors.New("端口配置错误2")
	}

	var c = serial.Config{
		Name:        pt.Portname,
		Baud:        pt.Baud,
		ReadTimeout: time.Duration(pt.ReadTimeout),
		Size:        byte(pt.Size),
		Parity:      parity,
		StopBits:    stop,
	}
	var err error
	pt.Port, err = serial.OpenPort(&c)
	if err != nil {
		pt.Ok = false
		fmt.Print(err)
		return err
	}

	pt.Ok = true
	return nil
}

func (pt *Porter) tcpserviceopen() {
	localaddr, _ := net.ResolveTCPAddr("tcp", pt.Ip+":"+pt.Portnum)
	ln, err := net.ListenTCP("tcp", localaddr)
	if err != nil {
		pt.Ok = false
	}
	go func() {
		for {
			conn, err := ln.AcceptTCP()
			if err != nil {
				continue
			}
			if pt.exit == true {
				conn.Close()
				pt.Ok = false
				pt.Port = nil
				return
			}
			pt.Ok = true
			pt.Port = conn
		}
	}()
}

//tcpopen 客户端tcp连接
func (pt *Porter) tcpopen() {
	if pt.Ctype == "服务端" {
		localaddr, _ := net.ResolveTCPAddr("tcp", pt.Ip+":"+pt.Portnum)
		ln, err := net.ListenTCP("tcp", localaddr)
		if err != nil {
			pt.Ok = false
			fmt.Println(err)
		}

		go func() {
			for {
				time.Sleep(10 * time.Millisecond)
				ln.SetDeadline(time.Now().Add(400 * time.Millisecond))
				conn, err := ln.AcceptTCP()
				ln.SetDeadline(time.Time{})
				if pt.exit == true {
					if pt.Port != nil {
						pt.Port.Close()
					}

					pt.Ok = false
					return
				}
				if err != nil {
					continue
				}

				pt.Ok = true
				pt.Port = conn
			}
		}()
	} else {
		go func() {
			for {
				if pt.Ok == false {
					pt.Lock()
					if pt.Port != nil {
						pt.Port.Close()
						fmt.Println("已关闭socket")
					}
					conn, err := net.Dial("tcp", pt.Ip+":"+pt.Portnum)
					fmt.Println("正在重连")
					if err != nil {
						pt.Ok = false
						pt.Port = nil
					} else {
						fmt.Println("已连接")
						pt.Ok = true
						pt.Port = conn
						pt.Disconnected = true
					}
					pt.Unlock()
				}
				time.Sleep(3 * time.Second)

				// conn, err := net.Dial("tcp", pt.Ip+":"+pt.Portnum)
				// fmt.Println("tcptcptcptcptcptcptcptcptcptcp", pt.Ip+":"+pt.Portnum)
				// if err != nil {
				// 	fmt.Println("kkkkk")
				// 	pt.Ok = false
				// 	pt.Port = nil
				// 	fmt.Println("nooonononononono", err)
				// } else {
				// 	if pt.Ok == false {
				// 		defer conn.Close()
				// 		pt.Ok = true
				// 		// conn.SetReadDeadline(time.Time{})
				// 		// conn.SetReadDeadline(time.Now().Add(time.Duration(1000000) * time.Millisecond))
				// 		// conn.SetWriteDeadline(time.Time{})
				//
				// 		// conn.SetWriteDeadline(time.Time{})
				// 		pt.Port = conn
				// 	}
				// }
				// time.Sleep(3 * time.Second)
			}
		}()
	}
}

//OpenAllPorter 打开所有端口
func (pl *PorterMap) OpenAllPorter() {
	for _, p := range *pl {
		err := p.Open()
		if err != nil {
			fmt.Print(err)
		}
	}
}

//Read 读取端口数据
func (pt *Porter) Read() ([]byte, error) {
	if !pt.Ok {
		time.Sleep(500 * time.Millisecond)
		return nil, errors.New("duduankouweibeichushihua")
	}
	pt.Lock()
	defer pt.Unlock()
	if pt.Port == nil {
		time.Sleep(500 * time.Millisecond)
		pt.Ok = false
		return nil, nil
	}
	pt.Port.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	//需要循环读取直到EOF
	n, _ := pt.Port.Read(pt.buff)
	if pt.Port == nil {
		return nil, nil
	}
	pt.Port.SetReadDeadline(time.Time{})
	for _, call := range pt.Callback {
		if call != nil && n > 0 {
			call(string(pt.buff[:n]), "收")
		}
	}
	fmt.Println("在读端口")
	return pt.buff[:n], nil
}

//Write 向端口写数据
func (pt *Porter) Write(data string, len int) error {
	if !pt.Ok {
		return errors.New("duankouweibeizhengquechushihua2")
	}
	for _, call := range pt.Callback {
		if call != nil {
			call(data, "发")
		}
	}
	pt.Lock()
	defer pt.Unlock()
	if pt.Port == nil {
		pt.Ok = false
		return nil
	}
	pt.Port.SetWriteDeadline(time.Now().Add(500 * time.Millisecond))
	n, err := pt.Port.Write([]byte(data))
	pt.Port.SetWriteDeadline(time.Time{})
	if err != nil {
		pt.Ok = false
		return errors.New("写数据异常")
	}

	if len != n {
		return errors.New("数据长度异常")
	}
	return nil
}

//Close 关闭端口
func (pt *Porter) Close() error {
	// if pt.Ptype == "com" {
	// pt.Port.Close()
	pt.Ok = false
	// }
	pt.exit = true
	return nil
}
