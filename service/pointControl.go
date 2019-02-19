package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mg-coms/config"
	"net/http"

	"github.com/gorilla/mux"
)

//Point 服务点表
type Point struct {
	Name    string
	Channel string
	Ycs     []config.Point //遥测
	Yxs     []config.Point //遥信
	Yms     []config.Point //遥脉
	Yks     []config.Point //遥控
	Yts     []config.Point //遥调
}

//ConfigPoint 点表管理
func (s *Service) ConfigPoint(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	switch r.Method {
	// case http.MethodGet:
	// 	s.SearchPoint(w, r)
	case http.MethodPut:
		s.modifiPoint(w, r)
	case http.MethodPost:
		s.appendPoint(w, r)
		// case http.MethodDelete:
		// 	s.removePoint(w, r)
	}
}

//SearchPoint 查看点表
func (s *Service) SearchPoint(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	// channelName := string(r.Form.Get("name"))
	// deviceID := string(r.Form.Get("id"))
	var bean ServiceBean
	var p Point
	data, err := s.config.GetDeviceByID(vars["channelname"], vars["deviceid"])
	if err != nil {
		bean.Error = true
		bean.Message = err.Error()
	}
	p.Name = data.Name
	p.Channel = data.Channel
	p.Ycs = data.Ycs
	p.Yks = data.Yks
	p.Yms = data.Yms
	p.Yts = data.Yts
	p.Yxs = data.Yxs
	bean.Error = false
	bean.Data = p
	bytes, _ := json.Marshal(bean)
	fmt.Fprint(w, string(bytes))
}

//modifiPoint 修改点表
func (s *Service) modifiPoint(w http.ResponseWriter, r *http.Request) {
	// tmp := string(r.Form.Get("kind"))
	r.ParseForm()
	var po map[string]interface{}
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &po)
	tmp := po["kind"].(string)
	channelName := po["name"].(string)
	deviceID := po["did"].(string)
	pointID := po["pid"].(string)
	value := po["scale"].(string)
	switch tmp {
	case "Ycs":
		s.modifiYcsPoint(w, r, channelName, deviceID, pointID, value)
	case "Yxs":
		s.modifiYxsPoint(w, r, channelName, deviceID, pointID, value)
	case "Yms":
		s.modifiYmsPoint(w, r, channelName, deviceID, pointID, value)
	case "Yks":
		s.modifiYksPoint(w, r, channelName, deviceID, pointID, value)
	case "Yts":
		s.modifiYtsPoint(w, r, channelName, deviceID, pointID, value)
	}
}

//modifiYcsPoint 修改遥测点
func (s *Service) modifiYcsPoint(w http.ResponseWriter, r *http.Request, channelName, deviceID, pointID, value string) {
	r.ParseForm()
	// channelName := string(r.Form.Get("name"))
	// deviceID := string(r.Form.Get("did"))
	// pointID := string(r.Form.Get("pid"))
	// value := string(r.Form.Get("scale"))
	var bean ServiceBean
	data, err := s.config.SetYcsPointByID(channelName, deviceID, pointID, value)
	if err != nil {
		bean.Error = true
		bean.Message = err.Error()
	}
	bean.Error = false
	bean.Data = data
	bytes, _ := json.Marshal(bean)
	fmt.Fprint(w, string(bytes))
}

//modifiYxsPoint 修改遥信点
func (s *Service) modifiYxsPoint(w http.ResponseWriter, r *http.Request, channelName, deviceID, pointID, value string) {
	r.ParseForm()
	// channelName := string(r.Form.Get("name"))
	// deviceID := string(r.Form.Get("did"))
	// pointID := string(r.Form.Get("pid"))
	// value := string(r.Form.Get("scale"))
	var bean ServiceBean
	data, err := s.config.SetYxsPointByID(channelName, deviceID, pointID, value)
	if err != nil {
		bean.Error = true
		bean.Message = err.Error()
	}
	bean.Error = false
	bean.Data = data
	bytes, _ := json.Marshal(bean)
	fmt.Fprint(w, string(bytes))
}

//modifiYmsPoint 修改遥脉点
func (s *Service) modifiYmsPoint(w http.ResponseWriter, r *http.Request, channelName, deviceID, pointID, value string) {
	r.ParseForm()
	// channelName := string(r.Form.Get("name"))
	// deviceID := string(r.Form.Get("did"))
	// pointID := string(r.Form.Get("pid"))
	// value := string(r.Form.Get("scale"))
	var bean ServiceBean
	data, err := s.config.SetYmsPointByID(channelName, deviceID, pointID, value)
	if err != nil {
		bean.Error = true
		bean.Message = err.Error()
	}
	bean.Error = false
	bean.Data = data
	bytes, _ := json.Marshal(bean)
	fmt.Fprint(w, string(bytes))
}

