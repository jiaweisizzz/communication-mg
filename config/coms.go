package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

//ComsConfig 通讯管理机配置
type ComsConfig struct {
	Key  string
	Port uint
}

//loadComs 加载通讯管理机配置
func loadComs(c *ComsConfig) error {
	f, err := ioutil.ReadFile("coms.json")
	if err != nil {
		return errors.New("Read Coms.json File error")
	}
	err = json.Unmarshal(f, c)
	if err != nil {
		return errors.New("Unmarshal error")
	}
	fmt.Println("主机配置加载完成")
	return nil
}
