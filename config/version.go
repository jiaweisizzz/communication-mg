package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

//VersionConfig 版本配置
type VersionConfig struct {
	Name   string
	Explan string
}

//AddVersion 增加版本
func (c *Config) AddVersion(name, explan string, v *[]VersionConfig) error {
	version := new(VersionConfig)
	version.Name = name
	version.Explan = explan
	*v = append(*v, *version)
	b, err := json.Marshal(v)
	if err != nil {
		return errors.New("Marshal error")
	}
	ioutil.WriteFile("version/version.json", b, os.ModePerm)
	return nil
}

//DeleteVersion 删除版本
func (c *Config) DeleteVersion(name string) error {
	var exist bool
	exist = false
	for index, version := range c.VersionConfigList {
		if name == version.Name {
			exist = true
			c.VersionConfigList = append(c.VersionConfigList[:index], c.VersionConfigList[index+1:]...)
			break
		}
	}
	if exist == false {
		return errors.New("版本不存在，删除失败")
	}

	b, err := json.Marshal(c.VersionConfigList)
	if err != nil {
		return errors.New("Marshal error")
	}
	ioutil.WriteFile("version/version.json", b, os.ModePerm)
	return nil
}

//loadVersion 加载版本
func loadVersion(v *[]VersionConfig) error {
	f, err := ioutil.ReadFile("version/version.json")
	if err != nil {
		return errors.New("Read Version.json File error")
	}
	err = json.Unmarshal(f, v)
	if err != nil {
		return errors.New("配置模板文件错误")
	}
	fmt.Println("配置模板加载完成")
	return nil
}