//modifiYksPoint 修改遥控点
func (s *Service) modifiYksPoint(w http.ResponseWriter, r *http.Request, channelName, deviceID, pointID, value string) {
	r.ParseForm()
	// channelName := string(r.Form.Get("name"))
	// deviceID := string(r.Form.Get("did"))
	// pointID := string(r.Form.Get("pid"))
	// value := string(r.Form.Get("scale"))
	var bean ServiceBean
	data, err := s.config.SetYksPointByID(channelName, deviceID, pointID, value)
	if err != nil {
		bean.Error = true
		bean.Message = err.Error()
	}
	bean.Error = false
	bean.Data = data
	bytes, _ := json.Marshal(bean)
	fmt.Fprint(w, string(bytes))
}

//modifiYtsPoint 修改遥调点
func (s *Service) modifiYtsPoint(w http.ResponseWriter, r *http.Request, channelName, deviceID, pointID, value string) {
	r.ParseForm()
	// channelName := string(r.Form.Get("name"))
	// deviceID := string(r.Form.Get("did"))
	// pointID := string(r.Form.Get("pid"))
	// value := string(r.Form.Get("scale"))
	var bean ServiceBean
	data, err := s.config.SetYtsPointByID(channelName, deviceID, pointID, value)
	if err != nil {
		bean.Error = true
		bean.Message = err.Error()
	}
	bean.Error = false
	bean.Data = data
	bytes, _ := json.Marshal(bean)
	fmt.Fprint(w, string(bytes))
}

//appendPoint 增加点
func (s *Service) appendPoint(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var po map[string]interface{}
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &po)
	tmp := po["kind"].(string)
	channelName := po["name"].(string)
	deviceID := po["did"].(string)
	pointName := po["pname"].(string)
	pointID := po["pid"].(string)
	switch tmp {
	case "Ycs":
		s.appendYcsPoint(w, r, channelName, deviceID, pointName, pointID)
	case "Yxs":
		s.appendYxsPoint(w, r, channelName, deviceID, pointName, pointID)
	case "Yms":
		s.appendYmsPoint(w, r, channelName, deviceID, pointName, pointID)
	case "Yks":
		s.appendYksPoint(w, r, channelName, deviceID, pointName, pointID)
	case "Yts":
		s.appendYtsPoint(w, r, channelName, deviceID, pointName, pointID)
	}
}

//appendYcsPoint 增加遥测点
func (s *Service) appendYcsPoint(w http.ResponseWriter, r *http.Request, channelName, deviceID, pointName, pointID string) {
	// r.ParseForm()

	// channelName := string(r.Form.Get("name"))
	// deviceID := string(r.Form.Get("did"))
	// pointName := string(r.Form.Get("pname"))
	// pointID := string(r.Form.Get("pid"))
	var bean ServiceBean
	err := s.config.AddYcsPoint(channelName, deviceID, pointName, pointID)
	if err != nil {
		bean.Error = true
		bean.Message = err.Error()
	}
	bean.Error = false
	bytes, _ := json.Marshal(bean)
	fmt.Fprint(w, string(bytes))
}

//appendYxsPoint 增加遥信点
func (s *Service) appendYxsPoint(w http.ResponseWriter, r *http.Request, channelName, deviceID, pointName, pointID string) {
	r.ParseForm()
	// channelName := string(r.Form.Get("name"))
	// deviceID := string(r.Form.Get("did"))
	// pointName := string(r.Form.Get("pname"))
	// pointID := string(r.Form.Get("pid"))
	var bean ServiceBean
	err := s.config.AddYxsPoint(channelName, deviceID, pointName, pointID)
	if err != nil {
		bean.Error = true
		bean.Message = err.Error()
	}
	bean.Error = false
	bytes, _ := json.Marshal(bean)
	fmt.Fprint(w, string(bytes))
}

//appendYmsPoint 增加遥脉点
func (s *Service) appendYmsPoint(w http.ResponseWriter, r *http.Request, channelName, deviceID, pointName, pointID string) {
	r.ParseForm()
	// channelName := string(r.Form.Get("name"))
	// deviceID := string(r.Form.Get("did"))
	// pointName := string(r.Form.Get("pname"))
	// pointID := string(r.Form.Get("pid"))
	var bean ServiceBean
	err := s.config.AddYmsPoint(channelName, deviceID, pointName, pointID)
	if err != nil {
		bean.Error = true
		bean.Message = err.Error()
	}
	bean.Error = false
	bytes, _ := json.Marshal(bean)
	fmt.Fprint(w, string(bytes))
}

