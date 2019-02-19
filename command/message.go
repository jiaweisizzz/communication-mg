package command

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

const (
	//EXIT 退出
	EXIT = iota
	//INIT 重新读取配置并启动
	INIT = iota
	//RESTART 重新启动
	RESTART = iota
	//UPDATEDATA 更新设备数据
	UPDATEDATA = iota
	SELECT     = iota
	EXE        = iota
	CANCEL     = iota
	//YC 下发更新遥测命令
	YC = iota
	//YX 下发更新遥信命令
	YX = iota
	//YM 下发更新遥脉命令
	YM = iota
	//YK 下发遥控命令
	YK = iota
	//YT 下发遥调命令
	YT = iota
)

//Message 消息结构
type Message struct {
	Cmd     int //YK or YT
	SubCmd  int // EXE or SELECT or CANCEL
	Success bool
	Channel string
	Device  int
	Point   int
	Value   interface{}
	// Time   interface{}
}

// //Init   初始化
// func (m *Message) Init() {
// 	// conn, err := net.Dial("tcp", "172.17.0.1:10000")
// 	// if err != nil {
// 	// 	panic(err)
// 	// }
// 	// // go readServer(conn)
// 	// // connectServer(conn)
// 	m.Cmd = INIT
// 	m.Device = "设备01"
// 	// m.Point = 3
// 	// m.Value = "重新读取配置并启动"
// 	// m.Time = time.Now()
// 	// //	go readServer(conn)
// 	// connectServer(conn)
// }

func connectServer(conn net.Conn) {
	input := bufio.NewReader(os.Stdin)
	for {
		name, _ := input.ReadString('\n')
		trimName := strings.Trim(name, "\r\n")
		conn.Write([]byte(trimName + "\n"))
	}
}

func readServer(conn net.Conn) string {
	br := bufio.NewReader(conn)
	for {
		data, err := br.ReadString('\n')
		if err == io.EOF {
			fmt.Println("err")
		}
		// fmt.Printf("%s", data)
		return data
	}
}

// var d string
//
// func connect() {
// 	conn, err := net.Dial("tcp", "172.17.0.1:10000")
// 	if err != nil {
// 		panic(err)
// 	}
// 	d = readServer(conn)
// 	connectServer(conn)
// }
