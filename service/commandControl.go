package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mg-coms/command"
	"net/http"
)

//CommandControl 命令
type CommandControl struct {
	Cmd      string
	Type     string
	Value    interface{}
	Dchannel string
	Did      int
	Pid      int
}

//Command 消息接收
func (s *Service) Command(w http.ResponseWriter, r *http.Request) {
	fmt.Println("下发消息:")
	r.ParseForm()
	var cn CommandControl
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &cn)
	channel := cn.Dchannel
	did := cn.Did
	pid := cn.Pid
	cmd := cn.Cmd
	kind := cn.Type
	value := cn.Value

	var bean ServiceBean
	bean.Data = ""
	bean.Error = false
	bean.Message = "消息下发成功"
	bytes, _ := json.Marshal(bean)
	var c int
	switch cmd {
	case "yk":
		c = command.YK
	case "yt":
		c = command.YT
	}
	var d int
	switch kind {
	case "select":
		d = command.SELECT
	case "exe":
		d = command.EXE
	case "cancel":
		d = command.CANCEL
	}
	msg := command.Message{
		Cmd:     c,
		Value:   value,
		SubCmd:  d,
		Channel: channel,
		Device:  did,
		Point:   pid,
	}
	s.coms.Commander.Send(msg)
	fmt.Fprint(w, string(bytes))
}

//SearchAllPotocol 查看所有规约
func (s *Service) SearchAllPotocol(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var bean ServiceBean
	data, err := s.config.GetPotocolConfig()
	if err != nil {
		bean.Error = true
		bean.Message = err.Error()
	}
	bean.Error = false
	bean.Data = data
	bytes, _ := json.Marshal(bean)
	fmt.Fprint(w, string(bytes))
}