//appendYksPoint 增加遥控点
func (s *Service) appendYksPoint(w http.ResponseWriter, r *http.Request, channelName, deviceID, pointName, pointID string) {
	r.ParseForm()
	// channelName := string(r.Form.Get("name"))
	// deviceID := string(r.Form.Get("did"))
	// pointName := string(r.Form.Get("pname"))
	// pointID := string(r.Form.Get("pid"))
	var bean ServiceBean
	err := s.config.AddYksPoint(channelName, deviceID, pointName, pointID)
	if err != nil {
		bean.Error = true
		bean.Message = err.Error()
	}
	bean.Error = false
	bytes, _ := json.Marshal(bean)
	fmt.Fprint(w, string(bytes))
}

//appendYtsPoint 增加遥调点
func (s *Service) appendYtsPoint(w http.ResponseWriter, r *http.Request, channelName, deviceID, pointName, pointID string) {
	r.ParseForm()
	// channelName := string(r.Form.Get("name"))
	// deviceID := string(r.Form.Get("did"))
	// pointName := string(r.Form.Get("pname"))
	// pointID := string(r.Form.Get("pid"))
	var bean ServiceBean
	err := s.config.AddYtsPoint(channelName, deviceID, pointName, pointID)
	if err != nil {
		bean.Error = true
		bean.Message = err.Error()
	}
	bean.Error = false
	bytes, _ := json.Marshal(bean)
	fmt.Fprint(w, string(bytes))
}

//RemovePoint 删除点
func (s *Service) RemovePoint(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	// tmp := string(r.Form.Get("kind"))
	switch vars["kind"] {
	case "Ycs":
		s.removeYcsPoint(w, r)
	case "Yxs":
		s.removeYxsPoint(w, r)
	case "Yms":
		s.removeYmsPoint(w, r)
	case "Yks":
		s.removeYksPoint(w, r)
	case "Yts":
		s.removeYtsPoint(w, r)
	}
}

//removeYcsPoint 删除遥测点
func (s *Service) removeYcsPoint(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	// channelName := string(r.Form.Get("name"))
	// deviceID := string(r.Form.Get("did"))
	// pointID := string(r.Form.Get("pid"))
	var bean ServiceBean
	err := s.config.DeleteYcsPoint(vars["channelname"], vars["deviceid"], vars["pointid"])
	if err != nil {
		bean.Error = true
		bean.Message = err.Error()
	}
	bean.Error = false
	bytes, _ := json.Marshal(bean)
	fmt.Fprint(w, string(bytes))
}

//removeYxsPoint 删除遥信点
func (s *Service) removeYxsPoint(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	// channelName := string(r.Form.Get("name"))
	// deviceID := string(r.Form.Get("did"))
	// pointID := string(r.Form.Get("pid"))
	var bean ServiceBean
	err := s.config.DeleteYxsPoint(vars["channelname"], vars["deviceid"], vars["pointid"])
	if err != nil {
		bean.Error = true
		bean.Message = err.Error()
	}
	bean.Error = false
	bytes, _ := json.Marshal(bean)
	fmt.Fprint(w, string(bytes))
}

//removeYmsPoint 删除遥脉点
func (s *Service) removeYmsPoint(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	// channelName := string(r.Form.Get("name"))
	// deviceID := string(r.Form.Get("did"))
	// pointID := string(r.Form.Get("pid"))
	var bean ServiceBean
	err := s.config.DeleteYmsPoint(vars["channelname"], vars["deviceid"], vars["pointid"])
	if err != nil {
		bean.Error = true
		bean.Message = err.Error()
	}
	bean.Error = false
	bytes, _ := json.Marshal(bean)
	fmt.Fprint(w, string(bytes))
}

//removeYksPoint 删除遥控点
func (s *Service) removeYksPoint(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	// channelName := string(r.Form.Get("name"))
	// deviceID := string(r.Form.Get("did"))
	// pointID := string(r.Form.Get("pid"))
	var bean ServiceBean
	err := s.config.DeleteYksPoint(vars["channelname"], vars["deviceid"], vars["pointid"])
	if err != nil {
		bean.Error = true
		bean.Message = err.Error()
	}
	bean.Error = false
	bytes, _ := json.Marshal(bean)
	fmt.Fprint(w, string(bytes))
}

//removeYtsPoint 删除遥调点
func (s *Service) removeYtsPoint(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	// channelName := string(r.Form.Get("name"))
	// deviceID := string(r.Form.Get("did"))
	// pointID := string(r.Form.Get("pid"))
	var bean ServiceBean
	err := s.config.DeleteYtsPoint(vars["channelname"], vars["deviceid"], vars["pointid"])
	if err != nil {
		bean.Error = true
		bean.Message = err.Error()
	}
	bean.Error = false
	bytes, _ := json.Marshal(bean)
	fmt.Fprint(w, string(bytes))
}
