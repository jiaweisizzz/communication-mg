package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

// User 服务用户
type User struct {
	Name     string
	Password string
	Browse   bool
	Operate  bool
	Maintain bool
}

//Login 登录
func (s *Service) Login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	usernameParam, found1 := r.Form["username"]
	passwordParam, found2 := r.Form["password"]
	if !(found1 && found2) {
		fmt.Fprint(w, "请勿非法访问")
		return
	}
	var bean ServiceBean
	username := usernameParam[0]
	password := passwordParam[0]

	u := "username:" + username + ",password:" + password
	fmt.Println(u)
	if username == "用户" && password == "123456" {
		bean.Error = true
		bean.Message = "登录成功"
	} else {
		bean.Error = false
		bean.Message = "用户名或密码不正确"
	}
	bytes, _ := json.Marshal(bean)
	fmt.Fprint(w, string(bytes))
}

//SearchUser　查用户
func (s *Service) SearchUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	// userName := string(r.Form["name"][0])
	var bean ServiceBean
	data, err := s.config.GetUserByName(vars["name"])
	if err != nil {
		bean.Error = true
		bean.Message = err.Error()
	}
	bean.Error = false
	bean.Data = data
	bytes, _ := json.Marshal(bean)
	fmt.Fprint(w, string(bytes))
}

//LoginUser 用户登录
func (s *Service) LoginUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	// userName := string(r.Form["name"][0])
	var bean ServiceBean
	data, err := s.config.GetUserByName(vars["name"])
	if err != nil {
		bean.Error = true
		bean.Message = err.Error()
	}
	if data.Password == vars["password"] {
		bean.Data = true
		bean.Error = false
	} else {
		bean.Error = true
		bean.Data = false
	}
	bytes, _ := json.Marshal(bean)
	fmt.Fprint(w, string(bytes))
}

//SearchAllUsers 查找所有用户
func (s *Service) SearchAllUsers(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var bean ServiceBean
	data, err := s.config.GetUserConfig()
	if err != nil {
		bean.Error = true
		bean.Message = err.Error()
	}
	var ul = make([]User, 0)
	for _, usr := range data {
		u := User{
			Name:     usr.Name,
			Browse:   usr.Browse,
			Operate:  usr.Operate,
			Maintain: usr.Maintain,
		}
		ul = append(ul, u)
	}
	bean.Error = false
	bean.Data = ul
	bytes, _ := json.Marshal(bean)
	fmt.Fprint(w, string(bytes))
}

//ModifiUser 改用户
func (s *Service) ModifiUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	// vars := mux.Vars(r)
	// w.WriteHeader(http.StatusOK)
	// userName := string(r.Form.Get("name"))
	// userPassword := string(r.Form.Get("password"))

	// var user map[string]interface{}
	// body, _ := ioutil.ReadAll(r.Body)
	// json.Unmarshal(body, &user)
	//
	// userName := user["name"].(string)
	// userPassword := user["password"].(string)
	fmt.Println("___________________________________________________", r.Form, r.Body)
	var user User
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &user)
	userName := user.Name
	userBrowse := user.Browse
	userOperate := user.Operate
	userMaintain := user.Maintain
	var bean ServiceBean
	data, err := s.config.SetUserByName(userName, userBrowse, userOperate, userMaintain)
	if err != nil {
		bean.Error = true
		bean.Message = err.Error()
		fmt.Println(err)
	}

	bean.Error = false
	bean.Data = data
	bytes, _ := json.Marshal(bean)
	fmt.Fprint(w, string(bytes))
}

//AppendUser 增用户
func (s *Service) AppendUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println("___________________________________________________", r.Form)
	var user User
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &user)

	userName := user.Name
	userPassword := user.Password
	userBrowse := user.Browse
	userOperate := user.Operate
	userMaintain := user.Maintain
	var bean ServiceBean
	err := s.config.AddUser(userName, userPassword, userBrowse, userOperate, userMaintain, &s.config.UsersConfigList)
	if err != nil {
		bean.Error = true
		bean.Message = err.Error()
	}
	bean.Error = false
	bytes, _ := json.Marshal(bean)
	fmt.Fprint(w, string(bytes))
}

//RemoveUser 删用户
func (s *Service) RemoveUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	// userName := string(r.Form["name"][0])
	var bean ServiceBean
	err := s.config.DeleteUser(vars["name"], s.config.UsersConfigList)
	if err != nil {
		bean.Error = true
		bean.Message = err.Error()
	}
	bean.Error = false
	bytes, _ := json.Marshal(bean)
	fmt.Fprint(w, string(bytes))
}
