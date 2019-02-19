package websocketio

import (
	"encoding/json"
	"fmt"
	"mg-coms/channel"
	"time"

	"golang.org/x/net/websocket"
)

// type Chan chan ConnectData

type Deal struct {
	// 注册了的连接器
	connections map[*Connection]bool
	// 从连接器中发入的信息
	broadcast chan *Connection
	// 从连接器中注册请求
	register chan *Connection
	// 从连接器中注销请求
	unregister chan *Connection
}

var d = Deal{
	broadcast:   make(chan *Connection),
	register:    make(chan *Connection),
	unregister:  make(chan *Connection),
	connections: make(map[*Connection]bool),
}

var ChannelMap channel.ChannelMap

type Connection struct {
	ws    *websocket.Conn
	send  chan []byte
	cData ConnectData
	buf   string
}

func (c *Connection) Send(m MSG) error {
	return websocket.JSON.Send(c.ws, &m)
}

type ConnectData struct {
	ChannelId int
	Temp      time.Time
}

func WsHandler(ws *websocket.Conn) {
	fmt.Println("WsHandler-----------------------------------------------------------")
	c := new(Connection)
	c.send = make(chan []byte)
	c.cData = ConnectData{}
	c.ws = ws
	c.buf = ""
	d.register <- c
	defer func() { d.unregister <- c }()
	// c.NoticeDataList()
	c.Read()
}

func Server() {
	go d.Run()
}

func (d *Deal) Run() {
	for {
		select {
		case c := <-d.register:
			d.connections[c] = true
			fmt.Println("connections=======", d.connections)
		case c := <-d.unregister:
			if _, ok := d.connections[c]; ok {
				fmt.Println("断开", c.cData)
				ChannelMap.RemoveCallBack(c.cData.ChannelId, c.cData.Temp)
				delete(d.connections, c)
				close(c.send)
			}
		case m := <-d.broadcast:
			fmt.Println("------", d.connections)
			for c := range d.connections {
				if m == c {
					fmt.Println("->->->", c.cData.Temp)
					ChannelMap.SetCallBack(m.cData.ChannelId, c.cData.Temp, c.CutData)
				}
			}
		}
	}
}

func (c *Connection) Read() {
	fmt.Println("-----Read")
	for {
		err := websocket.Message.Receive(c.ws, &c.buf)
		fmt.Println(",,,,,", c.buf)
		if err != nil {
			fmt.Println("=======", err)
			break
		}
		ChannelMap.RemoveCallBack(c.cData.ChannelId, c.cData.Temp)
		_ = json.Unmarshal([]byte(c.buf), &c.cData)
		if c.cData.ChannelId > 0 {
			c.cData.Temp = time.Now()
			fmt.Println("Read", c.cData)
			d.broadcast <- c
		}
	}
}

type MSG struct {
	Cmd  string
	Data string
}

func (c *Connection) CutData(data, parm string) {
	c.buf = data
	fmt.Println("data", c.buf)
	m := new(MSG)
	m.Cmd = "message"
	m.Data = fmt.Sprintf("%s :%x", parm, c.buf)
	// b, _ := json.Marshal(&m)
	err := websocket.JSON.Send(c.ws, &m)
	if err != nil {
		fmt.Println(err.Error())
	}
}

// func (c *Connection) NoticeDataList() {
// 	c.buf = notice.ReadMessage("../app/day")
// 	m := new(MSG)
// 	m.Cmd = "noticelist"
// 	m.Data = fmt.Sprintf("%s", c.buf)
// 	err := websocket.JSON.Send(c.ws, &m)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 	}
// }

func NoticeData(s string) {
	m := new(MSG)
	m.Cmd = "notice"
	m.Data = fmt.Sprintf("%s", s)
	for conn, _ := range d.connections {
		err := conn.Send(*m)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

}
