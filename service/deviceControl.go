package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

//Device 服务设备结构体
type Device struct {
	ID        int
	Tnumber   int //终端号
	Name      string
	Channel   string
	ChannelID int
	Template  string
	Status    bool
}

//OneDevice 某个实时设备
func (s *Service) OneDevice(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	// channelName := string(r.Form.Get("name"))
	// deviceID := string(r.Form.Get("id"))
	var bean ServiceBean
	data, err := s.coms.GetDeviceByID(vars["channelname"], vars["deviceid"])
	if err != nil {
		bean.Error = true
		bean.Message = err.Error()
	}
	bean.Error = false
	bean.Data = data
	bytes, _ := json.Marshal(bean)
	fmt.Fprint(w, string(bytes))
}

//AllDevices 所有实时设备
func (s *Service) AllDevices(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var bean ServiceBean
	data := s.coms.GetAllDevices()
	var dl = make([]Device, 0)
	for _, dev := range data {
		d := Device{
			ID:      dev.ID,
			Tnumber: dev.Tnumber,
			Name:    dev.Name,
			Channel: dev.Channel,
			Status:  dev.Status,
		}
		dl = append(dl, d)
	}
	bean.Error = false
	bean.Data = dl
	bytes, _ := json.Marshal(bean)
	fmt.Fprint(w, string(bytes))
}

//ConfigDevice 设备管理
func (s *Service) ConfigDevice(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	switch r.Method {
	// case http.MethodGet:
	// 	s.searchDevice(w, r)
	case http.MethodPut:
		s.ModifiDevice(w, r)
	case http.MethodPost:
		s.AppendDevice(w, r)
		// case http.MethodDelete:
		// 	s.removeDevice(w, r)
	}
}

//SearchOneDevice 查看某个设备
func (s *Service) SearchOneDevice(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	// channelName := string(r.Form.Get("name"))
	// deviceID := string(r.Form.Get("id"))
	var bean ServiceBean
	data, err := s.config.GetDeviceByID(vars["channelname"], vars["deviceid"])
	if err != nil {
		bean.Error = true
		bean.Message = err.Error()
	}
	bean.Error = false
	bean.Data = data
	bytes, _ := json.Marshal(bean)
	fmt.Fprint(w, string(bytes))
}

//SearchAllDevice 查看所有设备
func (s *Service) SearchAllDevice(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var bean ServiceBean
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	data, err := s.config.GetChannelDevice(vars["channelname"])
	if err != nil {
		bean.Error = true
		bean.Message = err.Error()
	}
	var dl = make([]Device, 0)
	for _, dev := range data {
		d := Device{
			ID:      dev.ID,
			Tnumber: dev.Tnumber,
			Name:    dev.Name,
			Channel: dev.Channel,
		}
		dl = append(dl, d)
	}
	bean.Error = false
	bean.Data = dl
	bytes, _ := json.Marshal(bean)
	fmt.Fprint(w, string(bytes))
}

//ModifiDevice 修改设备
func (s *Service) ModifiDevice(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	// channelName := string(r.Form.Get("name"))
	// deviceID := string(r.Form.Get("id"))
	// value := string(r.Form.Get("devname"))

	var dev Device
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &dev)

	channelName := dev.Channel
	deviceID := dev.ID
	dTnumber := dev.Tnumber
	dname := dev.Name
	dTemplate := dev.Template
	var bean ServiceBean
	// tmp, err := s.config.GetChannelDevice(channelName)
	// if err != nil {
	// 	fmt.Println("未找到通道")
	// }
	// for _, d := range tmp {
	// 	if d.Tnumber == dTnumber {
	// 		bean.Error = false
	// 		bean.Data = "终端号已存在"
	// 		bytes, _ := json.Marshal(bean)
	// 		fmt.Fprint(w, string(bytes))
	// 		return
	// 	}
	// }
	data, err := s.config.SetDeviceByID(channelName, dname, dTemplate, deviceID, dTnumber)
	if err != nil {
		bean.Error = true
		bean.Message = err.Error()
	} else {
		bean.Error = false
		bean.Data = data
	}
	bytes, _ := json.Marshal(bean)
	fmt.Fprint(w, string(bytes))
}

//AppendDevice 增加设备
func (s *Service) AppendDevice(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var dev Device
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &dev)
	channelName := dev.Channel
	channelID := dev.ChannelID
	deviceName := dev.Name
	deviceID := dev.ID
	dTemplate := dev.Template

	// channelName := string(r.Form.Get("channelname"))
	// deviceName := string(r.Form.Get("name"))
	// deviceID := string(r.Form.Get("id"))
	var bean ServiceBean
	err := s.config.AddDevice(channelName, deviceName, dTemplate, deviceID, channelID)
	if err != nil {
		bean.Error = true
		bean.Data = err.Error()
	}
	bean.Error = false
	bytes, _ := json.Marshal(bean)
	fmt.Fprint(w, string(bytes))
}

//RemoveDevice 删除设备
func (s *Service) RemoveDevice(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	// channelName := string(r.Form.Get("name"))
	// deviceID := string(r.Form.Get("id"))
	var bean ServiceBean
	err := s.config.DeleteDevice(vars["channelname"], vars["deviceid"])
	if err != nil {
		bean.Error = true
		bean.Data = err.Error()
	}
	bean.Error = false
	bytes, _ := json.Marshal(bean)
	fmt.Fprint(w, string(bytes))
}
