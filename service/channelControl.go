package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mg-coms/config"
	"net/http"
	"runtime"

	"github.com/gorilla/mux"
)

//Channel 通道
type Channel struct {
	ID        int
	Name      string
	Crc       string
	Source    bool
	DeviceNum int
	Running   bool
	Port      config.PortConfig
	Potocol   config.PotocolConfig
}

//Channel 通道
func (s *Service) Channel(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	// channelName := string(r.Form.Get("name"))
	var bean ServiceBean
	var ch Channel
	data, err := s.coms.GetChannelByName(vars["name"])
	if err != nil {
		bean.Error = true
		bean.Message = err.Error()
	}
	i := 0
	for _, dev := range data.DeviceMap {
		if dev != nil {
			i++
		}
	}
	// i := len(data.DeviceMap)
	ch.ID = data.ID
	ch.Name = data.Name
	ch.Running = data.Running
	ch.DeviceNum = i
	bean.Error = false
	bean.Data = ch
	bytes, _ := json.Marshal(bean)
	fmt.Fprint(w, string(bytes))
}

//SearchOneChannel 查看某个通道
func (s *Service) SearchOneChannel(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	var bean ServiceBean
	// channelName := string(r.Form.Get("name"))
	data, err := s.config.GetChannelByName(vars["channelname"])
	if err != nil {
		bean.Error = true
		bean.Message = err.Error()
	}
	data.Port.Name = page(data.Port.Name)
	fmt.Println(data)
	bean.Error = false
	bean.Data = data
	bytes, _ := json.Marshal(bean)
	fmt.Fprint(w, string(bytes))
}

func page(pname string) string {
	if runtime.GOOS == "linux" {
		switch pname {
		case "/dev/ttyS0":
			return "COM1"
		case "/dev/ttyS1":
			return "COM2"
		case "/dev/ttyS2":
			return "COM3"
		case "/dev/ttyS3":
			return "COM4"
		case "/dev/ttyS4":
			return "COM5"
		case "/dev/ttyS5":
			return "COM6"
		case "/dev/ttyS6":
			return "COM7"
		case "/dev/ttyS7":
			return "COM8"
		case "/dev/ttyS8":
			return "COM9"
		case "/dev/ttyS9":
			return "COM10"
		case "/dev/ttyS10":
			return "COM11"
		case "/dev/ttyS11":
			return "COM12"
		case "/dev/ttyS12":
			return "COM13"
		case "/dev/ttyS13":
			return "COM14"
		case "/dev/ttyS14":
			return "COM15"
		case "/dev/ttyS15":
			return "COM16"
		case "/dev/ttyS16":
			return "COM17"
		case "/dev/ttyS17":
			return "COM18"
		}
	}
	return pname
}

//SearchAllChannel 查看所有通道
func (s *Service) SearchAllChannel(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var bean ServiceBean
	data, err := s.config.GetChannelConfig()
	if err != nil {
		bean.Error = true
		bean.Message = err.Error()
	}
	var cl = make([]Channel, 0)
	for _, ch := range data {
		ch.Port.Name = page(ch.Port.Name)
		c := Channel{
			ID:        ch.ID,
			Name:      ch.Name,
			Crc:       ch.Crc,
			Source:    ch.Source,
			Potocol:   ch.Potocol,
			DeviceNum: len(ch.Device),
			Running:   true,
			Port:      ch.Port,
		}
		cl = append(cl, c)
	}
	bean.Error = false
	bean.Data = cl
	bytes, _ := json.Marshal(bean)
	fmt.Fprint(w, string(bytes))
}

//ModifiChannel 修改通道
func (s *Service) ModifiChannel(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	// vars := mux.Vars(r)
	// w.WriteHeader(http.StatusOK)
	// channelName := string(r.Form.Get("name"))
	// newID := string(r.Form.Get("value"))
	// fmt.Println(newID)
	var ch Channel
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &ch)
	// channelName := ch["Name"].(string)
	// Pname := ch["Port"]["Name"].(string)
	// fmt.Println(Pname)
	// Pbaud := ch["Port"]["Baud"].(string)
	// Psize := ch["Port"]["Size"].(string)
	// Pstopbits := ch["Port"]["StopBits"].(string)
	// Pparity := ch["Port"]["Parity"].(string)
	var bean ServiceBean
	channelName := ch.Name
	channelCrc := ch.Crc
	if ch.Port.Ptype == "com" {
		Pname := ch.Port.Name
		Pbaud := ch.Port.Baud
		Psize := ch.Port.Size
		Preadtimeout := ch.Port.ReadTimeout
		Pstopbits := ch.Port.StopBits
		Pparity := ch.Port.Parity

		err := s.config.SetPortChannelByName(channelName, channelCrc, Pname, Pparity, Pbaud, Psize, Preadtimeout, Pstopbits)
		if err != nil {
			bean.Error = true
			bean.Message = err.Error()
		}
	} else {
		pip := ch.Port.Ip
		pconntype := ch.Port.Conntype
		pportnum := ch.Port.Portnum
		err := s.config.SetTcpChannelByName(channelName, pconntype, pip, pportnum, ch.Potocol.Data)
		if err != nil {
			bean.Error = true
			bean.Message = err.Error()
		}
	}
	bean.Error = false
	bytes, _ := json.Marshal(bean)
	fmt.Fprint(w, string(bytes))
}

//AppendChannel 添加通道
func (s *Service) AppendChannel(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	// vars := mux.Vars(r)
	// w.WriteHeader(http.StatusOK)
	var ch Channel
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &ch)
	if err != nil {
		fmt.Println(err)
	}
	var bean ServiceBean
	channelName := ch.Name
	channelID := ch.ID
	channelCrc := ch.Crc
	channelPtype := ch.Port.Ptype
	channelPotocol := ch.Potocol.Name

	var channelSource bool
	if ch.Potocol.Name == "MODBUS" || ch.Potocol.Name == "MODBUSTCP" || ch.Potocol.Name == "645" {
		channelSource = true
	} else {
		channelSource = false
	}

	if channelPtype == "com" {
		Pname := ch.Port.Name
		Pbaud := ch.Port.Baud
		Psize := ch.Port.Size
		Preadtimeout := ch.Port.ReadTimeout
		Pstopbits := ch.Port.StopBits
		Pparity := ch.Port.Parity
		err = s.config.AddPortChannel(channelName, channelCrc, channelPtype, Pname, Pparity, channelPotocol, channelID, Pbaud, Psize, Preadtimeout, Pstopbits, channelSource)
		if err != nil {
			bean.Error = true
			bean.Message = err.Error()
		}
	} else {
		ctype := ch.Port.Conntype
		Pip := ch.Port.Ip
		Pportnum := ch.Port.Portnum
		err = s.config.AddTcpChannel(channelName, channelPtype, channelPotocol, ctype, Pip, Pportnum, channelID, channelSource, ch.Potocol.Data)
		if err != nil {
			bean.Error = true
			bean.Message = err.Error()
		}
	}
	bean.Error = false
	bytes, _ := json.Marshal(bean)
	fmt.Fprint(w, string(bytes))
}

//RemoveChannel 删除通道
func (s *Service) RemoveChannel(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	// channelName := string(r.Form.Get("name"))
	var bean ServiceBean
	err := s.config.DeleteChannel(vars["channelname"])
	if err != nil {
		bean.Error = true
		bean.Message = err.Error()
	}
	bean.Error = false
	bytes, _ := json.Marshal(bean)
	fmt.Fprint(w, string(bytes))
}
