package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

//UserConfig 用户配置
type UserConfig struct {
	Name     string
	Password string
	Browse   bool
	Operate  bool
	Maintain bool
}

//loadUsers 加载用户
func loadUsers(c *[]UserConfig) error {
	f, err := ioutil.ReadFile("users.json")
	if err != nil {
		return errors.New("Read users.json File error")
	}
	err = json.Unmarshal(f, c)
	if err != nil {
		return errors.New("用户配置文件错误")
	}
	fmt.Println("用户配置加载完成")
	return nil
}

//GetUserByName　通过名字获取用户
func (c *Config) GetUserByName(name string) (*UserConfig, error) {
	for _, u := range c.UsersConfigList {
		if name == u.Name {
			return &u, nil
		}
	}
	return nil, errors.New("用户不存在")
}

//SetUserByName 通过名字修改用户
func (c *Config) SetUserByName(name string, browse, operate, maintain bool) (*UserConfig, error) {
	var exist bool
	exist = false
	var user UserConfig
	for index, u := range c.UsersConfigList {
		if name == u.Name {
			exist = true
			c.UsersConfigList[index].Browse = browse
			c.UsersConfigList[index].Operate = operate
			c.UsersConfigList[index].Maintain = maintain
			user = c.UsersConfigList[index]
			break
		}
	}

	if exist == false {
		return nil, errors.New("用户不存在")
	}

	b, err := json.Marshal(c.UsersConfigList)
	if err != nil {
		return nil, errors.New("Marshal error")
	}
	ioutil.WriteFile("users.json", b, os.ModePerm)
	return &user, nil
}

//AddUser 增加用户
func (c *Config) AddUser(name, password string, browse, opretion, maintain bool, u *[]UserConfig) error {
	user := new(UserConfig)
	user.Name = name
	user.Password = password
	user.Browse = browse
	user.Operate = opretion
	user.Maintain = maintain

	*u = append(*u, *user)
	b, err := json.Marshal(u)
	if err != nil {
		return errors.New("Marshal error")
	}
	ioutil.WriteFile("users.json", b, os.ModePerm)
	return nil
}

//DeleteUser 删除用户
func (c *Config) DeleteUser(name string, u []UserConfig) error {
	var exist bool
	exist = false
	for index, user := range u {
		if name == user.Name {
			exist = true
			u = append(u[:index], u[index+1:]...)
			fmt.Println(u)
			break
		}
	}
	if exist == false {
		return errors.New("用户不存在，删除失败")
	}
	b, err := json.Marshal(u)
	if err != nil {
		return errors.New("Marshal error")
	}
	ioutil.WriteFile("users.json", b, os.ModePerm)
	return nil
}

// //Users 用户
// type Users struct {
// 	UsersConfigList []UserConfig
// }
//
// //GetUserConfig 获得用户配置
// func (u *Users) GetUserConfig() (Userconfig []UserConfig, err error) {
// 	usr := u.UsersConfigList
// 	if usr == nil {
// 		err = errors.New("用户不存在")
// 		return nil, err
// 	}
// 	return usr, nil
// }

//Init 规约初始化
// func (u *Users) Init() *Users {
// 	return u
// }
//
// //NewUsers 新建用户配置实例
// func NewUsers() *Users { return new(Users).Init() }
