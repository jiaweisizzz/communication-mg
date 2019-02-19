package service

import (
	"encoding/json"
	"fmt"
	"net/http"
)

//Process 运行指标
type Process struct {
	Cpu float64
	Mem float64
}

//SearchProcess 查看运行指标
func (s *Service) SearchProcess(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var bean ServiceBean
	var p Process
	sum, num := s.coms.GetProcess()
	p.Cpu = sum
	p.Mem = num
	bean.Error = false
	bean.Data = p
	bytes, _ := json.Marshal(bean)
	fmt.Fprint(w, string(bytes))
}
