package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mg-coms/config"
	"net/http"

	"github.com/gorilla/mux"
)

type PointaddrConfig struct {
	Name   string
	AddrYc []config.Ypoint
	AddrYx []config.Ypoint
	AddrYm []config.Ypoint
	AddrYk []config.Ypoint
	AddrYt []config.Ypoint
}

func (s *Service) SearchPointAddr(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	w.WriteHeader(http.StatusOK)
	var bean ServiceBean
	data, err := s.config.Getaddrlist(&s.config.AddrList)
	if err != nil {
		bean.Error = true
		bean.Message = err.Error()
	}
	bean.Error = false
	bean.Data = data
	bytes, _ := json.Marshal(bean)
	fmt.Fprint(w, string(bytes))
}

func (s *Service) ModifiPointAddr(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var pa PointaddrConfig
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &pa)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(pa)
	var bean ServiceBean
	cname := pa.Name
	pointaddryc := pa.AddrYc
	pointaddryx := pa.AddrYx
	pointaddrym := pa.AddrYm
	pointaddryk := pa.AddrYk
	pointaddryt := pa.AddrYt
	err = s.config.SetPointAddr(pointaddryc, pointaddryx, pointaddrym, pointaddryk, pointaddryt, cname)
	if err != nil {
		bean.Error = true
		bean.Message = err.Error()
	}
	bean.Error = false
	bytes, _ := json.Marshal(bean)
	fmt.Fprint(w, string(bytes))
	s.Restart()
}

func (s *Service) SearchNewPointAddr(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	var bean ServiceBean
	data, err := s.config.GetAddrPoint(vars["cname"])
	if err != nil {
		bean.Error = true
		bean.Message = err.Error()
	}
	bean.Error = false
	bean.Data = data
	bytes, _ := json.Marshal(bean)
	fmt.Fprint(w, string(bytes))
}

func (s *Service) DeletePointAddr(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var pa PointaddrConfig
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &pa)
	if err != nil {
		fmt.Println(err)
	}
	var bean ServiceBean
	cname := pa.Name
	pointaddryc := pa.AddrYc
	pointaddryx := pa.AddrYx
	pointaddrym := pa.AddrYm
	pointaddryk := pa.AddrYk
	pointaddryt := pa.AddrYt
	err = s.config.DeletePointaddr(pointaddryc, pointaddryx, pointaddrym, pointaddryk, pointaddryt, cname)
	if err != nil {
		bean.Error = true
		bean.Message = err.Error()
	}
	bean.Error = false
	bytes, _ := json.Marshal(bean)
	fmt.Fprint(w, string(bytes))
	s.Restart()
}
