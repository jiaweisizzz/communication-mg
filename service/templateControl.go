package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mg-coms/config"
	"net/http"

	"github.com/gorilla/mux"
)

//Template 服务模板
type Template struct {
	ID      int
	Name    string
	Channel string
	Potocol string
}

//TemplateConfig 模板配置
type TemplateConfig struct {
	ID      int
	Name    string
	Channel string
	Potocol string
	Ycs     []config.Point //遥测点表
	Yxs     []config.Point //遥信点表
	Yms     []config.Point //遥脉点表
	Yks     []config.Point //遥控点表
	Yts     []config.Point //遥调点表
}

//SearchOneTemplate 查找某个模板
func (s *Service) SearchOneTemplate(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	var bean ServiceBean
	data, err := s.config.GetTemplateByName(vars["templatename"])
	if err != nil {
		bean.Error = true
		bean.Message = err.Error()
	}
	bean.Error = false
	bean.Data = data
	bytes, _ := json.Marshal(bean)
	fmt.Fprint(w, string(bytes))
}

//SearchAllTemplate 查找所有模板
func (s *Service) SearchAllTemplate(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var bean ServiceBean
	data, err := s.config.GetTemplateConfig()
	if err != nil {
		bean.Error = true
		bean.Message = err.Error()
	}
	var tl = make([]Template, 0)
	for _, tp := range data {
		t := Template{
			ID:      tp.ID,
			Name:    tp.Name,
			Channel: tp.Channel,
			Potocol: tp.Potocol,
		}
		tl = append(tl, t)
	}
	bean.Error = false
	bean.Data = tl
	bytes, _ := json.Marshal(bean)
	fmt.Fprint(w, string(bytes))
}

//AppendTemplate 增加模板
func (s *Service) AppendTemplate(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println("fffffffffffffffffffffffffffffffffffffffff", r.Form, r.Body, r.PostForm)
	var tp Template
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &tp)
	if err != nil {
		fmt.Println(err)
	}
	var bean ServiceBean
	// templateID := tp.ID
	templateName := tp.Name
	// templateChannel := tp.Channel
	templatePotocol := tp.Potocol
	err = s.config.AddTemplate(templateName, templatePotocol)
	if err != nil {
		bean.Error = true
		bean.Message = err.Error()
	}
	bean.Error = false
	bytes, _ := json.Marshal(bean)
	fmt.Fprint(w, string(bytes))
}

//ModifiTemplate 修改模板
func (s *Service) ModifiTemplate(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var tp TemplateConfig
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &tp)
	fmt.Println("=============", tp)
	if err != nil {
		fmt.Println(err)
	}
	var bean ServiceBean
	templateName := tp.Name
	templateycs := tp.Ycs
	templateyxs := tp.Yxs
	templateyms := tp.Yms
	templateyks := tp.Yks
	templateyts := tp.Yts
	err = s.config.SetTemplateById(templateName, templateycs, templateyxs, templateyms, templateyks, templateyts)
	if err != nil {
		bean.Error = true
		bean.Message = err.Error()
	}
	bean.Error = false
	bean.Message = "保存成功"
	bytes, _ := json.Marshal(bean)
	fmt.Fprint(w, string(bytes))
}

//RemoveTemplate 删除模板
func (s *Service) RemoveTemplate(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	// channelName := string(r.Form.Get("name"))
	var bean ServiceBean
	err := s.config.DeleteTemplate(vars["templatename"])
	if err != nil {
		bean.Error = true
		bean.Message = err.Error()
	}
	bean.Error = false
	bytes, _ := json.Marshal(bean)
	fmt.Fprint(w, string(bytes))
}
